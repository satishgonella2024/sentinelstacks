#!/bin/bash
# Test script to run the images command without building the full project

set -e  # Exit immediately if a command fails

echo "Building and running the test..."
cd /Users/subrahmanyagonella/the-repo/sentinelstacks

# Build the test
go build -o test-images test-images.go override-images.go

# Make the executable
chmod +x test-images

# Run the test
./test-images

echo "Test complete!"
