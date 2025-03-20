package hardware_test

import (
	"context"
	"testing"
	"time"

	"github.com/josealecrim/audiototext/internal/hardware"
	"github.com/josealecrim/audiototext/test/helpers"
)

func TestCPUDetection(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 5*time.Second, func(ctx context.Context) {
		info, err := detector.GetCPUInfo(ctx)
		helpers.AssertNoError(t, err)

		// Basic CPU validations
		if info.NumCores < 1 {
			t.Error("CPU cores should be at least 1")
		}
		if info.Frequency <= 0 {
			t.Error("CPU frequency should be greater than 0")
		}
		if len(info.CacheInfo) == 0 {
			t.Error("CPU cache info should not be empty")
		}
	})
}

func TestGPUDetection(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 5*time.Second, func(ctx context.Context) {
		info, err := detector.GetGPUInfo(ctx)
		helpers.AssertNoError(t, err)

		// Even if no GPU is available, we should get valid info
		if info == nil {
			t.Error("GPU info should not be nil")
		}

		// Check GPU capabilities
		if info.HasCUDA {
			if info.CUDAVersion == "" {
				t.Error("CUDA version should be set when CUDA is available")
			}
			if info.AvailableMemory <= 0 {
				t.Error("GPU memory should be greater than 0 when CUDA is available")
			}
		}

		if info.HasOpenCL {
			if info.OpenCLVersion == "" {
				t.Error("OpenCL version should be set when OpenCL is available")
			}
		}
	})
}

func TestMemoryDetection(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 5*time.Second, func(ctx context.Context) {
		info, err := detector.GetMemoryInfo(ctx)
		helpers.AssertNoError(t, err)

		if info.TotalRAM <= 0 {
			t.Error("Total RAM should be greater than 0")
		}
		if info.AvailableRAM <= 0 {
			t.Error("Available RAM should be greater than 0")
		}
		if info.AvailableRAM > info.TotalRAM {
			t.Error("Available RAM cannot be greater than total RAM")
		}
	})
}

func TestResourceMonitoring(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	monitor, err := hardware.NewResourceMonitor(detector)
	helpers.AssertNoError(t, err)

	helpers.WithTimeout(t, 10*time.Second, func(ctx context.Context) {
		// Start monitoring
		err := monitor.Start(ctx)
		helpers.AssertNoError(t, err)

		// Wait for a few metrics to be collected
		time.Sleep(2 * time.Second)

		// Get current metrics
		metrics := monitor.GetCurrentMetrics()

		// Validate CPU metrics
		if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
			t.Error("CPU usage should be between 0 and 100")
		}

		// Validate memory metrics
		if metrics.MemoryUsage < 0 || metrics.MemoryUsage > 100 {
			t.Error("Memory usage should be between 0 and 100")
		}

		// If GPU is available, validate GPU metrics
		if metrics.HasGPU {
			if metrics.GPUUsage < 0 || metrics.GPUUsage > 100 {
				t.Error("GPU usage should be between 0 and 100")
			}
		}

		// Stop monitoring
		monitor.Stop()
	})
}
