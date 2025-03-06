package grpcclient

import (
	"context"
	"log"

	pb "github.com/glebsonmac/audiototext/pkg/pb/whisper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type whisperClient struct {
	conn   *grpc.ClientConn
	client pb.whisperServiceClient
}

func NewWhisperClient(addr string) *whisperClient {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to whisper service: %v", err)
	}

	client := pb.NewWhisperServiceClient(conn)
	return &whisperClient{conn, client}

}
func (c *whisperClient) Transcribe(ctx context.Context, audioData []byte, audioFormat pb.AudioFormat) (string, error) {
	resp, err := c.client.Transcribe(ctx, &pb.TranscribeRequest{
		AudioData:   audioData,
		AudioFormat: audioFormat,
	})

	if err != nil {
		return "", err
	}
	return resp.Transcript, nil

}
