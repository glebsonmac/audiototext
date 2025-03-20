# Sistema de Inferência ONNX

## Descrição
O Sistema de Inferência ONNX é responsável por executar modelos de machine learning já convertidos para o formato ONNX (Open Neural Network Exchange). Este sistema oferece uma camada de abstração para executar inferências de forma eficiente, aproveitando aceleradores de hardware quando disponíveis e fornecendo um pipeline completo de pré-processamento, inferência e pós-processamento.

## Componentes Principais

### Runtime ONNX
O componente de Runtime ONNX integra o ONNX Runtime para execução eficiente de modelos:
- Configuração de sessões de inferência
- Execução em CPU ou GPU (quando disponível)
- Otimização automática baseada no hardware
- Sistema de fallback para garantir execução em qualquer ambiente

### Pipeline de Inferência
O pipeline fornece um fluxo completo para processamento de áudio:
- Pré-processamento de áudio (normalização, segmentação)
- Execução do modelo de inferência
- Pós-processamento de resultados (formatação de texto)
- Otimização de performance

## Implementação

### Interface do Runtime
```go
// Runtime define a interface para execução de modelos ONNX
type Runtime interface {
    // Inicialização e configuração
    Initialize(modelPath string, options ...Option) error
    Close() error
    
    // Execução de inferência
    RunInference(input []float32) ([]float32, error)
    RunBatchInference(inputs [][]float32) ([][]float32, error)
    
    // Configuração e status
    SetExecutionProvider(provider string) error
    GetSupportedProviders() []string
    GetModelMetadata() (map[string]string, error)
}
```

### Implementação do Runtime
```go
// onnxRuntime implementa a interface Runtime
type onnxRuntime struct {
    session      *onnxruntime.Session
    inputNames   []string
    outputNames  []string
    options      *onnxruntime.SessionOptions
    initialized  bool
    deviceInfo   *hardware.DeviceInfo
}

// NewRuntime cria uma nova instância do runtime ONNX
func NewRuntime(deviceInfo *hardware.DeviceInfo) (Runtime, error)
```

## Execução do Pipeline

### Processo de Inferência
1. O áudio é recebido como entrada
2. O sistema aplica pré-processamento específico para o modelo
3. A entrada processada é passada para o runtime ONNX
4. O runtime executa a inferência usando o provedor apropriado (CPU/GPU)
5. Os resultados são pós-processados para gerar a transcrição final
6. A transcrição é retornada para o cliente

### Otimização de Performance
- Uso de batching para inferência em lote
- Ajuste automático de threads baseado em cores de CPU
- Utilização de GPU quando disponível
- Cache de resultados intermediários
- Balanceamento dinâmico de carga

## Casos de Uso

1. **Transcrição de Áudio em Tempo Real**:
   - Segmentação de stream de áudio
   - Inferência contínua em chunks
   - Combinação de resultados parciais

2. **Processamento em Lote**:
   - Inferência em múltiplos arquivos de áudio
   - Otimização para throughput máximo
   - Distribuição de carga em múltiplos dispositivos

3. **Adaptação a Diferentes Ambientes**:
   - Execução em dispositivos de baixa potência (CPU)
   - Execução em servidores com GPUs
   - Fallback automático entre dispositivos

4. **Integração com Sistemas Externos**:
   - Exposição do serviço via gRPC
   - Interface para streaming de áudio
   - Comunicação bidirecional

## Status da Implementação
- [x] Integração do ONNX Runtime
- [x] Configuração de sessões de inferência
- [x] Implementação de batch processing
- [x] Sistema de fallback CPU/GPU
- [x] Implementação do pré-processamento
- [x] Configuração da inferência
- [x] Implementação do pós-processamento
- [x] Sistema de otimização de performance 