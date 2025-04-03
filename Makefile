# Makefile for SentinelStacks

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=sentinel
API_EXAMPLE=api-example

# Main binary path
CMD_DIR=cmd/sentinel
API_EXAMPLE_DIR=cmd/api_example

# Output directories
BIN_DIR=bin
DATA_DIR=data

# Default target
all: clean build

# Build the main binary
build:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(API_EXAMPLE) ./$(API_EXAMPLE_DIR)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Run the main application
run: build
	$(BIN_DIR)/$(BINARY_NAME)

# Run the API example
run-example: build
	$(BIN_DIR)/$(API_EXAMPLE)

# Run specific examples
run-example-stack: build
	$(BIN_DIR)/$(API_EXAMPLE) stack

run-example-memory: build
	$(BIN_DIR)/$(API_EXAMPLE) memory

run-example-registry: build
	$(BIN_DIR)/$(API_EXAMPLE) registry

run-example-comprehensive: build
	$(BIN_DIR)/$(API_EXAMPLE) comprehensive

# Initialize the application with example data
init: build
	mkdir -p $(DATA_DIR)
	$(BIN_DIR)/$(API_EXAMPLE) init

# Verify dependencies
deps:
	$(GOMOD) verify
	$(GOMOD) tidy

# Update dependencies
deps-update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Install the application
install: build
	cp $(BIN_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Build for release (all platforms)
release:
	mkdir -p $(BIN_DIR)/release
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/release/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/release/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_DIR)/release/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

.PHONY: all build clean test test-coverage run run-example run-example-stack run-example-memory run-example-registry run-example-comprehensive init deps deps-update install release
