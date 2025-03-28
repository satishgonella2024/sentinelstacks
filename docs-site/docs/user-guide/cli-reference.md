# CLI Reference

The SentinelStacks CLI (`sentinel`) provides a command-line interface for managing and interacting with agents.

## Global Commands

### `sentinel --help`

Shows general help information and lists available commands.

### `sentinel --version`

Shows version information for the CLI.

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

### `sentinel agentfile convert`

Converts a natural language description to a structured YAML Agentfile.

```bash
sentinel agentfile convert <path/to/agentfile.natural.txt>
```

**Example:**
```bash
sentinel agentfile convert my-assistant/agentfile.natural.txt
```

### `sentinel agentfile validate`

Validates an Agentfile against the schema.

```bash
sentinel agentfile validate <path/to/agentfile.yaml>
```

**Example:**
```bash
sentinel agentfile validate my-assistant/agentfile.yaml
```

## Agent Commands

Commands for running and managing agents.

### `sentinel agent run`

Runs an agent interactively.

```bash
sentinel agent run <agent-name>
```

**Flags:**
- `--model <model-name>`: Override the model specified in the Agentfile
- `--endpoint <url>`: Override the model endpoint (e.g., for Ollama)

**Example:**
```bash
sentinel agent run my-assistant --model llama3
```

### `sentinel agent list`

Lists all available agents.

```bash
sentinel agent list
```

## Registry Commands

Commands for interacting with the agent registry.

### `sentinel registry push`

Pushes an agent to the registry.

```bash
sentinel registry push <agent-name>
```

**Flags:**
- `--visibility <public|private>`: Set the visibility of the agent (default: public)
- `--version <version>`: Specify a version (default: current version in Agentfile)

**Example:**
```bash
sentinel registry push my-assistant --visibility public --version 1.0.0
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

### `sentinel registry search`

Searches for agents in the registry.

```bash
sentinel registry search <query>
```

**Flags:**
- `--tags <tags>`: Filter by tags (comma-separated)
- `--model <model>`: Filter by compatible model

**Example:**
```bash
sentinel registry search research --tags academic,papers --model llama3
```

## Stack Commands

Commands for managing multi-agent workflows.

### `sentinel stack create`

Creates a new stack of agents.

```bash
sentinel stack create --name <stack-name>
```

### `sentinel stack run`

Runs a stack of agents.

```bash
sentinel stack run <stack-name>
```

## Model Commands

Commands for managing model connections.

### `sentinel model list`

Lists available models.

```bash
sentinel model list
```

### `sentinel model add`

Adds a new model connection.

```bash
sentinel model add --provider <provider> --name <name> --endpoint <url>
```

**Example:**
```bash
sentinel model add --provider ollama --name llama3 --endpoint http://localhost:11434
```
