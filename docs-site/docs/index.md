# SentinelStacks - Docker for AI Agents

SentinelStacks is a platform for creating, running, sharing, and orchestrating AI agents.

## What is SentinelStacks?

SentinelStacks provides tools similar to what Docker did for containers, but for AI agents:

- **Natural Language Agentfiles**: Define agents in plain English, which get converted to structured YAML
- **CLI-First Development**: Command-line interface for developers with full control
- **Desktop GUI**: Visual interface for monitoring and managing agents
- **Registry**: Discover, share, and reuse agents
- **Model-Agnostic**: Works with Ollama, OpenAI, Claude, and other LLM providers

## Key Features

- **Local-First Execution**: Run models locally for privacy and cost control
- **Natural Language Configuration**: Describe your agent in plain language
- **Standardized Format**: Consistent structure for agent definitions
- **Composable Agents**: Build complex systems from simpler agents
- **State Management**: Persist agent memory between runs

## Quick Start

```bash
# Install SentinelStacks
go install github.com/satishgonella2024/sentinelstacks/cmd/sentinel@latest

# Create a new agent
sentinel agentfile create --name my-assistant

# Edit the Agentfile natural language definition
# This will be automatically converted to YAML

# Run your agent
sentinel agent run my-assistant
```

## Project Status

SentinelStacks is currently in early development. See the [Roadmap](roadmap.md) for current status and planned features.

## License

[MIT License](https://github.com/satishgonella2024/sentinelstacks/blob/main/LICENSE)
