# Welcome to SentinelStacks

🤖 SentinelStacks is an AI-powered infrastructure management platform that helps organizations automate, secure, and manage their cloud resources using intelligent agents.

## Current Status (March 2024)

SentinelStacks is actively being developed with regular updates and improvements. Here's the current implementation status:

### ✅ Production Ready (100%)

- **Model Adapters**
  - Full support for OpenAI (GPT-3.5, GPT-4)
  - Anthropic Claude integration
  - Local models via Ollama
  - Streaming responses and error handling
  - Model-specific parameter configuration

- **CLI Tool Enhancements**
  - Animated progress indicators
  - Color-coded output
  - Interactive mode
  - Clear success/error states
  - Command completion

- **Tools Framework**
  - Calculator with advanced operations
  - URL fetcher with caching
  - Weather service integration
  - Terraform infrastructure management
  - Custom tool development API

### 🔄 In Development

- **Memory System (80%)**
  - Vector-based storage implemented
  - Semantic search capabilities
  - Basic context management
  - Persistence layer complete
  - Optimization and advanced features in progress

- **Desktop UI (35%)**
  - Tauri foundation established
  - React components in development
  - Basic agent management interface
  - Settings panel implementation
  - Real-time monitoring planned

- **Registry System (60%)**
  - Basic agent storage
  - Version tracking in progress
  - Search functionality in development
  - Sharing capabilities planned
  - Security features being implemented

## System Requirements

### Core Requirements
- Go 1.21 or later
- Node.js 18+ (for desktop UI)
- Docker 20+ (for containerized deployment)
- 4GB RAM minimum (8GB recommended)
- 2 CPU cores minimum (4 cores recommended)

### LLM Provider Requirements
Choose one of:
- OpenAI API key (GPT-3.5/4)
- Anthropic API key (Claude)
- Ollama installation (local models)

### Optional Components
- Rust (for desktop UI development)
- PostgreSQL 13+ (for local deployment)
- Redis 6+ (for caching)

## Installation Methods

### 1. Homebrew Installation (macOS)
```bash
# Install using Homebrew
brew install sentinelstacks

# Verify installation
sentinel --version
```

### 2. Docker Installation
```bash
# Pull the latest image
docker pull sentinelstacks/sentinelstacks:latest

# Run with persistent storage
docker run -v ~/.sentinel:/root/.sentinel sentinelstacks/sentinelstacks:latest
```

### 3. Direct Download
```bash
# Download and install
curl -sSL https://get.sentinelstacks.io | sh

# Verify installation
sentinel --version
```

### 4. Build from Source
```bash
# Clone repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build the project
./scripts/build.sh

# Run tests
go test ./...
```

## Configuration

### API Keys
```bash
# OpenAI configuration
export OPENAI_API_KEY=your_key_here

# Anthropic configuration
export ANTHROPIC_API_KEY=your_key_here

# Local model configuration
# Install Ollama from https://ollama.ai
```

### Basic Usage
```bash
# List available agents
sentinel registry list

# Create a new agent
sentinel agent create --name test-agent

# Run an agent interactively
sentinel agent run --name test-agent --interactive
```

## Features

### AI Integration
- Multiple model support with dynamic switching
- Streaming responses for real-time interaction
- Robust error handling and automatic retries
- Model-specific parameter optimization
- Context window management

### Memory System
- Vector-based storage for efficient retrieval
- Semantic search with customizable thresholds
- Persistent storage with backup options
- Context management with priority queuing
- Automatic garbage collection (planned)

### Tools Framework
- Built-in calculator with scientific functions
- URL fetcher with caching and retry logic
- Weather service with multiple provider support
- Terraform integration for infrastructure
- Custom tool development SDK

### Infrastructure
- Docker-based deployment with compose
- Nginx reverse proxy with SSL
- PostgreSQL for persistent storage
- Redis for high-speed caching
- Horizontal scaling support

## Known Limitations

### Memory System
- Limited context window size (8K tokens)
- Basic vector optimization algorithms
- No automatic memory pruning
- Limited cross-agent memory sharing

### Desktop UI
- Early development stage (35% complete)
- Basic functionality only
- No real-time updates yet
- Limited customization options

### Registry
- Basic CRUD operations only
- No semantic search in registry
- Limited version control
- Basic access controls

## Next Steps

- [Getting Started Guide](getting-started/index.md)
- [User Guide](user-guide/index.md)
- [API Reference](user-guide/api-reference.md)
- [Contributing Guide](developer-guide/contributing.md)
- [Architecture Overview](developer-guide/architecture.md)

## Community & Support

- [GitHub Repository](https://github.com/satishgonella2024/sentinelstacks.git)
- [Issue Tracker](https://github.com/satishgonella2024/sentinelstacks/issues)
- [Documentation](https://satishgonella2024.github.io/sentinelstacks/)
- [Community Discord](https://discord.gg/sentinelstacks)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/satishgonella2024/sentinelstacks/blob/main/LICENSE) file for details.
