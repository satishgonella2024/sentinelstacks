# Creating Your First Agent

This guide will walk you through creating your first AI agent with SentinelStacks.

## Creating a New Agent

1. Use the CLI to create a new agent:

```bash
sentinel agentfile create --name my-assistant
```

This will:
- Create a directory named `my-assistant`
- Initialize a basic agent configuration
- Create a natural language description file

## Understanding the Files

The command creates several files:

- **agentfile.natural.txt**: A natural language description of your agent
- **agentfile.yaml**: The structured YAML configuration generated from the description
- **agent.state.json**: An empty state file for persisting agent memory

## Editing the Natural Language Description

Open `my-assistant/agentfile.natural.txt` in your favorite text editor and describe what you want your agent to do:

```
This agent helps answer questions about programming languages.
It should focus on Python, JavaScript, and Go.
It should provide code examples when relevant.
It should be friendly and clear in its explanations.
It should use Llama3 as its model.
```

## Converting to YAML

After editing the natural language description, convert it to a structured YAML file:

```bash
sentinel agentfile convert my-assistant/agentfile.natural.txt
```

This will update the `agentfile.yaml` file with a configuration based on your description.

## Running Your Agent

To start interacting with your agent:

```bash
sentinel agent run my-assistant
```

This will start an interactive session where you can chat with your agent.

## Example Conversation

```
> sentinel agent run my-assistant
Running agent: my-assistant
Connected to Llama3 via Ollama

> How do I create a simple HTTP server in Go?

To create a simple HTTP server in Go, you can use the standard library's "net/http" package. Here's a basic example:

```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on port 8080...")
    http.ListenAndServe(":8080", nil)
}
```

Save this code to a file named `server.go` and run it with `go run server.go`. Then open your browser to http://localhost:8080.

> 
```

## Next Steps

- Explore more [advanced agent capabilities](../user-guide/agentfile-spec.md)
- Learn how to [share your agent](../user-guide/registry.md) with others
- Create [multi-agent workflows](../user-guide/stacks.md)
