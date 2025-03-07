package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	pb "github.com/glebsonmac/audiototext/pkg/pb/audio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Configurações
	grpcAddr := "localhost:50051"
	webPort := "8080"

	// Configurar cliente gRPC
	conn, err := grpc.Dial(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Falha ao conectar ao servidor gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewAudioServiceClient(conn)

	// Servir arquivos estáticos
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Endpoint de transcrição
	http.HandleFunc("/transcribe", func(w http.ResponseWriter, r *http.Request) {
		// Processar upload
		file, _, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Falha ao ler arquivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Ler dados do áudio
		audioData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erro ao ler dados do áudio", http.StatusInternalServerError)
			return
		}

		// Chamar servidor gRPC
		ctx := r.Context()
		resp, err := grpcClient.TranscribeAudio(ctx, &pb.TranscribeRequest{
			AudioData:   audioData,
			AudioFormat: pb.AudioFormat_AUDIO_FORMAT_MP3,
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Erro na transcrição: %v", err), http.StatusInternalServerError)
			return
		}

		// Retornar resposta
		w.Header().Set("Content-Type", "application/json")
		if resp.Error != "" {
			fmt.Fprintf(w, `{"error": %q}`, resp.Error)
		} else {
			fmt.Fprintf(w, `{"transcript": %q}`, resp.Transcript)
		}
	})

	log.Printf("Servidor web ouvindo em :%s", webPort)
	log.Fatal(http.ListenAndServe(":"+webPort, nil))
}
