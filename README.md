# SentinelStacks

![SentinelStacks Logo](docs/visualizations/sentinelstacks_logo.svg)

SentinelStacks is an open-source AI agent management system that makes it easy to create, deploy, and share sophisticated LLM-powered agents.

[![Go Tests](https://github.com/sentinelstacks/sentinel/actions/workflows/go-test.yml/badge.svg)](https://github.com/sentinelstacks/sentinel/actions/workflows/go-test.yml)
[![Deploy Documentation](https://github.com/sentinelstacks/sentinel/actions/workflows/docs-deploy.yml/badge.svg)](https://github.com/sentinelstacks/sentinel/actions/workflows/docs-deploy.yml)
[![Docker Build](https://github.com/sentinelstacks/sentinel/actions/workflows/docker-build.yml/badge.svg)](https://github.com/sentinelstacks/sentinel/actions/workflows/docker-build.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## 🚀 Features

- **Natural Language Agent Definition**: Define agents using natural language or structured YAML
- **Multi-LLM Support**: Run agents across different LLM backends including Claude, Ollama, and soon OpenAI
- **Agent Management**: Build, run, and manage agents with a simple CLI interface
- **State Management**: Maintain agent state between runs and across sessions
- **Tool Integration**: Provide agents with access to tools like web search and calculators
- **Multi-Agent Orchestration**: Create systems of agents that work together
- **Registry System**: Share and discover agents through a central registry
- **NLP-to-Agent Generation**: Create agents on-the-fly from natural language descriptions

## 📋 Prerequisites

- Go 1.21 or later
- Access to at least one supported LLM provider

## 🔧 Installation

### Binary Installation

```bash
# Install with Go
go install github.com/sentinelstacks/sentinel@latest

# Verify installation
sentinel version
```

### From Source

```bash
# Clone the repository
git clone https://github.com/sentinelstacks/sentinel.git
cd sentinel

# Build the project
go build -o bin/sentinel cmd/sentinel/main.go

# Run the compiled binary
./bin/sentinel version
```

### Docker

```bash
# Pull the Docker image
docker pull ghcr.io/sentinelstacks/sentinel:latest

# Run in Docker
docker run --rm ghcr.io/sentinelstacks/sentinel:latest version
```

## 🏁 Quick Start

### Create Your First Agent

```bash
# Initialize a new agent
sentinel init --name mychatbot

# Edit the Sentinelfile
nano mychatbot/Sentinelfile

# Build the agent
sentinel build -t myusername/mychatbot:v1 -f mychatbot/Sentinelfile

# Run the agent
sentinel run myusername/mychatbot:v1
```

### Example Sentinelfile

```yaml
name: basicchatbot
description: Create a friendly chatbot with a helpful personality.
capabilities:
  - Engage in casual conversation
  - Answer general knowledge questions
  - Maintain a consistent personality
model: 
  base: claude3
state:
  - conversation_history
initialization:
  introduction: "Hello! I'm ready to assist you."
tools:
  - web_search:
      purpose: For looking up factual information
personality:
  tone: friendly
  response_length: medium
```

## 📖 Documentation

For comprehensive documentation, visit our [Documentation Site](https://sentinelstacks.github.io/sentinel/).

- [Architecture Overview](https://sentinelstacks.github.io/sentinel/architecture/README/)
- [User Guides](https://sentinelstacks.github.io/sentinel/user-guides/README/)
- [Example Agents](https://sentinelstacks.github.io/sentinel/examples/chatbot/)
- [Advanced Agent Design](https://sentinelstacks.github.io/sentinel/user-guides/advanced_agents/)

## 💡 Example Agents

SentinelStacks comes with several example agents that demonstrate different capabilities:

- **Basic Chatbot**: Simple conversational agent
- **Research Assistant**: Advanced information gathering and synthesis
- **Team Collaboration**: Multi-agent system with specialized roles
- **Financial Advisor**: Domain-specific agent with compliance controls

## 🛣️ Roadmap

See our [Development Roadmap](https://sentinelstacks.github.io/sentinel/planning/roadmap/) for planned features and enhancements.

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

SentinelStacks is released under the MIT License. See [LICENSE](LICENSE) for details.
