# Sistema de Detecção de Hardware

## Descrição
O Sistema de Detecção de Hardware é responsável por identificar e fornecer informações sobre o hardware disponível no sistema, incluindo CPU, memória e futuras integrações para GPU. Essa funcionalidade é essencial para otimizar a execução de modelos de inferência, permitindo que a aplicação adapte-se automaticamente aos recursos disponíveis.

## Componentes Principais

### Detector
O componente `Detector` fornece métodos para obter informações detalhadas do hardware do sistema:
- Detecção de CPU (núcleos, frequência, cache)
- Detecção de memória (total, disponível, uso)
- Interface para futura integração com detecção de GPU (CUDA, OpenCL)

### Monitor de Recursos
O `ResourceMonitor` implementa monitoramento contínuo do uso de recursos:
- Monitoramento de uso de CPU
- Monitoramento de uso de memória
- Métricas em tempo real via Prometheus
- Alertas baseados em limites configuráveis

## Implementação

### Detector de Hardware
```go
// Detector fornece capacidades de detecção de hardware
type Detector struct {
    cpuInfo    []cpu.InfoStat
    memoryInfo *mem.VirtualMemoryStat
}

// NewDetector cria um novo detector de hardware
func NewDetector() (*Detector, error)

// GetCPUInfo retorna informações da CPU
func (d *Detector) GetCPUInfo() ([]cpu.InfoStat, error)

// GetMemoryInfo retorna informações da memória
func (d *Detector) GetMemoryInfo() (*mem.VirtualMemoryStat, error)
```

### Monitor de Recursos
```go
// ResourceMetrics contém métricas de uso de recursos
type ResourceMetrics struct {
    CPUUsage    float64
    MemoryUsage float64
    GPUUsage    float64
}

// ResourceMonitor fornece capacidades de monitoramento
type ResourceMonitor interface {
    Start() error
    Stop() error
    GetMetrics() ResourceMetrics
}
```

## Métricas e Integrações

### Prometheus
O sistema integra-se com Prometheus para exportar métricas em tempo real:
- `cpu_usage_percent`: Uso de CPU em percentual
- `memory_usage_percent`: Uso de memória em percentual
- `gpu_usage_percent`: Uso de GPU em percentual (preparado para implementação futura)

## Casos de Uso

1. **Inicialização do Sistema**:
   - Detecção automática do hardware disponível
   - Configuração inicial baseada nos recursos detectados

2. **Otimização de Execução**:
   - Ajuste de parâmetros de modelos com base nos recursos disponíveis
   - Seleção do dispositivo de inferência mais adequado (CPU/GPU)

3. **Monitoramento em Tempo Real**:
   - Visualização do uso de recursos durante a execução
   - Detecção de gargalos e limitações de hardware

4. **Alertas e Notificações**:
   - Alertas quando o uso de recursos excede limites configurados
   - Notificações de hardware inadequado para determinadas operações

## Status da Implementação
- [x] Implementação do detector de CPU
- [x] Implementação do monitor de recursos
- [x] Implementação das métricas Prometheus
- [x] Integração com o sistema de inferência 