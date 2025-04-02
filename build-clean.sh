#!/bin/bash
# Clean build script

# Remove problematic files
rm -f cmd/sentinel/multimodal/multimodal_fixed.go

# Run the build
./build-sentinel.sh
