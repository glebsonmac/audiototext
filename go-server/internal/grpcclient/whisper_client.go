package grpcclient

import (
	"context"
	"log"

	pb "github.com/glebsonmac/audiototext/pkg/pb/audio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WhisperClient struct {
	conn   *grpc.ClientConn
	client pb.AudioServiceClient
}

func NewWhisperClient(addr string) *WhisperClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to whisper service: %v", err)
	}

	client := pb.NewAudioServiceClient(conn)
	return &WhisperClient{conn, client}

}
func (c *WhisperClient) Transcribe(ctx context.Context, audioData []byte, audioFormat pb.AudioFormat) (string, error) {
	resp, err := c.client.TranscribeAudio(ctx, &pb.TranscribeRequest{
		AudioData:   audioData,
		AudioFormat: audioFormat,
	})

	if err != nil {
		return "", err
	}
	return resp.Transcript, nil

}
