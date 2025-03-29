#!/bin/bash

# Release script for SentinelStacks CLI

VERSION="1.0.0"
RELEASE_DIR="release"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Ensure clean release directory
rm -rf "${RELEASE_DIR}"
mkdir -p "${RELEASE_DIR}"

echo -e "${GREEN}Creating release packages for v${VERSION}...${NC}"

# Build binaries
./scripts/build.sh

# Package for macOS (Apple Silicon)
echo -e "${GREEN}Packaging for macOS (Apple Silicon)...${NC}"
cd dist/darwin-arm64
tar czf "../../${RELEASE_DIR}/sentinel-darwin-arm64.tar.gz" sentinel
cd ../..
ARM64_SHA=$(shasum -a 256 "${RELEASE_DIR}/sentinel-darwin-arm64.tar.gz" | cut -d ' ' -f 1)

# Package for macOS (Intel)
echo -e "${GREEN}Packaging for macOS (Intel)...${NC}"
cd dist/darwin-amd64
tar czf "../../${RELEASE_DIR}/sentinel-darwin-amd64.tar.gz" sentinel
cd ../..
AMD64_SHA=$(shasum -a 256 "${RELEASE_DIR}/sentinel-darwin-amd64.tar.gz" | cut -d ' ' -f 1)

# Package for Linux
echo -e "${GREEN}Packaging for Linux...${NC}"
cd dist/linux-amd64
tar czf "../../${RELEASE_DIR}/sentinel-linux-amd64.tar.gz" sentinel
cd ../..

echo -e "${GREEN}Release packages created!${NC}"
echo
echo "Release files are available in the '${RELEASE_DIR}' directory:"
ls -lh "${RELEASE_DIR}"
echo
echo "SHA256 hashes for Homebrew formula:"
echo "ARM64: ${ARM64_SHA}"
echo "AMD64: ${AMD64_SHA}"
echo
echo "Next steps:"
echo "1. Update Formula/sentinel.rb with the SHA256 hashes"
echo "2. Create a GitHub release with the .tar.gz files"
echo "3. Update the installation instructions in README.md" 