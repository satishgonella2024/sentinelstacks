#!/bin/bash
set -e

echo "Running SentinelStacks with fixes..."

# Navigate to project root
cd "$(dirname "$0")/.."

# 1. Fix Go dependencies
echo "Fixing Go dependencies..."
go get github.com/mattn/go-isatty@v0.0.20
go mod tidy

# 2. Run the application
echo "Starting the application..."
./scripts/dev.sh
