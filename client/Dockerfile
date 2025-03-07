# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache protobuf-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p pkg/pb/audio && \
    protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
    --proto_path=./api/proto/audio \
    ./api/proto/audio/audio.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/client

# Runtime stage (scratch image)
FROM scratch
COPY --from=builder /client /client
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/client"]