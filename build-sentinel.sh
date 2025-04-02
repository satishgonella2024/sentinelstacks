#!/bin/bash
# Script to build SentinelStacks with real functionality

set -e  # Exit on error

echo "Building SentinelStacks with real functionality..."

# Ensure directories exist
mkdir -p ~/.sentinel/images

# Build the project
echo "Building the project..."
go build -o sentinel main.go

# Make the executable
echo "Making the sentinel executable..."
chmod +x sentinel

echo "Build complete! You can now run SentinelStacks with real functionality."
echo ""
echo "Try these commands:"
echo "  ./sentinel init --template chatbot my-agent"
echo "  cd my-agent"
echo "  ../sentinel build -t myname/chatbot:latest"
echo "  ../sentinel images"
echo "  ../sentinel run myname/chatbot:latest"
echo ""
echo "For LLM integration, add your API keys to ~/.sentinel/config.yaml:"
echo "  llm:"
echo "    provider: claude"
echo "    api_key: your_claude_api_key"
echo "    model: claude-3.7-sonnet"
echo "  openai:"
echo "    api_key: your_openai_api_key"
echo "  ollama:"
echo "    endpoint: http://localhost:11434"
