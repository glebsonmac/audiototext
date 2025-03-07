package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "github.com/glebsonmac/audiototext/go-server/pkg/pb/audio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddr := flag.String("server", "localhost:50051", "Endereço do servidor gRPC")
	audioFile := flag.String("audio", "", "Caminho do arquivo de áudio")
	flag.Parse()

	if *audioFile == "" {
		log.Fatal("Por favor, especifique um arquivo de áudio com -audio")
	}

	// Carregar arquivo de áudio
	audioData, err := os.ReadFile(*audioFile)
	if err != nil {
		log.Fatalf("Falha ao ler arquivo: %v", err)
	}

	// Conectar ao servidor
	conn, err := grpc.Dial(*serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Conexão falhou: %v", err)
	}
	defer conn.Close()

	client := pb.NewAudioServiceClient(conn)

	// Criar requisição
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.TranscribeAudio(ctx, &pb.TranscribeRequest{
		AudioData:   audioData,
		AudioFormat: pb.AudioFormat_AUDIO_FORMAT_MP3,
	})

	if err != nil {
		log.Fatalf("Erro na transcrição: %v", err)
	}

	if resp.Error != "" {
		log.Printf("Erro do servidor: %s", resp.Error)
		return
	}

	log.Printf("Transcrição:\n%s", resp.Transcript)
}
