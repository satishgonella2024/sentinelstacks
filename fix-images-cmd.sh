#!/bin/bash
# Quick Fix for SentinelStacks Images Command

set -e  # Exit immediately if a command fails

echo "Starting images command fix..."
cd /Users/subrahmanyagonella/the-repo/sentinelstacks

# Build the project to apply changes to images.go
echo "Building the project..."
go build -o sentinel main.go

# Make the executable
echo "Making the sentinel executable..."
chmod +x sentinel

echo "Fix complete! Try running './sentinel images' now to see if it properly reads from the registry."
