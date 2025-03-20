package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// ResourceMetrics contém métricas de recursos do sistema
type ResourceMetrics struct {
	CPUUsage        map[int]float64
	MemoryUsed      uint64
	MemoryAvailable uint64
	GPUUsage        float64
	GPUMemoryUsed   uint64
	Timestamp       time.Time
}

// Monitor implementa o monitoramento de recursos
type Monitor struct {
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	metrics  chan ResourceMetrics
}

// NewMonitor cria um novo monitor de recursos
func NewMonitor(interval time.Duration) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
		metrics:  make(chan ResourceMetrics, 100),
	}
}

// Start inicia o monitoramento
func (m *Monitor) Start() error {
	go m.monitor()
	return nil
}

// Stop para o monitoramento
func (m *Monitor) Stop() error {
	m.cancel()
	close(m.metrics)
	return nil
}

// GetMetrics retorna o canal de métricas
func (m *Monitor) GetMetrics() <-chan ResourceMetrics {
	return m.metrics
}

// monitor é a goroutine principal de monitoramento
func (m *Monitor) monitor() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			metrics, err := m.collectMetrics()
			if err != nil {
				// TODO: Implementar logging adequado
				fmt.Printf("Erro ao coletar métricas: %v\n", err)
				continue
			}

			// Envia métricas de forma não bloqueante
			select {
			case m.metrics <- metrics:
			default:
				// Canal cheio, descarta métricas antigas
				<-m.metrics
				m.metrics <- metrics
			}
		}
	}
}

// collectMetrics coleta métricas do sistema
func (m *Monitor) collectMetrics() (ResourceMetrics, error) {
	metrics := ResourceMetrics{
		CPUUsage:  make(map[int]float64),
		Timestamp: time.Now(),
	}

	// Coleta uso de CPU
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return metrics, fmt.Errorf("erro ao coletar CPU: %w", err)
	}
	for i, usage := range percentages {
		metrics.CPUUsage[i] = usage
	}

	// Coleta uso de memória
	memory, err := mem.VirtualMemory()
	if err != nil {
		return metrics, fmt.Errorf("erro ao coletar memória: %w", err)
	}
	metrics.MemoryUsed = memory.Used
	metrics.MemoryAvailable = memory.Available

	// TODO: Implementar coleta real de métricas de GPU
	// Por enquanto, usamos valores simulados
	metrics.GPUUsage = 0
	metrics.GPUMemoryUsed = 0

	return metrics, nil
}

// Subscribe registra um canal para receber métricas
func (m *Monitor) Subscribe(ch chan<- ResourceMetrics) error {
	go func() {
		for metric := range m.metrics {
			ch <- metric
		}
	}()
	return nil
}
