package detection

// CPUInfo contém informações sobre a CPU do sistema
type CPUInfo struct {
	PhysicalCores int     `json:"physical_cores"`
	LogicalCores  int     `json:"logical_cores"`
	BaseFrequency float64 `json:"base_frequency_mhz"`
	MaxFrequency  float64 `json:"max_frequency_mhz"`
	CacheL1       int     `json:"cache_l1_kb"`
	CacheL2       int     `json:"cache_l2_kb"`
	CacheL3       int     `json:"cache_l3_kb"`
}

// GPUInfo contém informações sobre a GPU do sistema
type GPUInfo struct {
	HasCUDA           bool   `json:"has_cuda"`
	HasOpenCL         bool   `json:"has_opencl"`
	TotalMemory       uint64 `json:"total_memory_mb"`
	ComputeCapability string `json:"compute_capability"`
	Name              string `json:"name"`
}

// MemoryInfo contém informações sobre a memória do sistema
type MemoryInfo struct {
	TotalRAM     uint64 `json:"total_ram_mb"`
	AvailableRAM uint64 `json:"available_ram_mb"`
	SwapTotal    uint64 `json:"swap_total_mb"`
	SwapUsed     uint64 `json:"swap_used_mb"`
}

// HardwareProfile define o perfil de hardware para otimização
type HardwareProfile string

const (
	HighEnd  HardwareProfile = "high_end"
	MidRange HardwareProfile = "mid_range"
	LowEnd   HardwareProfile = "low_end"
)

// HardwareInfo agrega todas as informações de hardware
type HardwareInfo struct {
	CPU     CPUInfo         `json:"cpu"`
	GPU     GPUInfo         `json:"gpu"`
	Memory  MemoryInfo      `json:"memory"`
	Profile HardwareProfile `json:"profile"`
}
