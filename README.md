# SentinelStacks

SentinelStacks is a comprehensive system for creating, managing, and orchestrating AI agents using natural language definitions.

## Overview

SentinelStacks provides a Docker-like workflow for AI agents:

- Define agents using Sentinelfiles
- Build agent images from Sentinelfiles
- Run agents from images locally or from registries
- Share agents through registries
- Stack multiple agents together to create complex workflows

## Key Features

- **Natural Language Agent Definition**: Create agents by describing what they should do in natural language
- **Local Agent Execution**: Run agents on your local machine with simple commands
- **Agent Registry**: Share and reuse agents through a registry system
- **Multi-Agent Stacks**: Define complex workflows with multiple agents working together
- **Context Propagation**: Pass data between agents in a stack
- **CLI-First Design**: Complete command-line interface for all operations

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/subrahmanyagonella/the-repo/sentinelstacks.git
cd sentinelstacks

# Build the binary
make build
```

### Create Your First Agent

```bash
# Create an agent from a natural language description
sentinel init "An agent that summarizes text input" --name summarizer

# Build the agent
sentinel build -t summarizer:latest .

# Run the agent
sentinel run summarizer:latest
```

### Create a Multi-Agent Stack

```bash
# Initialize a stack
sentinel stack init my-analysis --template analyzer

# Run the stack
sentinel stack run -f Stackfile.yaml --input="This is a test input"
```

## Core Concepts

### Agents

Agents are the basic building blocks in SentinelStacks. An agent:
- Has a specific purpose
- Takes inputs and produces outputs
- Is defined by a Sentinelfile
- Is built into an agent image
- Can be run locally or pulled from a registry

### Stacks

Stacks are collections of agents arranged in a workflow:
- Define a directed acyclic graph (DAG) of agent execution
- Manage data flow between agents
- Execute agents in the correct order based on dependencies
- Allow complex processing pipelines with specialized agents

### Sentinelfile

A Sentinelfile defines an agent's behavior:
```yaml
name: text-summarizer
description: Summarizes text input
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: text
    type: string
    description: "Text to summarize"
  - name: max_length
    type: integer
    default: 100
    description: "Maximum length of summary"
output_format: "text"
system_prompt: |
  You are a specialized text summarization agent.
prompt_template: |
  Summarize the following text in {{max_length}} words or less:
  
  {{text}}
```

### Stackfile

A Stackfile defines a multi-agent workflow:
```yaml
name: data-analysis-pipeline
description: Analyzes and summarizes data
version: 1.0.0
agents:
  - id: extractor
    uses: data-extractor:latest
    params:
      format: "json"
  - id: analyzer
    uses: data-analyzer:latest
    inputFrom:
      - extractor
    params:
      analysis_type: "comprehensive"
  - id: summarizer
    uses: text-summarizer:latest
    inputFrom:
      - analyzer
    params:
      max_length: 200
```

## Command Reference

### Agent Commands

- `sentinel init` - Initialize a new agent
- `sentinel build` - Build an agent from a Sentinelfile
- `sentinel run` - Run an agent locally
- `sentinel push` - Push an agent to a registry
- `sentinel pull` - Pull an agent from a registry
- `sentinel stop` - Stop a running agent
- `sentinel logs` - View logs from a running agent
- `sentinel images` - List available agent images

### Stack Commands

- `sentinel stack init` - Initialize a new stack
- `sentinel stack run` - Run a multi-agent stack
- `sentinel stack list` - List available stacks
- `sentinel stack inspect` - Inspect a stack configuration

## Documentation

- [Getting Started Guide](docs/SETUP.md)
- [Agent Creation Guide](docs/AGENTS.md)
- [Stack Creation Guide](docs/STACK-README.md)
- [Registry Guide](docs/REGISTRY.md)
- [API Reference](docs/API.md)

## Architecture

SentinelStacks follows a modular architecture:

- **Core Runtime**: Executes agents and manages their lifecycle
- **Stack Engine**: Orchestrates multi-agent workflows
- **Registry System**: Manages agent sharing and discovery
- **Parser**: Converts natural language to agent configurations
- **CLI**: Provides command-line interface to all features
- **Web UI**: (Coming soon) Visual interface for managing agents and stacks

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
