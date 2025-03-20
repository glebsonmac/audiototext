package models_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/josealecrim/audiototext/internal/models"
	"github.com/josealecrim/audiototext/test/helpers"
)

func TestModelDownload(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)
	manager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 2*time.Minute, func(ctx context.Context) {
		// Test downloading a small test model
		modelID := "openai/whisper-tiny"
		err := manager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Verify model files exist
		modelPath := filepath.Join(cfg.ModelDir, modelID)
		if !manager.ModelExists(modelID) {
			t.Errorf("Model files not found in %s", modelPath)
		}

		// Test cache functionality
		start := time.Now()
		err = manager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)
		if time.Since(start) > 5*time.Second {
			t.Error("Cache doesn't seem to be working - download took too long")
		}
	})
}

func TestModelConversion(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)
	manager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 5*time.Minute, func(ctx context.Context) {
		modelID := "openai/whisper-tiny"

		// Download model first
		err := manager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Convert to ONNX
		onnxPath, err := manager.ConvertToONNX(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Verify ONNX file exists
		if !manager.ONNXModelExists(modelID) {
			t.Errorf("ONNX model not found at %s", onnxPath)
		}

		// Test model validation
		err = manager.ValidateONNXModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Test optimization
		optimizedPath, err := manager.OptimizeONNXModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		if !manager.OptimizedModelExists(modelID) {
			t.Errorf("Optimized model not found at %s", optimizedPath)
		}
	})
}

func TestVersionManagement(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)
	manager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, time.Minute, func(ctx context.Context) {
		modelID := "openai/whisper-tiny"

		// Get initial version
		version1, err := manager.GetModelVersion(modelID)
		helpers.AssertNoError(t, err)

		// Download model
		err = manager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Get version after download
		version2, err := manager.GetModelVersion(modelID)
		helpers.AssertNoError(t, err)

		if version1 == version2 {
			t.Error("Version should have changed after download")
		}

		// Test backup functionality
		err = manager.CreateBackup(ctx, modelID)
		helpers.AssertNoError(t, err)

		backups, err := manager.ListBackups(modelID)
		helpers.AssertNoError(t, err)

		if len(backups) == 0 {
			t.Error("No backups found after creating one")
		}

		// Test restore functionality
		err = manager.RestoreFromBackup(ctx, modelID, backups[0])
		helpers.AssertNoError(t, err)

		version3, err := manager.GetModelVersion(modelID)
		helpers.AssertNoError(t, err)

		if version3 != version2 {
			t.Error("Version changed after restore when it shouldn't")
		}
	})
}

func TestCacheManagement(t *testing.T) {
	cfg := helpers.SetupTestEnv(t)
	manager, err := models.NewManager(cfg.ModelDir, cfg.CacheDir)
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, time.Minute, func(ctx context.Context) {
		// Download a model to populate cache
		modelID := "openai/whisper-tiny"
		err := manager.DownloadModel(ctx, modelID)
		helpers.AssertNoError(t, err)

		// Get cache size
		size, err := manager.GetCacheSize()
		helpers.AssertNoError(t, err)
		if size == 0 {
			t.Error("Cache size should not be 0 after downloading model")
		}

		// Clear cache
		err = manager.ClearCache()
		helpers.AssertNoError(t, err)

		// Verify cache is empty
		size, err = manager.GetCacheSize()
		helpers.AssertNoError(t, err)
		if size != 0 {
			t.Error("Cache size should be 0 after clearing")
		}
	})
}
