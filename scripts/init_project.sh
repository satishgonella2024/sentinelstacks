#!/bin/bash

# SentinelStacks Project Initialization Script
# This script sets up the initial project structure for SentinelStacks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print banner
echo -e "${BLUE}"
echo "  _____            _   _            _  _____ _             _        "
echo " / ____|          | | (_)          | |/ ____| |           | |       "
echo "| (___   ___ _ __ | |_ _ _ __   ___| | (___ | |_ __ _  ___| | _____ "
echo " \___ \ / _ \ '_ \| __| | '_ \ / _ \ |\___ \| __/ _\` |/ __| |/ / __|"
echo " ____) |  __/ | | | |_| | | | |  __/ |____) | || (_| | (__|   <\__ \\"
echo "|_____/ \___|_| |_|\__|_|_| |_|\___|_|_____/ \__\__,_|\___|_|\_\___/"
echo -e "${NC}"
echo -e "${YELLOW}AI Agent Management System${NC}"
echo

# Check for dependencies
echo -e "${BLUE}Checking dependencies...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.20 or later.${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "- Go version: ${GREEN}${GO_VERSION}${NC}"

# Check for Git
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git is not installed. Please install Git.${NC}"
    exit 1
fi

GIT_VERSION=$(git --version | awk '{print $3}')
echo -e "- Git version: ${GREEN}${GIT_VERSION}${NC}"

# Create directory structure
echo -e "\n${BLUE}Creating project structure...${NC}"

# Main directories
mkdir -p cmd/sentinel
mkdir -p internal/{config,parser,runtime,registry,shim}
mkdir -p pkg/{agent,tools,api}
mkdir -p cmd/sentinel/{build,run,push,pull,config,init}

# Documentation
echo -e "- Documentation"
mkdir -p docs/{architecture,planning,user-guides,visualizations}

# Configuration
echo -e "- Configuration"
mkdir -p configs

# Scripts and tools
echo -e "- Scripts and tools"
mkdir -p scripts/tools
mkdir -p tools

# Examples
echo -e "- Examples"
mkdir -p examples

# Testing
echo -e "- Testing"
mkdir -p test/{unit,integration,e2e}

# Creating initial files
echo -e "\n${BLUE}Creating initial files...${NC}"

# Root main.go
cat > main.go << EOF
package main

import (
	"fmt"
	"os"

	"github.com/sentinelstacks/sentinel/cmd/sentinel"
)

func main() {
	if err := sentinel.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
EOF
echo -e "- Created main.go"

# Root mod file
cat > go.mod << EOF
module github.com/sentinelstacks/sentinel

go 1.20

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
)
EOF
echo -e "- Created go.mod"

# Root command file
mkdir -p cmd/sentinel
cat > cmd/sentinel/root.go << EOF
package sentinel

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "SentinelStacks - AI Agent Management System",
	Long: \`SentinelStacks is a comprehensive system for creating, managing,
and distributing AI agents using natural language definitions.

It provides a Docker-like workflow for AI agents:\`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	
	// Add commands
	// These will be implemented in their respective packages
}
EOF
echo -e "- Created cmd/sentinel/root.go"

# Makefile
cat > Makefile << EOF
.PHONY: build test lint clean

# Variables
BINARY_NAME=sentinel
GO_FILES=\$(shell find . -name "*.go" -type f)
VERSION=\$(shell git describe --tags --always --dirty || echo "dev")
LDFLAGS=-ldflags "-X github.com/sentinelstacks/sentinel/internal/config.Version=\$(VERSION)"

# Main commands
build:
	@echo "Building \$(BINARY_NAME)..."
	@go build \$(LDFLAGS) -o bin/\$(BINARY_NAME) .

install: build
	@echo "Installing \$(BINARY_NAME)..."
	@cp bin/\$(BINARY_NAME) \$(GOPATH)/bin/\$(BINARY_NAME)

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
	@./bin/\$(BINARY_NAME)

dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t sentinelstacks/sentinel:\$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	@docker run -it --rm sentinelstacks/sentinel:\$(VERSION)
EOF
echo -e "- Created Makefile"

# Git initialization
echo -e "\n${BLUE}Initializing Git repository...${NC}"
git init
git add .

# Create .gitignore
cat > .gitignore << EOF
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Build directories
dist/

# IDE directories
.idea/
.vscode/
*.sublime-workspace

# OS specific files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Config files that might contain secrets
.env
*.env
config.local.yaml
EOF
echo -e "- Created .gitignore"

# README.md
cat > README.md << EOF
# SentinelStacks

![SentinelStacks Logo](docs/visualizations/sentinelstacks_logo.svg)

## AI Agent Management System

SentinelStacks is a comprehensive system for creating, managing, and distributing AI agents using natural language. Built with inspiration from Docker's paradigm, SentinelStacks allows you to:

- Define agents using natural language in Sentinelfiles
- Build agent images that can be versioned and shared
- Run agents across different LLM backends (Claude, OpenAI, Llama, etc.)
- Manage agent state and orchestrate multi-agent systems
- Share agents through public and private registries

## Installation

\`\`\`bash
# Install SentinelStacks
go install github.com/sentinelstacks/sentinel@latest

# Verify installation
sentinel version
\`\`\`

## Documentation

- [Architecture Overview](docs/architecture/README.md)
- [User Guides](docs/user-guides/README.md)
- [Development Roadmap](docs/planning/roadmap.md)
- [API Reference](docs/architecture/api.md)

## License

MIT
EOF
echo -e "- Created README.md"

# Create Docker file
cat > Dockerfile << EOF
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /sentinel .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /sentinel /app/sentinel

ENTRYPOINT ["/app/sentinel"]
EOF
echo -e "- Created Dockerfile"

echo -e "\n${GREEN}Project initialization complete!${NC}"
echo -e "To build the project, run: ${YELLOW}make build${NC}"
echo -e "To run tests, run: ${YELLOW}make test${NC}"
echo -e "To install the CLI, run: ${YELLOW}make install${NC}"
echo
echo -e "${BLUE}Happy coding!${NC}"
