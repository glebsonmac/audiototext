package detection

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Detector implementa a detecção de hardware
type Detector struct {
	lastInfo HardwareInfo
}

// NewDetector cria uma nova instância do detector
func NewDetector() *Detector {
	return &Detector{}
}

// Detect realiza a detecção de hardware
func (d *Detector) Detect() (HardwareInfo, error) {
	info := HardwareInfo{}

	// Detecta CPU
	cpuInfo, err := d.detectCPU()
	if err != nil {
		return info, fmt.Errorf("erro ao detectar CPU: %w", err)
	}
	info.CPU = cpuInfo

	// Detecta GPU
	gpuInfo, err := d.detectGPU()
	if err != nil {
		return info, fmt.Errorf("erro ao detectar GPU: %w", err)
	}
	info.GPU = gpuInfo

	// Detecta Memória
	memInfo, err := d.detectMemory()
	if err != nil {
		return info, fmt.Errorf("erro ao detectar memória: %w", err)
	}
	info.Memory = memInfo

	// Determina o perfil
	info.Profile = d.determineProfile(info)

	d.lastInfo = info
	return info, nil
}

// detectCPU obtém informações da CPU
func (d *Detector) detectCPU() (CPUInfo, error) {
	info := CPUInfo{}

	cpus, err := cpu.Info()
	if err != nil {
		return info, err
	}

	if len(cpus) > 0 {
		info.PhysicalCores = runtime.NumCPU() / 2 // Aproximação para cores físicos
		info.LogicalCores = runtime.NumCPU()
		info.BaseFrequency = float64(cpus[0].Mhz)
	}

	return info, nil
}

// detectGPU obtém informações da GPU
func (d *Detector) detectGPU() (GPUInfo, error) {
	info := GPUInfo{}

	// TODO: Implementar detecção real de CUDA e OpenCL
	// Por enquanto, apenas detectamos se estamos em um sistema que pode ter GPU
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		info.HasCUDA = true // Simplificação
		info.HasOpenCL = true
	}

	return info, nil
}

// detectMemory obtém informações da memória
func (d *Detector) detectMemory() (MemoryInfo, error) {
	info := MemoryInfo{}

	v, err := mem.VirtualMemory()
	if err != nil {
		return info, err
	}

	info.TotalRAM = v.Total / 1024 / 1024 // Converte para MB
	info.AvailableRAM = v.Available / 1024 / 1024

	swap, err := mem.SwapMemory()
	if err != nil {
		return info, err
	}

	info.SwapTotal = swap.Total / 1024 / 1024
	info.SwapUsed = swap.Used / 1024 / 1024

	return info, nil
}

// determineProfile determina o perfil de hardware
func (d *Detector) determineProfile(info HardwareInfo) HardwareProfile {
	// Critérios para High-End:
	// - 8+ cores lógicos
	// - 16GB+ RAM
	// - GPU disponível
	if info.CPU.LogicalCores >= 8 &&
		info.Memory.TotalRAM >= 16*1024 && // 16GB
		(info.GPU.HasCUDA || info.GPU.HasOpenCL) {
		return HighEnd
	}

	// Critérios para Mid-Range:
	// - 4+ cores lógicos
	// - 8GB+ RAM
	if info.CPU.LogicalCores >= 4 &&
		info.Memory.TotalRAM >= 8*1024 { // 8GB
		return MidRange
	}

	return LowEnd
}

// GetProfile retorna o último perfil detectado
func (d *Detector) GetProfile() (HardwareProfile, error) {
	if d.lastInfo.Profile == "" {
		_, err := d.Detect()
		if err != nil {
			return "", err
		}
	}
	return d.lastInfo.Profile, nil
}
