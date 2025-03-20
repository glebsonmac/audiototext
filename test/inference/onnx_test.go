package inference_test

import (
	"context"
	"testing"
	"time"

	"github.com/josealecrim/audiototext/internal/inference"
	"github.com/josealecrim/audiototext/internal/models"
	"github.com/josealecrim/audiototext/test/helpers"
)

func TestONNXInference(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)

	// Setup model manager and download test model
	modelManager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	modelID := "openai/whisper-tiny"

	helpers.WithTimeout(t, 5*time.Minute, func(ctx context.Context) {
		// Download and convert model
		err := modelManager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		onnxPath, err := modelManager.ConvertToONNX(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Create ONNX runtime
		runtime, err := inference.NewONNXRuntime(onnxPath)
		helpers.AssertNoError(t, err)

		// Test CPU inference
		err = runtime.SetExecutionProvider("CPU")
		helpers.AssertNoError(t, err)

		// Create sample input (5 seconds of audio at 16kHz)
		sampleRate := 16000
		duration := 5 * time.Second
		samples := make([]float32, int(duration.Seconds()*float64(sampleRate)))

		// Run inference
		result, err := runtime.RunInference(ctx, samples)
		helpers.AssertNoError(t, err)

		if len(result.Text) == 0 {
			t.Error("Inference result should not be empty")
		}

		// Test GPU inference if available
		if runtime.HasGPUSupport() {
			err = runtime.SetExecutionProvider("CUDA")
			if err != nil {
				t.Log("CUDA not available, skipping GPU test")
			} else {
				result, err = runtime.RunInference(ctx, samples)
				helpers.AssertNoError(t, err)

				if len(result.Text) == 0 {
					t.Error("GPU inference result should not be empty")
				}
			}
		}
	})
}

func TestBatchProcessing(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)

	modelManager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	modelID := "openai/whisper-tiny"

	helpers.WithTimeout(t, 5*time.Minute, func(ctx context.Context) {
		// Setup model
		err := modelManager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		onnxPath, err := modelManager.ConvertToONNX(ctx, modelID)
		helpers.AssertNoError(t, err)

		runtime, err := inference.NewONNXRuntime(onnxPath)
		helpers.AssertNoError(t, err)

		// Create batch processor
		batchSize := 4
		processor, err := inference.NewBatchProcessor(runtime, batchSize)
		helpers.AssertNoError(t, err)

		// Create sample inputs
		sampleRate := 16000
		duration := 3 * time.Second
		samplesPerInput := int(duration.Seconds() * float64(sampleRate))

		inputs := make([][]float32, batchSize)
		for i := range inputs {
			inputs[i] = make([]float32, samplesPerInput)
		}

		// Run batch inference
		results, err := processor.ProcessBatch(ctx, inputs)
		helpers.AssertNoError(t, err)

		if len(results) != batchSize {
			t.Errorf("Expected %d results, got %d", batchSize, len(results))
		}

		for i, result := range results {
			if len(result.Text) == 0 {
				t.Errorf("Result %d should not be empty", i)
			}
		}
	})
}

func TestFallbackBehavior(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)

	modelManager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	modelID := "openai/whisper-tiny"

	helpers.WithTimeout(t, 5*time.Minute, func(ctx context.Context) {
		// Setup model
		err := modelManager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		onnxPath, err := modelManager.ConvertToONNX(ctx, modelID)
		helpers.AssertNoError(t, err)

		runtime, err := inference.NewONNXRuntime(onnxPath)
		helpers.AssertNoError(t, err)

		// Test fallback behavior
		providers := []string{"CUDA", "DirectML", "OpenVINO", "CPU"}

		var selectedProvider string
		for _, provider := range providers {
			err = runtime.SetExecutionProvider(provider)
			if err == nil {
				selectedProvider = provider
				break
			}
		}

		if selectedProvider == "" {
			t.Fatal("No execution provider available")
		}

		// Verify inference works with selected provider
		samples := make([]float32, 16000) // 1 second of audio
		result, err := runtime.RunInference(ctx, samples)
		helpers.AssertNoError(t, err)

		if len(result.Text) == 0 {
			t.Errorf("Inference with %s provider failed", selectedProvider)
		}
	})
}
