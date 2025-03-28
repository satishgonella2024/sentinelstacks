#!/bin/bash

# Version from git tag or default
VERSION=${1:-$(git describe --tags --always --dirty)}
RELEASE_DIR="sentinelstacks-${VERSION}"

# Create release directory
mkdir -p "${RELEASE_DIR}"

# Build binary
go build -o "${RELEASE_DIR}/sentinel" cmd/sentinel/main.go

# Copy documentation
cp README.md LICENSE "${RELEASE_DIR}/"

# Copy examples
cp -r examples "${RELEASE_DIR}/"

# Create archive
tar -czf "${RELEASE_DIR}.tar.gz" "${RELEASE_DIR}"

# Cleanup
rm -rf "${RELEASE_DIR}"

echo "Release package created: ${RELEASE_DIR}.tar.gz" 