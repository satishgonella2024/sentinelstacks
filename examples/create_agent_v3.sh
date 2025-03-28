#!/bin/bash
# Script to create and run a new agent with SentinelStacks

# Navigate to the SentinelStacks root directory
cd ..

echo "Creating a new assistant agent directory..."
mkdir -p my-coding-assistant

# Create a custom natural language definition
echo "Creating a custom natural language definition..."
cat > my-coding-assistant/agentfile.natural.txt << EOL
Create an AI assistant that helps users with coding tasks.
It should be able to:
- Generate code examples in Python, JavaScript, and Go
- Explain programming concepts in simple terms
- Debug common errors
- Recommend best practices

The assistant should be friendly, patient, and tailored for beginners.
It should use the Llama3 model with a temperature of 0.5 to provide consistent answers.
EOL

# Convert the natural language to YAML
echo "Converting natural language to YAML configuration..."
./sentinel agentfile convert my-coding-assistant/agentfile.natural.txt

# Create an empty state file
echo "Creating empty state file..."
echo "{}" > my-coding-assistant/agent.state.json

# Run the agent
echo "Starting the agent... (Press Ctrl+C to exit)"
./sentinel agent run my-coding-assistant
