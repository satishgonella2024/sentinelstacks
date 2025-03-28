# Agentfile Specification

The Agentfile is the core configuration format for SentinelStacks agents. It defines how an agent behaves, what capabilities it has, and how it interacts with users and other systems.

## File Formats

An agent can be defined using two file formats:

1. **Natural Language (agentfile.natural.txt)**:
   - Human-readable description of the agent's purpose and behavior
   - Automatically converted to structured YAML

2. **Structured YAML (agentfile.yaml)**:
   - Machine-readable configuration
   - Generated from natural language or created directly
   - Used by the runtime to execute the agent

## Natural Language Format

The natural language format allows you to describe your agent in plain English. Here's an example:

```
This agent helps with research tasks. It should be able to
summarize academic papers, extract key information, and
answer questions about the content. It should maintain a
bibliography of sources it has processed. The agent should
use Llama3 as its model and should have access to a PDF parser.
```

## YAML Structure

The structured YAML follows this schema:

```yaml
name: agent-name
version: "1.0.0"
description: "A brief description of the agent"

# Model configuration
model:
  provider: "ollama"  # Could be ollama, openai, claude, etc.
  name: "llama3"      # Model name within the provider
  options:            # Provider-specific options
    temperature: 0.7

# Agent capabilities
capabilities:
  - text_processing
  - conversation
  - summarization

# Memory configuration
memory:
  type: "simple"      # simple, vector, etc.
  persistence: true   # Whether to save state between runs

# Tools the agent can use
tools:
  - id: "pdf-parser"
    version: "^1.0.0"

# Security permissions
permissions:
  file_access: ["read"]  # read, write, none
  network: false         # Whether the agent can access the network
```

## Required Fields

At minimum, an Agentfile must specify:
- `name`
- `model.provider`
- `model.name`

## Field Descriptions

### Top-Level Fields

- `name`: Unique identifier for the agent
- `version`: Semantic version of the agent definition
- `description`: Human-readable description of the agent

### Model Configuration

- `model.provider`: The AI model provider (ollama, openai, claude)
- `model.name`: The specific model to use
- `model.options`: Provider-specific parameters like temperature, max tokens, etc.

### Capabilities

A list of capabilities the agent has, which will be used to generate the system prompt. Common capabilities include:
- `conversation`: General dialogue ability
- `summarization`: Creating summaries of content
- `code_generation`: Writing code
- `search`: Finding information
- `tools_use`: Using specialized tools

### Memory Configuration

- `memory.type`: How the agent stores information (simple, vector, hierarchical)
- `memory.persistence`: Whether to save state between runs
- `memory.ttl`: Optional time-to-live for memory items
- `memory.capacity`: Optional maximum number of items to store

### Tools Configuration

Tools extend the agent's capabilities:
- `tools[].id`: The tool identifier
- `tools[].version`: Version requirement for the tool
- `tools[].config`: Optional tool-specific configuration

### Permissions

Security boundaries for the agent:
- `permissions.file_access`: Level of filesystem access (read, write, none)
- `permissions.network`: Whether the agent can make network requests
- `permissions.tools`: Specific tools the agent is allowed to use

## Version Compatibility

The Agentfile format uses semantic versioning:
- `1.0.0`: Initial stable release
- `1.x.y`: Backward-compatible changes
- `2.0.0`: Changes that may require updating agents

## Examples

### Research Assistant

```yaml
name: research-assistant
version: "1.0.0"
description: "Helps analyze academic papers"
model:
  provider: "ollama"
  name: "llama3"
  options:
    temperature: 0.3
capabilities:
  - summarization
  - question_answering
  - citation_management
memory:
  type: "vector"
  persistence: true
tools:
  - id: "pdf-parser"
    version: "^1.0.0"
  - id: "citation-formatter"
    version: "^1.0.0"
permissions:
  file_access: ["read"]
  network: false
```

### Code Assistant

```yaml
name: code-assistant
version: "1.0.0"
description: "Helps write and debug code"
model:
  provider: "openai"
  name: "gpt-4"
  options:
    temperature: 0.1
capabilities:
  - code_generation
  - code_explanation
  - debugging
memory:
  type: "simple"
  persistence: true
tools:
  - id: "code-analyzer"
    version: "^1.0.0"
permissions:
  file_access: ["read", "write"]
  network: false
```

## Converting Natural Language to YAML

Use the following command to convert a natural language description to YAML:

```bash
sentinel agentfile convert my-agent/agentfile.natural.txt
```

This will create `my-agent/agentfile.yaml` with the structured configuration.
