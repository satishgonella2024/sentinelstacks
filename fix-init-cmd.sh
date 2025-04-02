#!/bin/bash
# Quick Fix for SentinelStacks Init Command

set -e  # Exit immediately if a command fails

echo "Starting init command fix..."
cd /Users/subrahmanyagonella/the-repo/sentinelstacks

# Build the project to apply existing changes in init.go
echo "Building the project..."
go build -o sentinel main.go

# Make the executable
echo "Making the sentinel executable..."
chmod +x sentinel

echo "Fix complete! Try running './sentinel init --template chatbot my-agent' now."
