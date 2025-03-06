package handler

import (
	"context"
	"log"

	"github.com/glebsonmac/audiototext/internal/grpcclient"
	pb "github.com/glebsonmac/audiototext/pkg/pb/audio"
)

type AudioHandler struct {
	pb.UnimplementedAudioServiceServer
	whisperClient *grpcclient.whisperClient
}

func NewAudioHandler(whisperAddr sting) *AudioHandler {
	return &AudioHandler{
		whisperClient: grpcclient.NewWhisperClient(whisperAddr),
	}
}

func (h *AudioHandler) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	// encaminhar para o serviço whisper
	transcript, err := h.whisperClient.Transcribe(ctx, req.AudioData, req.AudioFormat)
	if err != nil {
		log.Printf("failed to transcribe audio: %v", err)
		return &pb.TranscribeResponse{Error: err.Error()}, nil
	}
	return &pb.TranscribeResponse{Transcript: transcript}, nil
}
