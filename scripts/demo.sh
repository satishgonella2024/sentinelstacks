#!/bin/bash

# Exit on error
set -e

echo "🚀 Starting SentinelStacks Demo"

# Build the CLI
echo "📦 Building Sentinel CLI..."
go build -o sentinel cmd/sentinel/main.go
chmod +x sentinel

# Start development environment
echo "🌐 Starting development environment..."
./scripts/setup-dev.sh

# Login to registry
echo "🔑 Logging in to registry..."
./sentinel registry login --username admin --password admin

# Build and package Terraform agent
echo "🔧 Building Terraform agent..."
cd examples/terraform-agent
go build -o terraform-agent main.go
tar -czf terraform-agent.tar.gz terraform-agent agent.yaml
cd ../..

# Push agent to registry
echo "⬆️ Pushing agent to registry..."
./sentinel registry push --name terraform-agent --version latest

# Pull agent from registry
echo "⬇️ Pulling agent from registry..."
./sentinel registry pull --name terraform-agent --version latest

# List available agents
echo "📋 Listing available agents..."
./sentinel registry list

echo "✨ Demo completed successfully!"
echo "You can now run agents using:"
echo "  ./sentinel agent run --name terraform-agent --version latest" 