package hardware

// ResourceMetrics holds the current resource usage metrics
type ResourceMetrics struct {
	CPUUsage    float64
	MemoryUsage float64
	GPUUsage    float64
}

// ResourceMonitor provides resource monitoring capabilities
type ResourceMonitor interface {
	Start() error
	Stop() error
	GetMetrics() ResourceMetrics
}
