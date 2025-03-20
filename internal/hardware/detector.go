package hardware

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CPUInfo contém informações sobre a CPU
type CPUInfo struct {
	Model     string
	Cores     int32
	Frequency float64
	CacheSize int32
}

// GPUInfo contém informações sobre uma GPU
type GPUInfo struct {
	Name      string
	Memory    uint64
	HasCUDA   bool
	HasOpenCL bool
}

// MemoryInfo contém informações sobre a memória
type MemoryInfo struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Available uint64
}

// Detector é responsável por detectar e monitorar recursos de hardware
type Detector struct {
	cpuUsage    prometheus.Gauge
	memoryUsage prometheus.Gauge
	memoryTotal prometheus.Gauge
}

// NewDetector cria um novo detector de hardware
func NewDetector() (*Detector, error) {
	d := &Detector{
		cpuUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Current CPU usage in percent",
		}),
		memoryUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		}),
		memoryTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "memory_total_bytes",
			Help: "Total memory available in bytes",
		}),
	}

	// Initialize memory metrics
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get memory info")
	}
	d.memoryTotal.Set(float64(v.Total))

	return d, nil
}

// Start inicia o monitoramento dos recursos de hardware
func (d *Detector) Start() error {
	// Update metrics immediately
	if err := d.updateMetrics(); err != nil {
		return err
	}

	// Start monitoring in background
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := d.updateMetrics(); err != nil {
				fmt.Printf("Error updating metrics: %v\n", err)
			}
		}
	}()

	return nil
}

// updateMetrics atualiza todas as métricas de hardware
func (d *Detector) updateMetrics() error {
	// Update CPU usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return errors.Wrap(err, "failed to get CPU usage")
	}
	if len(cpuPercent) > 0 {
		d.cpuUsage.Set(cpuPercent[0])
	}

	// Update memory usage
	v, err := mem.VirtualMemory()
	if err != nil {
		return errors.Wrap(err, "failed to get memory info")
	}
	d.memoryUsage.Set(float64(v.Used))

	return nil
}

// GetCPUInfo retorna informações sobre a CPU
func (d *Detector) GetCPUInfo() (*CPUInfo, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CPU info")
	}

	if len(info) == 0 {
		return nil, errors.New("no CPU information available")
	}

	return &CPUInfo{
		Model:     info[0].ModelName,
		Cores:     info[0].Cores,
		Frequency: info[0].Mhz,
		CacheSize: int32(info[0].CacheSize),
	}, nil
}

// HasCUDAGPU verifica se há uma GPU CUDA disponível
func (d *Detector) HasCUDAGPU() bool {
	// TODO: Implement actual CUDA detection
	return false
}

// HasIntelGPU verifica se há uma GPU Intel disponível
func (d *Detector) HasIntelGPU() bool {
	// TODO: Implement actual Intel GPU detection
	return false
}

// GetGPUInfo retorna informações sobre as GPUs disponíveis
func (d *Detector) GetGPUInfo() []GPUInfo {
	// TODO: Implement actual GPU detection
	return []GPUInfo{}
}

// GetMemoryInfo retorna informações sobre a memória
func (d *Detector) GetMemoryInfo() (*MemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get memory info")
	}

	return &MemoryInfo{
		Total:     v.Total,
		Used:      v.Used,
		Free:      v.Free,
		Available: v.Available,
	}, nil
}

// GetHardwareInfo retorna informações gerais sobre o hardware
func (d *Detector) GetHardwareInfo() (map[string]interface{}, error) {
	cpuInfo, err := d.GetCPUInfo()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CPU info")
	}

	memInfo, err := d.GetMemoryInfo()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get memory info")
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CPU usage")
	}

	info := map[string]interface{}{
		"cpu": map[string]interface{}{
			"model":      cpuInfo.Model,
			"cores":      cpuInfo.Cores,
			"usage":      cpuPercent[0],
			"mhz":        cpuInfo.Frequency,
			"cache_size": cpuInfo.CacheSize,
		},
		"memory": map[string]interface{}{
			"total":     memInfo.Total,
			"used":      memInfo.Used,
			"free":      memInfo.Free,
			"available": memInfo.Available,
		},
	}

	return info, nil
}
