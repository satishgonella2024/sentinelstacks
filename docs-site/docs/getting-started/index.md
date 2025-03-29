# Getting Started with SentinelStacks

This guide will help you get up and running with SentinelStacks quickly.

## System Requirements

Before installing SentinelStacks, ensure your system meets these requirements:

### Minimum Requirements
- Go 1.21 or later
- Node.js 18+ (for desktop UI)
- Docker 20+ (for containerized deployment)

### LLM Provider Requirements
You'll need at least one of these:
- OpenAI API key
- Anthropic API key
- Ollama (for local models)

## Installation Methods

Choose the installation method that best suits your needs:

### 1. Homebrew Installation (macOS)
The easiest way to install on macOS:

```bash
brew install sentinelstacks
```

### 2. Docker Installation
For containerized deployment:

```bash
# Pull the latest image
docker pull sentinelstacks/sentinelstacks:latest

# Run with persistent storage
docker run -v ~/.sentinel:/root/.sentinel sentinelstacks/sentinelstacks:latest
```

### 3. Direct Download
Quick installation script:

```bash
curl -sSL https://get.sentinelstacks.io | sh
```

### 4. Build from Source
For development or customization:

```bash
# Clone the repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build the project
./scripts/build.sh
```

## Initial Configuration

### 1. Verify Installation
After installing, verify everything is working:

```bash
sentinel --version
```

### 2. Configure API Keys
Set up your preferred LLM provider:

```bash
# For OpenAI
export OPENAI_API_KEY=your_key_here

# For Anthropic
export ANTHROPIC_API_KEY=your_key_here

# For Ollama
# Make sure Ollama is installed and running
# No API key needed
```

### 3. Basic Usage

List available agents:
```bash
sentinel registry list
```

Create a new agent:
```bash
sentinel agent create --name test-agent
```

Run an agent:
```bash
sentinel agent run --name test-agent --interactive
```

## Next Steps

- [Using the CLI](../user-guide/cli.md)
- [Working with Agents](../user-guide/agents.md)
- [Memory System](../user-guide/memory.md)
- [Tool Integration](../user-guide/tools.md)

## Troubleshooting

### Common Issues

1. **API Key Not Found**
   ```bash
   # Check if API key is set
   echo $OPENAI_API_KEY
   # Set it if needed
   export OPENAI_API_KEY=your_key_here
   ```

2. **Docker Volume Permissions**
   ```bash
   # Fix permissions on sentinel directory
   sudo chown -R $USER:$USER ~/.sentinel
   ```

3. **Build Errors**
   ```bash
   # Update Go dependencies
   go mod tidy
   # Clean build artifacts
   go clean
   ```

### Getting Help

- Check our [GitHub Issues](https://github.com/satishgonella2024/sentinelstacks/issues)
- Read the [FAQ](../user-guide/faq.md)
- Review the [Architecture Overview](../developer-guide/architecture.md) 