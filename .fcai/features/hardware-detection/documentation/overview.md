# Hardware Detection Feature

## Overview
O sistema de detecção de hardware é responsável por identificar e monitorar os recursos disponíveis no sistema para otimizar o processamento de áudio e inferência do modelo ONNX.

## Componentes Principais

### 1. Detector de Hardware
- Detecção de CPU (cores, frequência, cache)
- Detecção de GPU (CUDA, OpenCL)
- Análise de memória RAM
- Sistema de perfis de hardware

### 2. Monitor de Recursos
- Monitoramento em tempo real de CPU, GPU e memória
- Sistema de métricas com Prometheus
- Logging de eventos e performance
- Sistema de alertas

## Interfaces

### Hardware Detection API
```go
type HardwareInfo struct {
    CPU     CPUInfo
    GPU     GPUInfo
    Memory  MemoryInfo
    Profile HardwareProfile
}

type HardwareDetector interface {
    Detect() (HardwareInfo, error)
    Monitor() (<-chan HardwareMetrics, error)
    GetProfile() (HardwareProfile, error)
}
```

### Resource Monitor API
```go
type ResourceMonitor interface {
    Start() error
    Stop() error
    GetMetrics() ResourceMetrics
    Subscribe(chan<- ResourceMetrics) error
}
```

## Métricas Coletadas

### CPU
- Número de cores físicos e lógicos
- Frequência base e boost
- Utilização por core
- Temperatura
- Cache L1, L2, L3

### GPU
- Disponibilidade CUDA/OpenCL
- Memória total e disponível
- Utilização
- Temperatura
- Compute capability

### Memória
- RAM total e disponível
- Swap utilizado
- Velocidade de acesso
- Fragmentação

## Integração com ONNX Runtime
O sistema de detecção de hardware é utilizado para:
1. Escolher o backend mais apropriado (CPU/GPU)
2. Otimizar parâmetros de inferência
3. Ajustar batch size dinamicamente
4. Gerenciar recursos de forma eficiente

## Alertas e Logs
- Alertas de sobrecarga de CPU/GPU
- Alertas de memória baixa
- Logs de performance
- Métricas históricas

## Dependências
- github.com/shirou/gopsutil/v3
- github.com/NVIDIA/gpu-monitoring-tools
- github.com/prometheus/client_golang 