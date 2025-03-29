#!/bin/bash

# Build script for SentinelStacks CLI

VERSION="1.0.0"
BINARY_NAME="sentinel"

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m'

echo "Building SentinelStacks CLI v${VERSION}..."

# Build for macOS (Apple Silicon)
echo -e "${GREEN}Building for macOS (Apple Silicon)...${NC}"
GOOS=darwin GOARCH=arm64 go build -o "dist/darwin-arm64/${BINARY_NAME}" cmd/sentinel/main.go

# Build for macOS (Intel)
echo -e "${GREEN}Building for macOS (Intel)...${NC}"
GOOS=darwin GOARCH=amd64 go build -o "dist/darwin-amd64/${BINARY_NAME}" cmd/sentinel/main.go

# Build for Linux
echo -e "${GREEN}Building for Linux...${NC}"
GOOS=linux GOARCH=amd64 go build -o "dist/linux-amd64/${BINARY_NAME}" cmd/sentinel/main.go

echo -e "${GREEN}Build complete!${NC}"
echo
echo "Binaries are available in:"
echo "  - dist/darwin-arm64/sentinel (macOS Apple Silicon)"
echo "  - dist/darwin-amd64/sentinel (macOS Intel)"
echo "  - dist/linux-amd64/sentinel  (Linux x86_64)" 