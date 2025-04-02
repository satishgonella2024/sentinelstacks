# Makefile for SentinelStacks

# Variables
BINARY_NAME=sentinel
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean test lint vet fmt deps install uninstall run stack-test

all: clean test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/sentinel

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@golint ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Install binary
install: build
	@echo "Installing binary..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/

# Uninstall binary
uninstall:
	@echo "Uninstalling binary..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Test stack functionality
stack-test: build
	@echo "Running stack tests..."
	@chmod +x ./scripts/test-stack.sh
	@./scripts/test-stack.sh
