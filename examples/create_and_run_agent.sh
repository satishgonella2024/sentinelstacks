#!/bin/bash
# Example script to create and run a new agent with SentinelStacks

# Create a new agent
echo "Creating a new assistant agent..."
./sentinel agentfile create --name my-assistant

# Navigate to the agent directory
cd my-assistant

# Edit the natural language definition
echo "Creating a custom natural language definition..."
cat > agentfile.natural.txt << EOL
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
../sentinel agentfile convert agentfile.natural.txt

# Run the agent
echo "Starting the agent... (Press Ctrl+C to exit)"
cd ..
./sentinel agent run my-assistant
