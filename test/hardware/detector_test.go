package hardware_test

import (
	"testing"
	"time"

	"github.com/josealecrim/audiototext/internal/hardware"
	"github.com/josealecrim/audiototext/test/helpers"
)

func TestCPUDetection(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	// Test CPU info
	cpuInfo, err := detector.GetCPUInfo()
	helpers.AssertNoError(t, err)

	// Basic CPU validations
	if len(cpuInfo) < 1 {
		t.Error("Should detect at least one CPU")
	}

	// Test CPU cores
	numCPUs := detector.GetNumCPUs()
	if numCPUs < 1 {
		t.Error("Number of CPUs should be at least 1")
	}

	// Test CPU model name
	modelName := detector.GetCPUModelName()
	if modelName == "" {
		t.Error("CPU model name should not be empty")
	}

	// Test CPU frequency
	cpuFreq := detector.GetCPUFrequency()
	if cpuFreq <= 0 {
		t.Error("CPU frequency should be greater than 0")
	}

	// Test CPU cores count
	cpuCores := detector.GetCPUCores()
	if cpuCores < 1 {
		t.Error("CPU cores should be at least 1")
	}
}

func TestMemoryDetection(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	// Get memory info
	memInfo, err := detector.GetMemoryInfo()
	helpers.AssertNoError(t, err)

	// Basic memory validations
	if memInfo.Total <= 0 {
		t.Error("Total memory should be greater than 0")
	}
	if memInfo.Available <= 0 {
		t.Error("Available memory should be greater than 0")
	}
	if memInfo.Available > memInfo.Total {
		t.Error("Available memory cannot be greater than total memory")
	}

	// Test memory getter methods
	totalMem := detector.GetTotalMemory()
	if totalMem <= 0 {
		t.Error("Total memory should be greater than 0")
	}

	availableMem := detector.GetAvailableMemory()
	if availableMem <= 0 {
		t.Error("Available memory should be greater than 0")
	}

	memUsage := detector.GetMemoryUsagePercent()
	if memUsage < 0 || memUsage > 100 {
		t.Error("Memory usage percentage should be between 0 and 100")
	}
}

func TestHardwareInfo(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	// Get hardware info summary
	hwInfo := detector.GetHardwareInfo()

	// Validate CPU info in summary
	cpuInfo, ok := hwInfo["cpu"].(map[string]interface{})
	if !ok {
		t.Error("Hardware info should contain CPU information")
	} else {
		if _, exists := cpuInfo["cores"]; !exists {
			t.Error("CPU info should contain cores information")
		}
		if _, exists := cpuInfo["frequency"]; !exists {
			t.Error("CPU info should contain frequency information")
		}
		if _, exists := cpuInfo["model"]; !exists {
			t.Error("CPU info should contain model information")
		}
	}

	// Validate memory info in summary
	memInfo, ok := hwInfo["memory"].(map[string]interface{})
	if !ok {
		t.Error("Hardware info should contain memory information")
	} else {
		if _, exists := memInfo["total"]; !exists {
			t.Error("Memory info should contain total information")
		}
		if _, exists := memInfo["available"]; !exists {
			t.Error("Memory info should contain available information")
		}
		if _, exists := memInfo["usage"]; !exists {
			t.Error("Memory info should contain usage information")
		}
	}
}

func TestResourceMonitoring(t *testing.T) {
	detector, err := hardware.NewDetector()
	helpers.AssertNoError(t, err)

	monitor, err := hardware.NewMonitor(detector)
	helpers.AssertNoError(t, err)

	// Start monitoring
	err = monitor.Start()
	helpers.AssertNoError(t, err)

	// Wait for a few metrics to be collected
	time.Sleep(2 * time.Second)

	// Get current metrics
	metrics := monitor.GetMetrics()

	// Validate CPU metrics
	if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
		t.Error("CPU usage should be between 0 and 100")
	}

	// Validate memory metrics
	if metrics.MemoryUsage < 0 || metrics.MemoryUsage > 100 {
		t.Error("Memory usage should be between 0 and 100")
	}

	// Stop monitoring
	err = monitor.Stop()
	helpers.AssertNoError(t, err)
}
