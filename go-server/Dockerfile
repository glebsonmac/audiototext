# Build stage
FROM golang:1.23.4-alpine AS builder

# Configurar ambiente Go
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct

# Instalar dependências de sistema
RUN apk add --no-cache \
    protobuf-dev \
    git \
    make \
    curl

# Criar diretório de trabalho
WORKDIR /app

# 1. Copiar apenas os arquivos de módulo primeiro
COPY go.mod go.sum ./

# 2. Validar estrutura do módulo
RUN test -f go.mod && echo "Módulo Go válido" || (echo "ERRO: go.mod não encontrado" && exit 1)

# 3. Baixar dependências com verificação
RUN go mod download -x 2>&1 | tee /tmp/mod-download.log

# 4. Instalar protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# 5. Copiar o resto do código
COPY . .

# Build de todos componentes

#Build dos componentes com caminhos absolutos
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/web-server ./cmd/web/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/client ./cmd/client

# Runtime stage
FROM alpine:3.18

# 1. Criar diretório para os binários
RUN mkdir -p /app/bin

# 2. Copiar binários do estágio builder
COPY --from=builder /app/bin/ /app/bin/

# 3. Dar permissão de execução
RUN chmod +x /app/bin/*

# 4. Expor portas
EXPOSE 50051 8080

# 5. Definir comando de execução com caminho absoluto
CMD ["/app/bin/server"]