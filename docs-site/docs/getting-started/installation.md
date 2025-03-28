# Installation

This guide will help you install SentinelStacks on your system.

## Prerequisites

- **Go**: Version 1.18 or higher
- **Ollama** (optional): For local model execution
- **Git**: For version control and repository operations

## Installing from Source

1. Clone the repository:

```bash
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks
```

2. Build and install the CLI:

```bash
go install ./cmd/sentinel
```

This will install the `sentinel` command in your `$GOPATH/bin` directory, which should be in your `PATH`.

## Verifying Installation

To verify that SentinelStacks is installed correctly, run:

```bash
sentinel version
```

You should see the version information for SentinelStacks.

## Installing Ollama (Recommended)

SentinelStacks works best with Ollama for local model execution. To install Ollama:

1. Visit [Ollama's website](https://ollama.ai/) and download the appropriate version for your system.

2. Follow the installation instructions.

3. Test that Ollama is working by running:

```bash
ollama run llama3
```

## Configuration

By default, SentinelStacks looks for Ollama at `http://localhost:11434`. If your Ollama instance is running elsewhere, you can configure it when creating or running agents:

```bash
# When converting an Agentfile
sentinel agentfile convert --endpoint http://your-ollama-server:11434 my-agent/agentfile.natural.txt

# When running an agent
sentinel agent run --endpoint http://your-ollama-server:11434 my-agent
```

## Next Steps

Now that you have SentinelStacks installed, you can [create your first agent](first-agent.md).
