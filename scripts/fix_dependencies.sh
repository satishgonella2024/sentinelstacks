#!/bin/bash
set -e

echo "Fixing dependencies..."

# Navigate to the project root
cd "$(dirname "$0")/.."

# Add the missing dependencies
go get github.com/mattn/go-isatty@v0.0.20
go mod tidy

echo "Dependencies fixed!"
