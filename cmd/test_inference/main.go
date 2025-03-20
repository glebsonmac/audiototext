package main

import (
	"context"
	"fmt"
	"log"

	"github.com/josealecrim/audiototext/internal/hardware"
	"github.com/josealecrim/audiototext/internal/inference"
	"github.com/josealecrim/audiototext/internal/models"
)

func main() {
	// Initialize hardware detector
	hwDetector := hardware.NewDetector()
	if err := hwDetector.Initialize(); err != nil {
		log.Fatalf("Failed to initialize hardware detector: %v", err)
	}

	// Print hardware information
	printHardwareInfo(hwDetector)

	// Create session manager
	manager := inference.NewManager(hwDetector)

	// Create a test model
	model := &models.ONNXModel{
		Model: &models.Model{
			Path:     "models/whisper.onnx",
			Checksum: "test-checksum",
		},
		OptimizationLevel: 1,
		IsQuantized:       false,
		ID:                "test-model",
	}

	// Create session configuration
	config := &inference.SessionConfig{
		ExecutionProvider:      inference.CPUExecutionProvider,
		InterOpNumThreads:      4,
		IntraOpNumThreads:      4,
		GraphOptimizationLevel: 99,
		EnableMemoryPattern:    true,
		EnableCPUMemArena:      true,
		EnableProfiling:        true,
	}

	// Create batch configuration
	batchConfig := &inference.BatchConfig{
		MaxBatchSize:        8,
		DynamicBatching:     true,
		MaxLatencyMs:        100,
		PreferredBatchSizes: []int{1, 2, 4, 8},
	}

	// Create session
	ctx := context.Background()
	session, err := manager.CreateSession(ctx, model, config, batchConfig)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer manager.CloseSession(model.ID)

	// Create inference handler
	inf := inference.NewInference(session)

	// Create test audio data (1 second of silence at 16kHz)
	audioData := make([]float32, 16000)

	// Test single inference
	fmt.Println("\nTesting single inference...")
	result, err := inf.ProcessAudio(ctx, audioData)
	if err != nil {
		log.Fatalf("Failed to process audio: %v", err)
	}
	printResult(result)

	// Test batch inference
	fmt.Println("\nTesting batch inference...")
	batch := make([][]float32, 4)
	for i := range batch {
		batch[i] = make([]float32, 16000)
	}

	results, err := inf.ProcessBatch(ctx, batch)
	if err != nil {
		log.Fatalf("Failed to process batch: %v", err)
	}
	for i, result := range results {
		fmt.Printf("\nBatch result %d:\n", i+1)
		printResult(result)
	}

	// Print final statistics
	fmt.Println("\nFinal statistics:")
	printStats(inf.GetStats())
}

func printHardwareInfo(hwDetector *hardware.Detector) {
	fmt.Println("Hardware Information:")

	// Print CPU info
	cpuInfo := hwDetector.GetCPUInfo()
	fmt.Printf("\nCPU:\n")
	fmt.Printf("  Model: %s\n", cpuInfo.Model)
	fmt.Printf("  Cores: %d\n", cpuInfo.Cores)
	fmt.Printf("  Frequency: %.2f MHz\n", cpuInfo.Frequency)
	fmt.Printf("  Cache Size: %d bytes\n", cpuInfo.CacheSize)

	// Print GPU info
	gpuInfo := hwDetector.GetGPUInfo()
	fmt.Printf("\nGPUs:\n")
	if len(gpuInfo) == 0 {
		fmt.Println("  No GPUs detected")
	} else {
		for i, gpu := range gpuInfo {
			fmt.Printf("  GPU %d:\n", i+1)
			fmt.Printf("    Name: %s\n", gpu.Name)
			fmt.Printf("    Memory: %d bytes\n", gpu.Memory)
			fmt.Printf("    CUDA Support: %v\n", gpu.HasCUDA)
			fmt.Printf("    OpenCL Support: %v\n", gpu.HasOpenCL)
		}
	}

	// Print memory info
	memInfo := hwDetector.GetMemoryInfo()
	fmt.Printf("\nMemory:\n")
	fmt.Printf("  Total: %d bytes\n", memInfo.Total)
	fmt.Printf("  Available: %d bytes\n", memInfo.Available)
	fmt.Printf("  Used: %d bytes\n", memInfo.Used)
}

func printResult(result *inference.Result) {
	fmt.Printf("  Transcription: %s\n", result.Transcription)
	fmt.Printf("  Confidence: %.2f\n", result.Confidence)
	fmt.Printf("  Timestamp: %.2f - %.2f seconds\n", result.TimestampStart, result.TimestampEnd)
	fmt.Printf("  Processing Time: %.3f seconds\n", result.ProcessingTime)
}

func printStats(stats inference.Stats) {
	fmt.Printf("  Total Inferences: %d\n", stats.TotalInferences)
	fmt.Printf("  Average Processing Time: %.3f seconds\n", stats.AverageProcessingTime)
	fmt.Printf("  Total Processing Time: %.3f seconds\n", stats.TotalProcessingTime)
	fmt.Printf("  Total Batches: %d\n", stats.TotalBatches)
	fmt.Printf("  Average Batch Size: %.2f\n", stats.AverageBatchSize)
	if stats.PeakMemoryUsage > 0 {
		fmt.Printf("  Peak Memory Usage: %d bytes\n", stats.PeakMemoryUsage)
	}
	if stats.CurrentMemoryUsage > 0 {
		fmt.Printf("  Current Memory Usage: %d bytes\n", stats.CurrentMemoryUsage)
	}
}
