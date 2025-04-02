#!/bin/bash

# Test script for stack functionality
echo "Testing stack functionality..."

# Check if sentinel binary exists
if [ ! -f "./sentinel" ]; then
    echo "Error: sentinel binary not found. Run 'make build' first."
    exit 1
fi

# Initialize memory subsystem
echo "Initializing memory subsystem..."
./sentinel memory init

# Create a simple test stack file
cat > test-stack.yaml <<EOF
name: test-stack
description: A simple test stack for validation
version: 1.0.0
agents:
  - id: processor
    uses: processor
    params:
      format: "json"
  - id: analyzer
    uses: analyzer
    inputFrom:
      - processor
    params:
      analysis_type: "comprehensive"
  - id: summarizer
    uses: summarizer
    inputFrom:
      - analyzer
    params:
      format: "bullet_points"
EOF

# Run the stack with verbose output
echo "Running the stack..."
./sentinel stack run -f test-stack.yaml -v

# Clean up
echo "Cleaning up..."
rm test-stack.yaml

echo "Stack test completed."
