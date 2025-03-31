.PHONY: build test lint clean

# Variables
BINARY_NAME=sentinel
GO_FILES=$(shell find . -name "*.go" -type f)
VERSION=$(shell git describe --tags --always --dirty || echo "dev")
LDFLAGS=-ldflags "-X github.com/sentinelstacks/sentinel/internal/config.Version=$(VERSION)"

# Main commands
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./... -v

lint:
	@echo "Running linters..."
	@golangci-lint run

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Development helpers
run: build
	@./bin/$(BINARY_NAME)

dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t sentinelstacks/sentinel:$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	@docker run -it --rm sentinelstacks/sentinel:$(VERSION)
