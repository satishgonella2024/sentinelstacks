#!/bin/bash

# Prepare for commit
echo "Preparing for commit..."

# Make all scripts executable
chmod +x scripts/*.sh

# Run go mod tidy to clean up dependencies
go mod tidy

# Check for build issues
go build -o sentinel-test ./main.go
if [ $? -ne 0 ]; then
    echo "Error: Build failed!"
    exit 1
fi

# Clean up test binary
rm -f sentinel-test

echo "Ready for commit. The following files have been modified:"
git status --porcelain

echo "
Commit message suggestion:
feat: Implement stack engine with memory management

- Renamed 'compose' to 'stack' with full functionality
- Added complete memory management system with multiple backends
- Implemented flexible agent runtime system
- Added CLI commands for managing memory
- Updated documentation with implementation progress
"
