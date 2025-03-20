package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CPU Metrics
	cpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "Current CPU usage percentage",
	}, []string{"core"})

	cpuTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current CPU temperature in Celsius",
	})

	// Memory Metrics
	memoryUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_bytes",
		Help: "Current memory usage in bytes",
	})

	memoryAvailable = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "memory_available_bytes",
		Help: "Available memory in bytes",
	})

	// GPU Metrics
	gpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gpu_usage_percent",
		Help: "Current GPU usage percentage",
	})

	gpuMemoryUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gpu_memory_usage_bytes",
		Help: "Current GPU memory usage in bytes",
	})

	// Performance Metrics
	inferenceLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "inference_latency_seconds",
		Help:    "Time taken for model inference",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	})

	batchSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "current_batch_size",
		Help: "Current batch size for inference",
	})
)

// UpdateCPUMetrics atualiza as métricas de CPU
func UpdateCPUMetrics(coreUsage map[int]float64, temperature float64) {
	for core, usage := range coreUsage {
		cpuUsage.WithLabelValues(string(core)).Set(usage)
	}
	cpuTemperature.Set(temperature)
}

// UpdateMemoryMetrics atualiza as métricas de memória
func UpdateMemoryMetrics(used, available uint64) {
	memoryUsage.Set(float64(used))
	memoryAvailable.Set(float64(available))
}

// UpdateGPUMetrics atualiza as métricas de GPU
func UpdateGPUMetrics(usage float64, memoryUsed uint64) {
	gpuUsage.Set(usage)
	gpuMemoryUsage.Set(float64(memoryUsed))
}

// RecordInferenceLatency registra a latência de inferência
func RecordInferenceLatency(seconds float64) {
	inferenceLatency.Observe(seconds)
}

// UpdateBatchSize atualiza o tamanho do batch atual
func UpdateBatchSize(size int) {
	batchSize.Set(float64(size))
}
