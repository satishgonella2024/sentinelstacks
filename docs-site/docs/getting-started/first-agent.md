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

> Can you explain how the handler function works?

In Go's HTTP server, the handler function is responsible for processing HTTP requests and generating responses. Let's break down how it works:

1. The function signature: `func handler(w http.ResponseWriter, r *http.Request)`
   - It takes two parameters:
     - `w http.ResponseWriter`: An interface used to send the response back to the client
     - `r *http.Request`: Contains all the information about the incoming request

2. Inside the function:
   - `fmt.Fprintf(w, "Hello, World!")` writes the text "Hello, World!" to the response
   - The ResponseWriter interface implements methods for writing HTTP responses

3. How it gets called:
   - `http.HandleFunc("/", handler)` registers this function to handle requests to the root path "/"
   - When a request comes in, the server calls this function with the appropriate parameters

Handler functions can do much more, like reading request data, setting headers, or returning different status codes.
```

## Sharing Your Agent

Once you've created and tested your agent, you can share it with others through the registry:

```bash
sentinel registry push my-assistant
```

This will:
- Package your agent files
- Upload them to the registry
- Make the agent available for others to discover and use

## Pulling an Agent from the Registry

To use an agent someone else has created:

```bash
sentinel registry pull username/agent-name
```

This will download the agent to your local machine, allowing you to run it.

## Next Steps

- [Learn about the Agentfile specification](../user-guide/agentfile-spec.md)
- [Explore the registry to discover agents](../user-guide/registry.md)
- [Create more complex agents with tools](../user-guide/tools.md)
