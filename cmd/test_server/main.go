package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/josealecrim/audiototext/internal/hardware"
	"github.com/josealecrim/audiototext/internal/inference"
	"github.com/josealecrim/audiototext/internal/models"
	"github.com/josealecrim/audiototext/internal/server"
	pb "github.com/josealecrim/audiototext/pkg/transcription"
	"google.golang.org/grpc"
)

func main() {
	// Initialize hardware detector
	hwDetector, err := hardware.NewDetector()
	if err != nil {
		log.Fatalf("Failed to create hardware detector: %v", err)
	}

	if err := hwDetector.Start(); err != nil {
		log.Fatalf("Failed to start hardware detector: %v", err)
	}

	// Create model manager
	modelManager := models.NewManager()

	// Add test model
	testModel := &models.ONNXModel{
		Model: &models.Model{
			Path:     "models/whisper.onnx",
			Checksum: "test-checksum",
		},
		OptimizationLevel: 1,
		IsQuantized:       false,
		ID:                "test-model",
	}
	if err := modelManager.AddModel(testModel); err != nil {
		log.Fatalf("Failed to add test model: %v", err)
	}

	// Create inference manager
	inferenceManager := inference.NewManager(hwDetector)

	// Create gRPC server
	srv := server.NewServer(inferenceManager, modelManager)

	// Create gRPC server instance
	grpcServer := grpc.NewServer()
	pb.RegisterTranscriptionServiceServer(grpcServer, srv)

	// Create listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	// Start server
	log.Printf("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func testServer(client pb.TranscriptionServiceClient) {
	ctx := context.Background()

	// Test GetStatus
	status, err := client.GetStatus(ctx, &pb.GetStatusRequest{})
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}
	fmt.Printf("Server status: %+v\n", status)

	// Test GetModels
	models, err := client.GetModels(ctx, &pb.GetModelsRequest{})
	if err != nil {
		log.Fatalf("Failed to get models: %v", err)
	}
	fmt.Printf("Available models: %+v\n", models)

	// Test Transcribe
	audioData := make([]byte, 16000*2) // 1 second of 16-bit PCM audio at 16kHz
	resp, err := client.Transcribe(ctx, &pb.TranscribeRequest{
		AudioData: audioData,
		Format:    pb.AudioFormat_AUDIO_FORMAT_WAV,
		Config: &pb.TranscriptionConfig{
			ModelId:  "test-model",
			Language: "en-US",
		},
	})
	if err != nil {
		log.Fatalf("Failed to transcribe: %v", err)
	}
	fmt.Printf("Transcription result: %+v\n", resp)
}
