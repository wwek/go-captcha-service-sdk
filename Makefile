# Makefile for go-captcha-service-sdk project

# Variables
BINARY_NAME=go-captcha-service-sdk
GO=go

# Default Target
.PHONY: all
all: build

# Install golang Dependencies
.PHONY: deps-go
deps-go:
	$(GO) mod tidy
	$(GO) mod download
	@if ! command -v protoc >/dev/null; then \
		echo "Installing protoc..."; \
		$(GO) install github.com/golang/protobuf/protoc-gen-go@latest; \
	fi
	@if ! command -v protobufjs >/dev/null; then \
		echo "Installing protobufjs..."; \
		npm install -g protobufjs \
	fi

# Generate gRPC code for golang
.PHONY: proto-go
proto-go:
	protoc --go_out=./golang --go-grpc_out=./golang ./gocaptcha-service-api.proto

# Run tests
.PHONY: test
test: proto
	$(GO) test -v ./...

# Coverage report
.PHONY: cover
cover: proto
	$(GO) test -cover -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	$(GO) fmt ./...

# Help Information
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps-go            : Install golang dependencies"
	@echo "  proto-go           : Generate Protobuf code for golang"
	@echo "  test               : Run tests"
	@echo "  cover              : Generate test coverage report"
	@echo "  clean              : Remove build artifacts"
	@echo "  help               : Show this help message"