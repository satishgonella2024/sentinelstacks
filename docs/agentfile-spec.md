# Agentfile Specification

This document describes the format and structure of Agentfiles used by SentinelStacks. An Agentfile defines an AI agent's behavior, capabilities, and configuration in a structured YAML format.

## Overview

Agentfiles can be defined in two ways:
1. Natural language description (`agentfile.natural.txt`)
2. Structured YAML configuration (`agentfile.yaml`)

The natural language description can be automatically converted to YAML using the SentinelStacks CLI:
```bash
sentinel agentfile convert agentfile.natural.txt
```

## Schema

An Agentfile must conform to the following schema:

```yaml
name: string                # Name of the agent
version: string             # Semantic version (e.g., "0.1.0")
description: string         # Short description of the agent's purpose
model:                      # Configuration for the AI model
  provider: string          # Model provider (ollama, openai, claude)
  name: string              # Model name (llama3, gpt-4, claude-3-sonnet, etc.)
  options:                  # Provider-specific options
    temperature: float      # Creativity/randomness (0.0-1.0)
    max_tokens: int         # Maximum output tokens (optional)
    top_p: float            # Nucleus sampling parameter (optional)
capabilities:               # List of agent capabilities
  - string                  # e.g., conversation, code_generation, etc.
memory:                     # Memory/state configuration
  type: string              # Memory type (simple, vector)
  persistence: boolean      # Whether to persist state between runs
tools:                      # Optional list of tools the agent can use
  - id: string              # Tool identifier
    version: string         # Tool version (semantic version)
permissions:                # Optional security constraints
  file_access: [string]     # File access permissions (read, write, none)
  network: boolean          # Whether network access is allowed
author: string              # Optional author information
tags: [string]              # Optional tags for categorization
registry:                   # Optional registry metadata
  source: string            # Registry source
  visibility: string        # public or private
```

## Fields in Detail

### `name`
Identifies the agent. Should be unique, kebab-case, and descriptive.

### `version`
Semantic versioning string in the format `MAJOR.MINOR.PATCH`. 

### `description`
A concise description of what the agent does.

### `model`
Configuration for the AI model powering the agent.

#### `provider`
One of:
- `ollama`: Local models via Ollama
- `openai`: OpenAI API (requires OPENAI_API_KEY)
- `claude`: Anthropic Claude API (requires ANTHROPIC_API_KEY)

#### `name`
Model name depends on the provider:
- Ollama: `llama3`, `mistral`, etc.
- OpenAI: `gpt-4`, `gpt-3.5-turbo`, etc.
- Claude: `claude-3-opus-20240229`, `claude-3-sonnet-20240229`, etc.

#### `options`
Provider-specific parameters:
- `temperature`: Controls randomness/creativity (0.0-1.0)
- `max_tokens`: Maximum response length in tokens
- `top_p`: Nucleus sampling parameter (alternative to temperature)

### `capabilities`
List of capabilities the agent supports. Common values include:
- `conversation`: Basic chat functionality
- `code_generation`: Writing code
- `explanation`: Explaining concepts
- `summarization`: Summarizing content
- `research`: Looking up information
- `translation`: Translating languages

### `memory`
Configuration for the agent's state management.

#### `type`
Memory storage type:
- `simple`: Basic key-value storage
- `vector`: Vector-based memory for semantic search (not yet implemented)

#### `persistence`
Boolean indicating whether state should be saved between sessions.

### `tools` (Optional)
External tools the agent can use. Each tool has:
- `id`: Tool identifier
- `version`: Semantic version

### `permissions` (Optional)
Security constraints for the agent:

#### `file_access`
Array of file access permissions:
- `read`: Read-only access
- `write`: Write access
- `none`: No file access

#### `network`
Boolean indicating whether network access is allowed.

### `author` (Optional)
Username or organization that created the agent.

### `tags` (Optional)
Array of tags for categorization.

### `registry` (Optional)
Metadata for registry operations.

#### `source`
Registry source URL or identifier.

#### `visibility`
- `public`: Visible to all users
- `private`: Restricted visibility

## Example

```yaml
name: code-assistant
version: "0.1.0"
description: "An agent that helps with coding tasks and explains concepts"
model:
  provider: ollama
  name: llama3
  options:
    temperature: 0.5
capabilities:
  - conversation
  - code_generation
  - explanation
memory:
  type: simple
  persistence: true
permissions:
  file_access: ["read"]
  network: false
author: sentinelstacks
tags:
  - coding
  - tutorial
```

## Using Environment Variables

Some configuration options can be set via environment variables:
- `OLLAMA_ENDPOINT`: Override the default Ollama API endpoint
- `OPENAI_API_KEY`: API key for OpenAI models
- `ANTHROPIC_API_KEY`: API key for Claude models

## Best Practices

1. **Start Simple**: Begin with basic capabilities and add more as needed
2. **Use Natural Language**: Define your agent in plain English first, then convert to YAML
3. **Test Thoroughly**: Verify agent behavior with different inputs
4. **Version Carefully**: Increment versions when making significant changes
5. **Document Well**: Include clear descriptions and usage examples
