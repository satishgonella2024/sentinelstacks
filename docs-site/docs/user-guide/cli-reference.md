# CLI Reference

The SentinelStacks CLI (`sentinel`) provides a command-line interface for managing and interacting with agents.

## Global Commands

### `sentinel version`

Shows version information for the CLI.

```bash
sentinel version
```

**Output Example:**
```
SentinelStacks v0.1.0
Git commit: development
Built on: unknown
```

### `sentinel --help`

Shows general help information and lists available commands.

```bash
sentinel --help
```

## Agentfile Commands

Commands for managing agent definitions.

### `sentinel agentfile create`

Creates a new agent definition.

```bash
sentinel agentfile create --name <agent-name>
```

**Flags:**
- `--name, -n`: Name of the agent (required)

**Example:**
```bash
sentinel agentfile create --name my-assistant
```

**Output:**
```
Creating new Agentfile 'my-assistant'
✓ Initialized agentfile.yaml
✓ Created agentfile.natural.txt for natural language definition
✓ Added default state schema

Done! Edit agentfile.natural.txt to define your agent's purpose and behavior.
```

### `sentinel agentfile convert`

Converts a natural language description to a structured YAML Agentfile.

```bash
sentinel agentfile convert <path/to/agentfile.natural.txt>
```

**Flags:**
- `--endpoint, -e`: Override the model endpoint URL (default: http://localhost:11434)
- `--verbose, -v`: Enable verbose output

**Example:**
```bash
sentinel agentfile convert my-assistant/agentfile.natural.txt
```

**Output:**
```
Converting 'my-assistant/agentfile.natural.txt' to YAML using endpoint 'http://localhost:11434'...
✓ Successfully converted to 'my-assistant/agentfile.yaml'
```

## Agent Commands

Commands for running and managing agents.

### `sentinel agent run`

Runs an agent interactively.

```bash
sentinel agent run <agent-name>
```

**Flags:**
- `--endpoint, -e`: Override the model endpoint URL

**Example:**
```bash
sentinel agent run my-assistant
```

**Output:**
```
Running agent: my-assistant
Agent 'my-assistant' is ready. Type 'exit' to quit.

> Hello
Hello! How can I help you today?

> exit
```

## Registry Commands

Commands for interacting with the agent registry.

### `sentinel registry push`

Pushes an agent to the registry.

```bash
sentinel registry push <agent-name>
```

**Flags:**
- `--visibility, -v`: Set the visibility of the agent (`public` or `private`, default: `public`)

**Example:**
```bash
sentinel registry push my-assistant --visibility public
```

### `sentinel registry pull`

Pulls an agent from the registry.

```bash
sentinel registry pull <username/agent-name>
```

**Example:**
```bash
sentinel registry pull satishgonella/research-assistant
```

**Output:**
```
Pulling agent 'satishgonella/research-assistant' from registry...
Successfully pulled agent to 'research-assistant'
```

### `sentinel registry search`

Searches for agents in the registry.

```bash
sentinel registry search <query>
```

**Flags:**
- `--tags, -t`: Filter by tags (comma-separated)

**Example:**
```bash
sentinel registry search research --tags academic,papers
```

### `sentinel registry list`

Lists all agents in the registry.

```bash
sentinel registry list
```

**Output Example:**
```
Listing agents in registry...
Found 2 agents:

1. satishgonella/research-assistant@0.1.0
   Description: A research assistant that helps with academic papers
   Models: ollama/llama3
   Capabilities: summarization, question_answering

2. satishgonella/code-helper@0.2.0
   Description: Helps write and debug code
   Models: ollama/llama3
   Capabilities: code_generation, debugging
```

## Environment Variables

SentinelStacks respects the following environment variables:

- `SENTINEL_CONFIG`: Path to the configuration file (default: `~/.sentinelstacks/config.yaml`)
- `SENTINEL_REGISTRY`: Path to the local registry (default: `~/.sentinelstacks/registry`)
- `SENTINEL_MODEL_ENDPOINT`: Default model endpoint URL (default: `http://localhost:11434`)
