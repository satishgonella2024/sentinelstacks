#!/bin/bash

# Installation script for SentinelStacks CLI

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture to Go format
case "${ARCH}" in
    x86_64)  ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo -e "${RED}Unsupported architecture: ${ARCH}${NC}" && exit 1 ;;
esac

# Set binary path based on OS and architecture
BINARY_PATH="dist/${OS}-${ARCH}/sentinel"

# Check if binary exists
if [ ! -f "${BINARY_PATH}" ]; then
    echo -e "${RED}Binary not found: ${BINARY_PATH}${NC}"
    echo "Please run ./scripts/build.sh first"
    exit 1
fi

# Installation directory
INSTALL_DIR="/usr/local/bin"

# Create directory if it doesn't exist
sudo mkdir -p "${INSTALL_DIR}"

# Copy binary
echo -e "${GREEN}Installing SentinelStacks CLI to ${INSTALL_DIR}/sentinel...${NC}"
sudo cp "${BINARY_PATH}" "${INSTALL_DIR}/sentinel"
sudo chmod +x "${INSTALL_DIR}/sentinel"

echo -e "${GREEN}Installation complete!${NC}"
echo
echo "You can now use the 'sentinel' command. Try:"
echo "  sentinel --help" 