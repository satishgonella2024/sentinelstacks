#!/bin/bash

# Exit on error
set -e

# Create agent directory
AGENT_DIR="$HOME/.sentinel/agents/terraform-agent/latest"
mkdir -p "$AGENT_DIR"

# Build the agent
echo "Building Terraform agent..."
cd examples/terraform-agent
go build -o "$AGENT_DIR/terraform-agent" main.go

# Copy configuration
cp agent.yaml "$AGENT_DIR/"

echo "Agent installed successfully!"
echo "Location: $AGENT_DIR" 