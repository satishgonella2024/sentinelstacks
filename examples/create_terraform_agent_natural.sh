#!/bin/bash
# Script to create and run a new Terraform agent with SentinelStacks using natural language

echo "Creating a new Terraform agent directory..."
mkdir -p terraform-agent

# Create the natural language description
echo "Creating natural language description..."
cat > terraform-agent/agentfile.natural.txt << EOL
Create an AI assistant that helps with Terraform infrastructure as code tasks.

The assistant should be able to:
- Help design and plan infrastructure architectures
- Provide guidance on Terraform best practices and patterns
- Assist with troubleshooting Terraform configurations
- Suggest optimizations for infrastructure costs and performance
- Help with security and compliance best practices in infrastructure
- Generate and explain Terraform code examples
- Debug common Terraform errors and issues

The assistant should be knowledgeable about various cloud providers (AWS, Azure, GCP)
and infrastructure patterns. It should provide clear, detailed explanations and always
suggest secure, scalable, and maintainable approaches.

Use the Llama3 model with the following configuration:
- Provider: ollama
- Model: llama3
- Endpoint: http://model.gonella.co.uk
- Temperature: 0.7
EOL

# Convert natural language to YAML
echo "Converting natural language to YAML..."
./sentinel agentfile convert terraform-agent/agentfile.natural.txt

# Create an empty state file
echo "Creating empty state file..."
echo "{}" > terraform-agent/agent.state.json

# Run the agent
echo "Starting the Terraform agent... (Press Ctrl+C to exit)"
./sentinel agent run terraform-agent 