#!/bin/bash
set -e

echo "Building registry package..."
go build ./cmd/sentinel/registry
if [ $? -eq 0 ]; then
  echo "✅ Registry package built successfully!"
else
  echo "❌ Registry package build failed."
  exit 1
fi

echo "Building all packages..."
go build ./...
if [ $? -eq 0 ]; then
  echo "✅ All packages built successfully!"
else
  echo "❌ Full build failed. Some import cycles may still exist."
  exit 1
fi

echo "All done!"
