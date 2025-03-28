# Getting Started with SentinelStacks

This guide will help you get up and running with SentinelStacks, a platform for creating, running, and sharing AI agents.

## Installation

### Prerequisites
- Go 1.18 or higher
- An Ollama installation (for local model execution)

### Installing SentinelStacks CLI

```bash
# Clone the repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build and install the CLI
go install ./cmd/sentinel
```

## Creating Your First Agent

1. Create a new agent:

```bash
sentinel agentfile create --name my-assistant
```

2. Edit the natural language description:

```bash
# This will create and open agentfile.natural.txt
nano my-assistant/agentfile.natural.txt
```

Add a description of what you want your agent to do, for example:

```
This agent helps answer questions about Go programming.
It should provide code examples when asked and explain
concepts clearly. It should be friendly and helpful.
```

3. Convert the natural language to a structured Agentfile:

```bash
sentinel agentfile convert my-assistant/agentfile.natural.txt
```

This will create `my-assistant/agentfile.yaml` with the structured configuration.

## Running Your Agent

To run your agent:

```bash
sentinel agent run my-assistant
```

This will start an interactive session with your agent.

## Sharing Your Agent

To share your agent with others:

```bash
sentinel registry push my-assistant
```

## Next Steps

- Explore more advanced agent capabilities
- Create multi-agent workflows
- Set up the desktop application for visual management
