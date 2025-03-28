#!/bin/bash
# Script to create and run a new agent with SentinelStacks using direct YAML

# Navigate to the SentinelStacks root directory
cd ..

echo "Creating a new assistant agent directory..."
mkdir -p code-helper

# Create a YAML file directly
echo "Creating a YAML configuration file..."
cat > code-helper/agentfile.yaml << EOL
name: code-helper
version: "0.1.0"
description: "An AI assistant that helps users with coding tasks, tailored for beginners"
model:
  provider: ollama
  name: llama3
  options:
    temperature: 0.5
capabilities:
  - conversation
  - code_generation
  - debugging
  - explanation
memory:
  type: simple
  persistence: true
EOL

# Create an empty state file
echo "Creating empty state file..."
echo "{}" > code-helper/agent.state.json

# Run the agent
echo "Starting the agent... (Press Ctrl+C to exit)"
./sentinel agent run code-helper
EOL
