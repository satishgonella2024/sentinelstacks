# Agentfile Specification

The Agentfile is the core configuration format for SentinelStacks agents. It defines how an agent behaves, what capabilities it has, and how it interacts with users and other systems.

## File Formats

An agent can be defined using two file formats:

### Natural Language (agentfile.natural.txt)

This human-readable format describes the agent's purpose and behavior in plain English. For example:

```
This agent helps with research tasks. It should be able to
summarize academic papers, extract key information, and
answer questions about the content. It should maintain a
bibliography of sources it has processed. The agent should
use Llama3 as its model and should have access to a PDF parser.
```

### Structured YAML (agentfile.yaml)

This machine-readable format is used by the runtime to execute the agent. It's either generated from the natural language description or created directly.

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

## YAML Structure

### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique identifier for the agent |
| `version` | string | Yes | Semantic version of the agent definition |
| `description` | string | No | Human-readable description of the agent |

### Model Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `model.provider` | string | Yes | The AI model provider (ollama, openai, claude) |
| `model.name` | string | Yes | The specific model to use |
| `model.options` | object | No | Provider-specific parameters |

### Capabilities

A list of capabilities the agent has, which will be used to generate the system prompt:

- `conversation`: General dialogue ability
- `summarization`: Creating summaries of content
- `code_generation`: Writing code
- `search`: Finding information
- `tools_use`: Using specialized tools

### Memory Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `memory.type` | string | Yes | How the agent stores information |
| `memory.persistence` | boolean | Yes | Whether to save state between runs |
| `memory.ttl` | string | No | Time-to-live for memory items |
| `memory.capacity` | integer | No | Maximum number of items to store |

### Tools Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tools[].id` | string | Yes | The tool identifier |
| `tools[].version` | string | No | Version requirement for the tool |
| `tools[].config` | object | No | Tool-specific configuration |

### Permissions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `permissions.file_access` | array | No | Level of filesystem access |
| `permissions.network` | boolean | No | Whether the agent can make network requests |

## Examples

### Conversation Assistant

```yaml
name: conversation-assistant
version: "1.0.0"
description: "General purpose conversational agent"
model:
  provider: "ollama"
  name: "llama3"
  options:
    temperature: 0.7
capabilities:
  - conversation
  - question_answering
memory:
  type: "simple"
  persistence: true
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

## Converting Between Formats

Use the CLI to convert natural language to YAML:

```bash
sentinel agentfile convert myagent/agentfile.natural.txt
```

This will create or update `myagent/agentfile.yaml` based on the natural language description.
