# Sistema gRPC

## Descrição
O Sistema gRPC é responsável pela comunicação entre clientes e o serviço de transcrição de áudio, fornecendo uma interface eficiente e de alto desempenho para streaming de áudio e retorno de transcrições. Utilizando Protocol Buffers e HTTP/2, o sistema oferece streaming bidirecional, compressão eficiente de dados, e suporte para múltiplos clientes simultâneos.

## Componentes Principais

### Protobuf e Geração de Código
O sistema utiliza Protocol Buffers para definição e serialização de mensagens:
- Definição de schemas protobuf para mensagens e serviços
- Geração automática de código para cliente e servidor
- Versionamento de APIs
- Documentação integrada

### Servidor gRPC
O componente de servidor gerencia as conexões e processa as requisições:
- Implementação do servidor gRPC
- Suporte a streaming bidirecional
- Sistema de reconexão automática
- Balanceamento de carga
- Segurança e autenticação

## Implementação

### Definição do Serviço
```protobuf
service TranscriptionService {
  // Transcrição de áudio em streaming
  rpc TranscribeAudioStream (stream AudioChunk) returns (stream TranscriptionResult);
  
  // Transcrição de áudio em batch
  rpc TranscribeAudio (AudioData) returns (TranscriptionResult);
  
  // Listar modelos disponíveis
  rpc ListModels (ListModelsRequest) returns (ListModelsResponse);
  
  // Status do serviço
  rpc GetServiceStatus (ServiceStatusRequest) returns (ServiceStatusResponse);
}
```

### Interface do Servidor
```go
// Server define a interface para o servidor gRPC
type Server interface {
    // Inicialização e gerenciamento
    Start() error
    Stop() error
    
    // Status e monitoramento
    GetStatus() (*Status, error)
    
    // Registro de serviços
    RegisterService(service interface{})
}
```

### Implementação do Servidor
```go
// server implementa a interface Server
type server struct {
    grpcServer *grpc.Server
    addr       string
    config     *Config
    services   []interface{}
    status     *Status
    mu         sync.RWMutex
}

// NewServer cria um novo servidor gRPC
func NewServer(addr string, options ...Option) (Server, error)
```

## Fluxos de Comunicação

### Streaming de Áudio
1. O cliente inicia um stream bidirecional com o servidor
2. Chunks de áudio são enviados continuamente para o servidor
3. O servidor processa os chunks utilizando o sistema de inferência
4. Resultados parciais são retornados ao cliente em tempo real
5. O cliente pode encerrar o stream a qualquer momento

### Monitoramento de Status
1. O cliente solicita o status atual do serviço
2. O servidor coleta informações de todos os subsistemas
3. Um relatório completo é retornado ao cliente
4. O cliente pode configurar alertas baseados no status

## Casos de Uso

1. **Aplicação Web com WebAssembly**:
   - Frontend WebAssembly conecta-se ao servidor gRPC
   - Streaming de áudio capturado pelo navegador
   - Exibição de resultados em tempo real

2. **Aplicações Desktop/Mobile**:
   - Clientes nativos conectam-se ao servidor
   - Suporte a funcionalidades offline quando necessário
   - Sincronização quando a conexão é restaurada

3. **Integração com Sistemas Externos**:
   - APIs gRPC para integração com outros serviços
   - Autenticação e autorização
   - Suporte a múltiplos protocolos

4. **Escalabilidade Horizontal**:
   - Balanceamento de carga entre múltiplos servidores
   - Descoberta automática de serviços
   - Resiliência a falhas

## Otimizações e Características

- **Performance**:
  - Serialização eficiente com Protocol Buffers
  - Compressão de dados
  - Multiplexação com HTTP/2

- **Resiliência**:
  - Reconexão automática
  - Timeout e retry configuráveis
  - Graceful shutdown

- **Observabilidade**:
  - Métricas de performance
  - Tracing distribuído
  - Logging estruturado

## Status da Implementação
- [x] Definição de schemas protobuf
- [x] Configuração de geração de código
- [x] Implementação de interfaces
- [x] Criação de testes de API
- [x] Implementação do servidor gRPC
- [x] Configuração de streaming
- [x] Implementação de reconexão
- [x] Sistema de load balancing 