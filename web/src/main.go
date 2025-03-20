package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	pb "github.com/josealecrim/audiototext/pkg/transcription"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Global state
	state = struct {
		client     pb.TranscriptionServiceClient
		recording  bool
		processing bool
		connected  bool
		language   string
		modelID    string
	}{
		recording:  false,
		processing: false,
		connected:  false,
		language:   "en-US",
		modelID:    "whisper-base",
	}

	// Global callbacks
	callbacks = make(map[string]js.Func)
)

func main() {
	fmt.Println("WebAssembly module loaded")

	// Register JavaScript callbacks
	registerCallbacks()

	// Create channel to keep main goroutine alive
	done := make(chan struct{})
	<-done
}

func registerCallbacks() {
	// Initialize gRPC client
	callbacks["initClient"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return js.ValueOf(map[string]interface{}{
				"error": "server address is required",
			})
		}

		serverAddr := args[0].String()
		err := initClient(serverAddr)
		if err != nil {
			return js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return js.ValueOf(map[string]interface{}{
			"success": true,
		})
	})

	// Get available models
	callbacks["getModels"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !state.connected {
			return js.ValueOf(map[string]interface{}{
				"error": "not connected to server",
			})
		}

		models, err := getModels()
		if err != nil {
			return js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return js.ValueOf(map[string]interface{}{
			"models": models,
		})
	})

	// Start recording
	callbacks["startRecording"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if state.recording {
			return js.ValueOf(map[string]interface{}{
				"error": "already recording",
			})
		}

		state.recording = true
		return js.ValueOf(map[string]interface{}{
			"success": true,
		})
	})

	// Stop recording
	callbacks["stopRecording"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !state.recording {
			return js.ValueOf(map[string]interface{}{
				"error": "not recording",
			})
		}

		state.recording = false
		return js.ValueOf(map[string]interface{}{
			"success": true,
		})
	})

	// Process audio
	callbacks["processAudio"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return js.ValueOf(map[string]interface{}{
				"error": "audio data is required",
			})
		}

		audioData := make([]byte, args[0].Length())
		js.CopyBytesToGo(audioData, args[0])

		result, err := processAudio(audioData)
		if err != nil {
			return js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return js.ValueOf(map[string]interface{}{
			"text":       result.Text,
			"confidence": result.Confidence,
		})
	})

	// Set language
	callbacks["setLanguage"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return js.ValueOf(map[string]interface{}{
				"error": "language code is required",
			})
		}

		state.language = args[0].String()
		return js.ValueOf(map[string]interface{}{
			"success": true,
		})
	})

	// Set model
	callbacks["setModel"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return js.ValueOf(map[string]interface{}{
				"error": "model ID is required",
			})
		}

		state.modelID = args[0].String()
		return js.ValueOf(map[string]interface{}{
			"success": true,
		})
	})

	// Get server status
	callbacks["getStatus"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !state.connected {
			return js.ValueOf(map[string]interface{}{
				"error": "not connected to server",
			})
		}

		status, err := getStatus()
		if err != nil {
			return js.ValueOf(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return js.ValueOf(map[string]interface{}{
			"status": status,
		})
	})

	// Register callbacks in JavaScript
	js.Global().Set("audioToText", js.ValueOf(map[string]interface{}{
		"initClient":     callbacks["initClient"],
		"getModels":      callbacks["getModels"],
		"startRecording": callbacks["startRecording"],
		"stopRecording":  callbacks["stopRecording"],
		"processAudio":   callbacks["processAudio"],
		"setLanguage":    callbacks["setLanguage"],
		"setModel":       callbacks["setModel"],
		"getStatus":      callbacks["getStatus"],
	}))
}

func initClient(serverAddr string) error {
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	state.client = pb.NewTranscriptionServiceClient(conn)
	state.connected = true
	return nil
}

func getModels() (interface{}, error) {
	ctx := context.Background()
	resp, err := state.client.GetModels(ctx, &pb.GetModelsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %v", err)
	}

	// Convert to JSON-friendly format
	models := make([]map[string]interface{}, len(resp.Models))
	for i, model := range resp.Models {
		models[i] = map[string]interface{}{
			"id":                  model.Id,
			"name":                model.Name,
			"description":         model.Description,
			"languages":           model.Languages,
			"size":                model.Size,
			"supportsStreaming":   model.SupportsStreaming,
			"supportsDiarization": model.SupportsDiarization,
			"version":             model.Version,
		}
	}

	return models, nil
}

func processAudio(audioData []byte) (*pb.TranscribeResponse, error) {
	ctx := context.Background()
	resp, err := state.client.Transcribe(ctx, &pb.TranscribeRequest{
		AudioData: audioData,
		Format:    pb.AudioFormat_AUDIO_FORMAT_WAV,
		Config: &pb.TranscriptionConfig{
			ModelId:  state.modelID,
			Language: state.language,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe: %v", err)
	}

	return resp, nil
}

func getStatus() (interface{}, error) {
	ctx := context.Background()
	resp, err := state.client.GetStatus(ctx, &pb.GetStatusRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %v", err)
	}

	// Convert to JSON-friendly format
	status := map[string]interface{}{
		"isReady":        resp.IsReady,
		"load":           resp.Load,
		"activeSessions": resp.ActiveSessions,
		"memoryUsage":    resp.MemoryUsage,
		"gpuMemoryUsage": resp.GpuMemoryUsage,
		"details":        resp.Details,
	}

	return status, nil
}
