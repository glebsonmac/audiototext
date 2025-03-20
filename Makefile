.PHONY: all build test clean run lint wasm proto

# Variables
BINARY_NAME=audiototext
WASM_BINARY=web/wasm/main.wasm
PROTO_DIR=proto

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

all: lint test build

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/server

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f $(WASM_BINARY)

run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/server
	./bin/$(BINARY_NAME)

lint:
	$(GOLINT) run

wasm:
	GOOS=js GOARCH=wasm $(GOBUILD) -o $(WASM_BINARY) ./cmd/wasm

proto:
	protoc --go_out=. --go-grpc_out=. $(PROTO_DIR)/*.proto

deps:
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy

# Development tools installation
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest 