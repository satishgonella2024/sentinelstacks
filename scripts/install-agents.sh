#!/bin/bash

# Exit on error
set -e

# Install Terraform agent
echo "Installing Terraform agent..."
AGENT_DIR="$HOME/.sentinel/agents/terraform-agent/latest"
mkdir -p "$AGENT_DIR"
cd examples/terraform-agent
go build -o "$AGENT_DIR/terraform-agent" main.go
cp agent.yaml "$AGENT_DIR/"

# Install Kubernetes agent
echo "Installing Kubernetes agent..."
AGENT_DIR="$HOME/.sentinel/agents/kubernetes-agent/latest"
mkdir -p "$AGENT_DIR"
cd ../kubernetes-agent
go build -o "$AGENT_DIR/kubernetes-agent" main.go
cp agent.yaml "$AGENT_DIR/"

# Push agents to registry
echo "Pushing agents to registry..."
cd ../..
./sentinel registry push -name terraform-agent -version latest
./sentinel registry push -name kubernetes-agent -version latest

echo "Agents installed and pushed successfully!"
echo "You can now list available agents with:"
echo "  ./sentinel registry list" 