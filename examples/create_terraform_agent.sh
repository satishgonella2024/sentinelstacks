#!/bin/bash
# Script to create and run a new Terraform agent with SentinelStacks

echo "Creating a new Terraform agent directory..."
mkdir -p terraform-agent

# Create a YAML file directly
echo "Creating a YAML configuration file..."
cat > terraform-agent/agentfile.yaml << EOL
name: terraform-agent
version: "0.1.0"
description: "An AI assistant that helps users with Terraform infrastructure as code tasks"
model:
  provider: ollama
  name: llama3
  endpoint: "http://model.gonella.co.uk"
  options:
    temperature: 0.7
capabilities:
  - terraform_planning
  - infrastructure_design
  - resource_optimization
  - security_best_practices
  - troubleshooting
memory:
  type: simple
  persistence: true
EOL

# Create an empty state file
echo "Creating empty state file..."
echo "{}" > terraform-agent/agent.state.json

# Run the agent
echo "Starting the Terraform agent... (Press Ctrl+C to exit)"
./sentinel agent run terraform-agent 