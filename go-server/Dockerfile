# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    protobuf-dev \
    grpc-dev \
    git \
    make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p pkg/pb/audio && \
    protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb \
    --proto_path=./api/proto/audio \
    ./api/proto/audio/audio.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Runtime stage
FROM alpine:3.18
COPY --from=builder /server /server
EXPOSE 50051
CMD ["/server"]