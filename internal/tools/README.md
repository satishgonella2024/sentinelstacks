# SentinelStacks Tool Integration

This directory contains the implementation of the tool integration framework for SentinelStacks. Tools allow agents to perform actions outside their conversation context, such as accessing files, searching the web, or interacting with external APIs.

## Architecture

The tool integration system consists of the following key components:

1. **Tool Interface**: Defines the contract for all tools, including name, description, parameters, execution, and permission requirements.

2. **Registry**: A global registry for managing available tools.

3. **Permission System**: Controls which tools are available to specific agents.

4. **Tool Executor**: Handles the execution of tools with proper parameter validation and permission checks.

5. **LLM Integration**: Augments LLM inputs with tool descriptions and handles function calling.

## Available Tools

### File Tools

- `file/read`: Read the contents of a file from the file system.
- `file/write`: Write content to a file in the file system.
- `file/list`: List the contents of a directory in the file system.

### Web Tools

- `web/search`: Search for information on the web and return structured results.

## Using Tools in Sentinelfiles

To enable tools for an agent, specify them in the Sentinelfile:

```yaml
name: MyAgent
description: An agent that can use tools
baseModel: claude-3-opus-20240229

tools:
  - file/read
  - file/write
  - web/search

toolSettings:
  web/search:
    default_results: 5
    safe_search: true
```

## Permission System

Tools require specific permissions to be granted to agents. The available permissions are:

- `none`: No permissions required
- `file`: File system access permissions
- `network`: Network access permissions
- `shell`: Shell access permissions
- `api`: API access permissions
- `all`: All permissions

Permissions can be managed using the CLI:

```bash
# List permissions for an agent
sentinel tools perms [agent-id]

# Grant a permission to an agent
sentinel tools grant [agent-id] [permission]

# Revoke a permission from an agent
sentinel tools revoke [agent-id] [permission]
```

## Extending with New Tools

To add a new tool to SentinelStacks, follow these steps:

1. Create a new package for your tool category if it doesn't exist.
2. Implement the `Tool` interface.
3. Create a registration function to add your tool to the registry.
4. Update the initialization code to register your tool.

Example:

```go
package mytools

import (
	"context"
	"fmt"

	"github.com/sentinelstacks/sentinel/internal/tools"
)

type MyTool struct {
	tools.BaseTool
}

func NewMyTool() *MyTool {
	return &MyTool{
		BaseTool: tools.BaseTool{
			Name:        "my/tool",
			Description: "Description of my tool",
			Parameters: []tools.Parameter{
				{
					Name:        "param1",
					Type:        "string",
					Description: "Description of parameter",
					Required:    true,
				},
			},
			Permission: tools.PermissionNone,
		},
	}
}

func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Implementation
	return "result", nil
}

func RegisterMyTools() error {
	// Register tools
	registry := tools.GetRegistry()
	if err := registry.RegisterTool(NewMyTool()); err != nil {
		return err
	}
	return nil
}
```

## Future Enhancements

1. **Shell Tools**: Add tools for executing shell commands with proper sandboxing.
2. **API Tools**: Add tools for interacting with external APIs with configuration.
3. **Enhanced Permission System**: Implement more granular permissions for better security.
4. **Tool Usage Metrics**: Add monitoring and limits for tool usage.
5. **Improved Error Handling**: Enhance error handling and recovery for tool executions.
