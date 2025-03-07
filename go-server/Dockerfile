# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    protobuf-dev \
    git \
    make \
    curl

# Configurar módulo Go
ENV GO111MODULE=on

# Instalar dependências do protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH $PATH:/root/go/bin

WORKDIR /app

# Copiar apenas os arquivos necessários primeiro
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código
COPY . .

# Gerar código gRPC
RUN mkdir -p pkg/pb/audio && \
    protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
    --proto_path=./api/proto/audio \
    ./api/proto/audio/audio.proto

# Build dos componentes
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-server ./web/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/client

# Runtime stage
FROM alpine:3.18

COPY --from=builder /server /server
COPY --from=builder /web-server /web-server
COPY --from=builder /client /client
COPY --from=builder /app/web/static /web/static

EXPOSE 50051 8080
CMD ["/server"]