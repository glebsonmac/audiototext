package main

import (
	"log"
	"net"

	"github.com/glebsonmac/audiototext/internal/handler"
	pb "github.com/glebsonmac/audiototext/pkg/pb/audio"
	"google.golang.org/grpc"
)

func main() {
	listenAddr := ":50051"
	whisperAddr := "localhost:50052" // Servidor whisper

	// Inicializa o servidor gRPC
	grpcServer := grpc.NewServer()
	audioHandler := handler.NewAudioHandler(whisperAddr)
	pb.RegisterAudioServiceServer(grpcServer, audioHandler)

	//iniciar listener

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
