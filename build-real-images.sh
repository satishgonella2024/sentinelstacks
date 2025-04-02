#!/bin/bash
# Build and run the real-images program

set -e  # Exit immediately if a command fails

echo "Building real-images..."
cd /Users/subrahmanyagonella/the-repo/sentinelstacks

go build -o real-images real-images.go
chmod +x real-images

echo "Running real-images..."
./real-images

echo "Done!"
