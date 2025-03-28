# SentinelStacks

SentinelStacks is a "Docker for AI agents" - providing tools to define, run, share, and orchestrate AI agents across different model providers.

## Features

- Natural language Agentfile definitions
- CLI-first development with GUI desktop interface
- Registry for discovering and sharing agents
- Pluggable foundation model backends (Ollama, OpenAI, Claude, etc.)
- State management for persistent agent memory

## Installation

### Prerequisites

- Go 1.20 or later
- [Ollama](https://ollama.com/) for local LLM support (optional)

### Installing from Source

```bash
# Clone the repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build the CLI
go build -o sentinel cmd/sentinel/main.go

# Add to your PATH (optional)
sudo mv sentinel /usr/local/bin/
```

### Environment Setup

If you want to use OpenAI or Claude models, set the following environment variables:

```bash
# For OpenAI models
export OPENAI_API_KEY="your-api-key"

# For Claude models
export ANTHROPIC_API_KEY="your-api-key"

# For custom Ollama endpoint
export OLLAMA_ENDPOINT="http://your-ollama-server:11434"
```

## Quick Start

```bash
# Create a new agent
sentinel agentfile create --name my-assistant

# Edit the Agentfile natural language definition
vim my-assistant/agentfile.natural.txt

# Convert to YAML
sentinel agentfile convert my-assistant/agentfile.natural.txt

# Run your agent
sentinel agent run my-assistant

# Share your agent
sentinel registry push my-assistant
```

## Usage

### Creating an Agent

```bash
# Create a new agent with a name
sentinel agentfile create --name my-agent-name
```

This creates a directory with the agent name containing:
- `agentfile.natural.txt`: Natural language definition
- `agentfile.yaml`: YAML configuration
- `agent.state.json`: Empty state file

### Converting Natural Language to YAML

```bash
# Convert natural language to structured YAML
sentinel agentfile convert path/to/agentfile.natural.txt
```

This uses an LLM to convert your natural language description into a structured YAML configuration.

### Running an Agent

```bash
# Run an agent with the default model endpoint
sentinel agent run my-agent-name

# Run with a custom model endpoint
sentinel agent run my-agent-name --endpoint http://custom-model-server:11434
```

### Registry Operations

```bash
# Push an agent to the registry
sentinel registry push my-agent-name

# Pull an agent from the registry
sentinel registry pull username/agent-name

# Search the registry
sentinel registry search "coding assistant" --tags python,tutorial

# List all agents in the registry
sentinel registry list
```

### Advanced Usage

#### Using Different Model Providers

Edit the `model` section in your `agentfile.yaml`:

```yaml
model:
  provider: openai  # or claude, ollama
  name: gpt-4       # model name
  options:
    temperature: 0.7
```

Don't forget to set the required environment variables for API access.

## Project Structure

```
sentinelstacks/
├── cmd/
│   └── sentinel/           # CLI application
│       ├── commands/       # Command implementations
│       └── main.go         # Entry point
├── pkg/
│   ├── agentfile/          # Agentfile parser and schema
│   ├── models/             # Model adapters (Ollama, OpenAI)
│   ├── registry/           # Registry client and server
│   └── runtime/            # Agent execution runtime
├── docs/                   # Documentation assets
├── docs-site/              # MkDocs documentation site
│   ├── docs/               # Markdown documentation
│   └── mkdocs.yml          # MkDocs configuration
└── examples/               # Example agent definitions
```

## Documentation

Comprehensive documentation is available in the `/docs-site` directory. To view the documentation locally:

1. Install MkDocs and the Material theme:
   ```bash
   pip install mkdocs-material
   ```

2. Run the local development server:
   ```bash
   cd docs-site
   mkdocs serve
   ```

3. Open your browser to [http://localhost:8000](http://localhost:8000)

## Development Status

SentinelStacks is currently in early development. See the [ROADMAP.md](ROADMAP.md) for current status and planned features.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](docs-site/docs/developer-guide/contributing.md) for guidelines.

## License

[MIT License](LICENSE)
