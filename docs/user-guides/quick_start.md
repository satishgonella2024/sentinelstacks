# SentinelStacks Quick Start Guide

This guide will help you get started with SentinelStacks by walking through the installation and basic usage to create, build, run, and share your first AI agent.

## Installation

### Prerequisites

- Go 1.20 or later
- Git
- Access to an LLM API (Claude, OpenAI, etc.)

### Installing SentinelStacks CLI

```bash
# Install the Sentinel CLI
go install github.com/sentinelstacks/cli@latest

# Verify installation
sentinel version
```

### Configuration

Set up your LLM API access:

```bash
# Configure your default LLM provider
sentinel config set llm.provider claude

# Add your API key
sentinel config set llm.api_key your_api_key_here
```

## Creating Your First Agent

### Step 1: Initialize a New Agent

```bash
# Create a new directory for your agent
mkdir my-first-agent
cd my-first-agent

# Initialize a new Sentinelfile
sentinel init
```

This will create a basic Sentinelfile template in your current directory.

### Step 2: Edit the Sentinelfile

Open the Sentinelfile in your favorite text editor and describe your agent in natural language:

```
# Sentinelfile for WeatherAssistant

Create an agent that provides weather forecasts and recommendations.

The agent should be able to:
- Check current weather conditions for a location
- Provide 5-day forecasts
- Suggest clothing based on weather conditions
- Alert about severe weather events

The agent should use claude-3.7-sonnet as its base model.

It should maintain state about the user's location preferences and recent queries.

When the conversation starts, the agent should introduce itself as a weather assistant and ask for the user's location if not already known.

Allow the agent to access the following tools:
- Weather API
- Geolocation service

Set default_unit to metric.
Set refresh_interval to 30 minutes.
```

### Step 3: Build Your Agent

```bash
# Build the agent image
sentinel build -t username/weather-assistant:latest

# Verify the build
sentinel images
```

This will parse your natural language Sentinelfile into a structured definition and create a Sentinel Image.

### Step 4: Run Your Agent

```bash
# Run your agent
sentinel run username/weather-assistant:latest

# Or run with custom parameters
sentinel run --env refresh_interval=15 username/weather-assistant:latest
```

This will start your agent in interactive mode where you can chat with it directly.

### Step 5: Push to Registry

Share your agent with others by pushing it to the Sentinel Registry:

```bash
# Log in to the registry
sentinel login

# Push your agent
sentinel push username/weather-assistant:latest
```

## Using Sentinel Desktop

For a graphical interface to manage your agents:

1. Download Sentinel Desktop from the [official website](https://sentinelstacks.com/download)
2. Install and launch the application
3. Log in with your Sentinel Registry credentials
4. Access your agents from the "My Agents" section
5. Build, run, and monitor agents using the GUI

## Next Steps

- **Explore Examples**: Check out the [examples directory](https://github.com/sentinelstacks/examples) for more agent templates
- **Connect Tools**: Learn how to [connect external tools](https://docs.sentinelstacks.com/guides/connecting-tools) to your agents
- **Create Agent Networks**: Discover how to [orchestrate multiple agents](https://docs.sentinelstacks.com/guides/agent-networks) to work together
- **Customize Runtime**: Explore [advanced runtime configurations](https://docs.sentinelstacks.com/guides/runtime-config) for your agents

## Common Commands

```bash
# List running agents
sentinel ps

# Stop a running agent
sentinel stop agent_id

# View agent logs
sentinel logs agent_id

# Inspect an agent image
sentinel inspect username/weather-assistant:latest

# List available templates
sentinel templates list

# Create from template
sentinel init --template customer-support
```

## Troubleshooting

If you encounter issues:

1. Check the logs: `sentinel logs agent_id`
2. Verify your API key: `sentinel config check llm.api_key`
3. Update to the latest version: `go install github.com/sentinelstacks/cli@latest`
4. Seek help: `sentinel support`

For more detailed information, visit the [SentinelStacks documentation](https://docs.sentinelstacks.com).
