package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CPUInfo contains information about the CPU
type CPUInfo struct {
	NumCores  int
	Frequency float64
	CacheInfo []CacheInfo
}

// CacheInfo contains information about CPU cache levels
type CacheInfo struct {
	Level int
	Size  int
}

// GPUInfo contains information about available GPUs
type GPUInfo struct {
	HasCUDA         bool
	HasOpenCL       bool
	CUDAVersion     string
	OpenCLVersion   string
	AvailableMemory int64
}

// MemoryInfo contains system memory information
type MemoryInfo struct {
	TotalRAM     int64
	AvailableRAM int64
}

// Detector provides hardware detection capabilities
type Detector struct {
	cpuInfo    []cpu.InfoStat
	memoryInfo *mem.VirtualMemoryStat
}

// NewDetector creates a new hardware detector
func NewDetector() (*Detector, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}

	return &Detector{
		cpuInfo:    cpuInfo,
		memoryInfo: memInfo,
	}, nil
}

// GetCPUInfo returns CPU information
func (d *Detector) GetCPUInfo() ([]cpu.InfoStat, error) {
	return d.cpuInfo, nil
}

// GetMemoryInfo returns memory information
func (d *Detector) GetMemoryInfo() (*mem.VirtualMemoryStat, error) {
	return d.memoryInfo, nil
}

// GetNumCPUs returns the number of CPUs
func (d *Detector) GetNumCPUs() int {
	return runtime.NumCPU()
}

// GetTotalMemory returns the total memory in bytes
func (d *Detector) GetTotalMemory() uint64 {
	return d.memoryInfo.Total
}

// GetAvailableMemory returns the available memory in bytes
func (d *Detector) GetAvailableMemory() uint64 {
	return d.memoryInfo.Available
}

// GetMemoryUsagePercent returns the memory usage percentage
func (d *Detector) GetMemoryUsagePercent() float64 {
	return d.memoryInfo.UsedPercent
}

// GetCPUModelName returns the CPU model name
func (d *Detector) GetCPUModelName() string {
	if len(d.cpuInfo) > 0 {
		return d.cpuInfo[0].ModelName
	}
	return "Unknown"
}

// GetCPUCores returns the number of CPU cores
func (d *Detector) GetCPUCores() int32 {
	if len(d.cpuInfo) > 0 {
		return d.cpuInfo[0].Cores
	}
	return 0
}

// GetCPUFrequency returns the CPU frequency in MHz
func (d *Detector) GetCPUFrequency() float64 {
	if len(d.cpuInfo) > 0 {
		return d.cpuInfo[0].Mhz
	}
	return 0
}

// GetHardwareInfo returns a summary of hardware information
func (d *Detector) GetHardwareInfo() map[string]interface{} {
	return map[string]interface{}{
		"cpu": map[string]interface{}{
			"model":     d.GetCPUModelName(),
			"cores":     d.GetCPUCores(),
			"frequency": d.GetCPUFrequency(),
		},
		"memory": map[string]interface{}{
			"total":     d.GetTotalMemory(),
			"available": d.GetAvailableMemory(),
			"usage":     d.GetMemoryUsagePercent(),
		},
	}
}
