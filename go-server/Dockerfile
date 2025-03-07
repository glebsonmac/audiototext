# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    protobuf-dev \
    git \
    make

WORKDIR /app

# Copiar todos os arquivos do projeto
COPY . .

# Gerar código gRPC
RUN mkdir -p pkg/pb/audio && \
    protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
    --proto_path=./api/proto/audio \
    ./api/proto/audio/audio.proto

# Build do server
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Build do client
RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./client

# Runtime stage
FROM alpine:3.18

# Copiar binários
COPY --from=builder /server /server
COPY --from=builder /client /app/bin/client

# Copiar web assets
COPY ./web /web

# Portas
EXPOSE 50051 8080

# Comando principal
CMD ["/server"]