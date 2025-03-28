# SentinelStacks

SentinelStacks is a "Docker for AI agents" - providing tools to define, run, share, and orchestrate AI agents across different model providers.

## Features

- Natural language Agentfile definitions
- CLI-first development with GUI desktop interface
- Registry for discovering and sharing agents
- Pluggable foundation model backends (Ollama, OpenAI, Claude, etc.)
- State management for persistent agent memory

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

# Share your agent
sentinel registry push my-assistant
```

## Development Status

SentinelStacks is currently in early development. See the [ROADMAP.md](ROADMAP.md) for current status and planned features.

## License

[MIT License](LICENSE)
