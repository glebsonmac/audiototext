# Sistema de Gerenciamento de Modelos

## Descrição
O Sistema de Gerenciamento de Modelos é responsável por todas as operações relacionadas aos modelos de machine learning utilizados na aplicação, incluindo download, armazenamento, conversão, otimização e gerenciamento de versões. Este componente é fundamental para garantir que os modelos adequados estejam disponíveis para inferência e que sejam utilizados de forma eficiente.

## Componentes Principais

### Manager
O componente `Manager` implementa a interface principal para gerenciamento de modelos:
- Download assíncrono de modelos do Hugging Face
- Cache local de modelos para uso offline
- Conversão de modelos para formato ONNX
- Otimização de modelos para diferentes dispositivos (CPU/GPU)
- Sistema de versionamento e backup

## Implementação

### Interface do Gerenciador de Modelos
```go
// Manager define a interface para operações de gerenciamento de modelos
type Manager interface {
    // Gerenciamento de Modelos
    GetModel(id string) error
    ListModels() ([]string, error)
    AddModel(id string) error
    RemoveModel(id string) error
    UpdateModel(id string) error

    // Operações de Modelos
    DownloadModel(ctx context.Context, modelID string) error
    ModelExists(modelID string) bool
    ConvertToONNX(ctx context.Context, modelID string) (string, error)
    ONNXModelExists(modelID string) bool
    ValidateONNXModel(ctx context.Context, modelID string) error
    OptimizeONNXModel(ctx context.Context, modelID string) (string, error)
    OptimizedModelExists(modelID string) bool
    GetModelVersion(modelID string) (string, error)

    // Operações de Backup
    CreateBackup(ctx context.Context, modelID string) error
    ListBackups(modelID string) ([]string, error)
    RestoreFromBackup(ctx context.Context, modelID string, backupID string) error

    // Operações de Cache
    GetCacheSize() (int64, error)
    ClearCache() error
}
```

### Implementação do Gerenciador
```go
// manager implementa a interface Manager
type manager struct {
    modelDir string
    cacheDir string
    mu       sync.RWMutex
    versions map[string]string
    cache    map[string][]byte
}

// NewManager cria um novo gerenciador de modelos
func NewManager(modelDir, cacheDir string) (Manager, error)
```

## Fluxos de Trabalho

### Download e Conversão de Modelos
1. O usuário solicita um modelo específico
2. O sistema verifica se o modelo já existe localmente
3. Se não existir, o sistema inicia o download do Hugging Face
4. Após o download, o modelo é convertido para formato ONNX
5. O modelo convertido é otimizado para o hardware disponível
6. O modelo otimizado é armazenado no cache local

### Gerenciamento de Versões
1. Cada modelo é associado a uma versão específica
2. O sistema mantém um registro das versões disponíveis
3. O usuário pode selecionar a versão desejada
4. O sistema pode criar backups de versões anteriores
5. É possível restaurar modelos a partir de backups

### Otimização de Modelos
1. O sistema detecta o hardware disponível
2. O modelo é otimizado para o hardware específico
3. Otimizações diferentes são aplicadas para CPU e GPU
4. Os parâmetros de otimização são ajustados automaticamente

## Casos de Uso

1. **Inicialização da Aplicação**:
   - Verificação de modelos pré-instalados
   - Download automático de modelos necessários

2. **Atualização de Modelos**:
   - Verificação de novas versões disponíveis
   - Download e conversão de novas versões
   - Criação de backup da versão anterior

3. **Seleção de Modelos**:
   - Listagem de modelos disponíveis
   - Seleção do modelo adequado para a tarefa
   - Carregamento do modelo otimizado

4. **Gerenciamento de Espaço**:
   - Verificação do tamanho do cache
   - Limpeza automática de modelos não utilizados
   - Priorização de modelos frequentemente utilizados

## Status da Implementação
- [x] Implementação do cliente Hugging Face
- [x] Sistema de download assíncrono
- [x] Gerenciamento de cache local
- [x] Sistema de versionamento
- [x] Conversor para ONNX
- [x] Sistema de otimização de modelos
- [x] Validação de modelos convertidos
- [x] Sistema de backup 