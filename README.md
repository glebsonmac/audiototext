Aqui está um arquivo `README.md` bilíngue explicando o projeto:

```markdown
# Audio Transcription gRPC Server / Servidor gRPC de Transcrição de Áudio

![Architecture Diagram](https://via.placeholder.com/800x400.png?text=Microservices+Architecture) <!-- Add real diagram if available -->

## 🇧🇷 Português

### Visão Geral
Servidor gRPC em Go que recebe arquivos de áudio e utiliza o Whisper (OpenAI) para transcrição em texto. O projeto segue uma arquitetura de microsserviços com comunicação gRPC.

### Funcionalidades Principais
- 🎙️ Recebimento de áudios via gRPC
- 🔄 Suporte a múltiplos formatos (WAV, MP3, FLAC)
- ⚡ Processamento assíncrono
- 🔒 Comunicação segura entre serviços

### Arquitetura
```
.
├── docker-compose.yml
├── go-server/
│   ├── Dockerfile
│   ├── cmd/
│   │   └── server/
│   │   |    └── main.go
│   |   ├── client/
│   │   |   └── main.go
│   ├── web/
│   │   └── static/
│   │       └── index.html
│   │   └── server/
│   │       └── main.go
│   └── pkg/
└── whisper-server/
    └── Dockerfile
```

### Pré-requisitos
- Go 1.19+
- protoc 3.21+
- Serviço Whisper em execução

### Instalação
1. Clonar repositório:
```bash
git clone https://github.com/seuusuario/audio-transcriber.git
cd audio-transcriber
```

2. Instalar dependências:
```bash
go mod download
```

3. Gerar código gRPC:
```bash
protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
  --proto_path=./api/proto/audio \
  ./api/proto/audio/audio.proto
```

### Uso
1. Iniciar servidor:
```bash
go run cmd/server/main.go
```

2. Executar cliente de teste:
```bash
# Criar arquivo de teste
dd if=/dev/urandom of=testaudio.wav bs=1M count=1 status=none

go run cmd/client/main.go
```

### Configuração do Whisper
```bash
# No diretório do serviço Whisper (Python)
em andamento
```

## 🇺🇸 English

### Overview
Go gRPC server that receives audio files and uses Whisper (OpenAI) for text transcription. The project follows a microservices architecture with gRPC communication.

### Key Features
- 🎙️ Audio reception via gRPC
- 🔄 Multiple format support (WAV, MP3, FLAC)
- ⚡ Async processing
- 🔒 Secure service communication

### Architecture
```
.
├── api/              # Protobuf definitions
├── cmd/              # Entrypoints
├── internal/         # Internal logic
├── pkg/              # Generated code
└── client/           # Test client
```

### Prerequisites
- Go 1.19+
- protoc 3.21+
- Running Whisper service

### Installation
1. Clone repository:
```bash
git clone https://github.com/youruser/audio-transcriber.git
cd audio-transcriber
```

2. Install dependencies:
```bash
go mod download
```

3. Generate gRPC code:
```bash
protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
  --proto_path=./api/proto/audio \
  ./api/proto/audio/audio.proto
```

### Usage
1. Start server:
```bash
go run cmd/server/main.go
```

2. Run test client:
```bash
# Create test file
dd if=/dev/urandom of=testaudio.wav bs=1M count=1 status=none

go run cmd/client/main.go
```

### Whisper Setup
```bash
# In Whisper service directory (Python)
on going
```

## 📝 License / Licença
MIT License © 2024 [Seu Nome]  
MIT License © 2024 [Your Name]

```

Este README fornece:
1. Explicação técnica em ambos os idiomas
2. Instruções de instalação e uso
3. Diagrama arquitetural (placeholder)
4. Seções claramente separadas para cada idioma
5. Informações essenciais para desenvolvedores

Você pode personalizar com:
- Informações reais de contato
- Diagramas específicos
- Detalhes adicionais de configuração
- Logos ou badges do projeto