package hardware

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
)

// monitorMetrics holds Prometheus metrics
var (
	monitorCPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "Current CPU usage in percent",
	})

	monitorMemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_percent",
		Help: "Current memory usage in percent",
	})

	monitorGPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gpu_usage_percent",
		Help: "Current GPU usage in percent",
	})
)

func init() {
	prometheus.MustRegister(monitorCPUUsage)
	prometheus.MustRegister(monitorMemoryUsage)
	prometheus.MustRegister(monitorGPUUsage)
}

// monitor implements the ResourceMonitor interface
type monitor struct {
	detector *Detector
	stopCh   chan struct{}
	metrics  ResourceMetrics
	mu       sync.RWMutex
}

// NewMonitor creates a new resource monitor
func NewMonitor(detector *Detector) (ResourceMonitor, error) {
	return &monitor{
		detector: detector,
		stopCh:   make(chan struct{}),
	}, nil
}

// Start begins monitoring system resources
func (m *monitor) Start() error {
	go m.monitorLoop()
	return nil
}

// Stop stops monitoring system resources
func (m *monitor) Stop() error {
	close(m.stopCh)
	return nil
}

// GetMetrics returns the current resource metrics
func (m *monitor) GetMetrics() ResourceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics
}

// monitorLoop periodically updates resource metrics
func (m *monitor) monitorLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.updateMetrics()
		}
	}
}

// updateMetrics updates the current resource metrics
func (m *monitor) updateMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update CPU metrics
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		m.metrics.CPUUsage = cpuPercent[0]
		monitorCPUUsage.Set(m.metrics.CPUUsage)
	}

	// Update memory metrics
	memInfo, err := m.detector.GetMemoryInfo()
	if err == nil {
		m.metrics.MemoryUsage = memInfo.UsedPercent
		monitorMemoryUsage.Set(m.metrics.MemoryUsage)
	}

	// GPU metrics would be implemented here if needed
	m.metrics.GPUUsage = 0.0
	monitorGPUUsage.Set(m.metrics.GPUUsage)
}
