package server

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/josealecrim/audiototext/internal/inference"
	"github.com/josealecrim/audiototext/internal/models"
	pb "github.com/josealecrim/audiototext/pkg/transcription"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the TranscriptionService gRPC server
type Server struct {
	pb.UnimplementedTranscriptionServiceServer
	// inferenceManager manages ONNX Runtime sessions
	inferenceManager *inference.Manager
	// modelManager manages model downloads and caching
	modelManager *models.Manager
	// activeSessions tracks active transcription sessions
	activeSessions sync.Map
	// stats tracks server statistics
	stats *Stats
}

// Stats tracks server statistics
type Stats struct {
	mu sync.RWMutex
	// activeSessionCount is the number of active transcription sessions
	activeSessionCount int32
	// totalRequests is the total number of requests processed
	totalRequests int64
	// totalErrors is the total number of errors encountered
	totalErrors int64
	// startTime is when the server started
	startTime time.Time
}

// NewServer creates a new transcription server
func NewServer(inferenceManager *inference.Manager, modelManager *models.Manager) *Server {
	return &Server{
		inferenceManager: inferenceManager,
		modelManager:     modelManager,
		stats:            &Stats{startTime: time.Now()},
	}
}

// Transcribe implements the one-shot transcription RPC
func (s *Server) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	if err := s.validateTranscribeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get the model
	model, err := s.modelManager.GetModel(req.Config.ModelId)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("model not found: %v", err))
	}

	// Create inference session
	session, err := s.inferenceManager.CreateSession(ctx, model, nil, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create inference session: %v", err))
	}
	defer s.inferenceManager.CloseSession(model.ID)

	// Create inference handler
	inf := inference.NewInference(session)

	// Convert audio data to float32 array
	audioData, err := s.convertAudioData(req.AudioData, req.Format)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to convert audio data: %v", err))
	}

	// Process audio
	result, err := inf.ProcessAudio(ctx, audioData)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to process audio: %v", err))
	}

	// Convert result to response
	response := s.convertResultToResponse(result)

	return response, nil
}

// TranscribeStream implements the streaming transcription RPC
func (s *Server) TranscribeStream(stream pb.TranscriptionService_TranscribeStreamServer) error {
	ctx := stream.Context()
	sessionID := fmt.Sprintf("stream-%d", time.Now().UnixNano())

	// Create channels for audio processing
	audioCh := make(chan []float32, 10)
	resultCh := make(chan *inference.Result, 10)
	errorCh := make(chan error, 1)

	// Start processing goroutine
	go s.processAudioStream(ctx, sessionID, audioCh, resultCh, errorCh)

	// Start result sending goroutine
	go s.sendResults(stream, resultCh, errorCh)

	// Receive audio chunks
	var config *pb.TranscriptionConfig
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			close(audioCh)
			return nil
		}
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to receive audio chunk: %v", err))
		}

		// Get config from first chunk
		if config == nil {
			config = chunk.Config
			if err := s.validateConfig(config); err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// Convert audio data
		audioData, err := s.convertAudioData(chunk.AudioData, chunk.Format)
		if err != nil {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("failed to convert audio data: %v", err))
		}

		// Send audio data for processing
		select {
		case audioCh <- audioData:
		case err := <-errorCh:
			return status.Error(codes.Internal, fmt.Sprintf("processing error: %v", err))
		case <-ctx.Done():
			return status.Error(codes.Canceled, "stream canceled")
		}
	}
}

// GetModels implements the GetModels RPC
func (s *Server) GetModels(ctx context.Context, req *pb.GetModelsRequest) (*pb.GetModelsResponse, error) {
	models, err := s.modelManager.ListModels()
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list models: %v", err))
	}

	response := &pb.GetModelsResponse{
		Models: make([]*pb.Model, 0, len(models)),
	}

	for _, model := range models {
		info := s.convertModelToInfo(model)
		response.Models = append(response.Models, info)
	}

	return response, nil
}

// GetStatus implements the GetStatus RPC
func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	s.stats.mu.RLock()
	defer s.stats.mu.RUnlock()

	// Get hardware stats
	memStats := s.inferenceManager.GetStats()
	var gpuMemoryUsage int64
	var currentMemoryUsage int64

	// Get the first model's stats or use 0 if no models are loaded
	for modelID, stats := range memStats {
		if stats.PeakMemoryUsage > gpuMemoryUsage {
			gpuMemoryUsage = stats.PeakMemoryUsage
		}
		if modelID == "default" {
			currentMemoryUsage = stats.CurrentMemoryUsage
		}
	}

	response := &pb.GetStatusResponse{
		IsReady:        true,
		Load:           float32(s.stats.activeSessionCount) / 100.0, // Arbitrary scale
		ActiveSessions: s.stats.activeSessionCount,
		MemoryUsage:    currentMemoryUsage,
		GpuMemoryUsage: gpuMemoryUsage,
		Details: map[string]string{
			"uptime":         time.Since(s.stats.startTime).String(),
			"total_requests": fmt.Sprintf("%d", s.stats.totalRequests),
			"total_errors":   fmt.Sprintf("%d", s.stats.totalErrors),
		},
	}

	return response, nil
}

// Helper functions

func (s *Server) validateTranscribeRequest(req *pb.TranscribeRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if len(req.AudioData) == 0 {
		return errors.New("audio data cannot be empty")
	}
	if req.Config == nil {
		return errors.New("config cannot be nil")
	}
	return s.validateConfig(req.Config)
}

func (s *Server) validateConfig(config *pb.TranscriptionConfig) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}
	if config.ModelId == "" {
		return errors.New("model ID cannot be empty")
	}
	if config.Language == "" {
		return errors.New("language cannot be empty")
	}
	return nil
}

func (s *Server) convertAudioData(data []byte, format pb.AudioFormat) ([]float32, error) {
	// TODO: Implement actual audio conversion
	// This will be implemented in the audio processing sprint
	return make([]float32, len(data)/4), nil
}

func (s *Server) convertResultToResponse(result *inference.Result) *pb.TranscribeResponse {
	response := &pb.TranscribeResponse{
		Text:       result.Transcription,
		Confidence: result.Confidence,
		Segments: []*pb.Segment{
			{
				Text:       result.Transcription,
				StartTime:  result.TimestampStart,
				EndTime:    result.TimestampEnd,
				Confidence: result.Confidence,
			},
		},
		Metadata: map[string]string{
			"processing_time": fmt.Sprintf("%.3f", result.ProcessingTime),
		},
	}

	return response
}

func (s *Server) processAudioStream(ctx context.Context, sessionID string, audioCh <-chan []float32, resultCh chan<- *inference.Result, errorCh chan<- error) {
	defer close(resultCh)
	defer close(errorCh)

	// TODO: Implement actual streaming audio processing
	// This will be implemented when we integrate streaming support
	for audioData := range audioCh {
		select {
		case <-ctx.Done():
			return
		default:
			// Process audio chunk
			result := &inference.Result{
				Transcription:  "Placeholder transcription",
				Confidence:     1.0,
				TimestampStart: 0,
				TimestampEnd:   float32(len(audioData)) / 16000,
				ProcessingTime: 0.1,
			}
			resultCh <- result
		}
	}
}

func (s *Server) sendResults(stream pb.TranscriptionService_TranscribeStreamServer, resultCh <-chan *inference.Result, errorCh chan error) {
	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				return
			}
			response := &pb.TranscribeResponse{
				Text:       result.Transcription,
				Confidence: result.Confidence,
				Segments: []*pb.Segment{
					{
						Text:       result.Transcription,
						StartTime:  result.TimestampStart,
						EndTime:    result.TimestampEnd,
						Confidence: result.Confidence,
					},
				},
			}
			if err := stream.Send(response); err != nil {
				select {
				case errorCh <- err:
				default:
				}
				return
			}
		case <-stream.Context().Done():
			return
		}
	}
}

func (s *Server) convertModelToInfo(model *models.ONNXModel) *pb.Model {
	return &pb.Model{
		Id:                  model.ID,
		Name:                "Whisper Model", // TODO: Get actual name
		Description:         "Automatic Speech Recognition Model",
		Languages:           []string{"en-US", "pt-BR"}, // TODO: Get actual languages
		Size:                0,                          // TODO: Get actual size
		SupportsStreaming:   true,
		SupportsDiarization: false,
		Version:             "1.0", // TODO: Get actual version
	}
}
