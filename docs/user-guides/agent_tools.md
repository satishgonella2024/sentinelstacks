# Agent Tools

SentinelStacks provides a powerful tools system that enables agents to interact with external resources and perform actions beyond text generation. This guide explains how to use and configure tools for your agents.

## Overview

Tools in SentinelStacks allow agents to:

- Access and manipulate files
- Search the web for information
- Call external APIs
- Execute specific operations with proper permissions

Each tool has defined parameters, requirements, and associated permissions to ensure security.

## Available Tools

SentinelStacks includes several built-in tools:

### File Tools

| Tool | Description | Required Permission |
|------|-------------|---------------------|
| `file/read` | Read the contents of a file | `file` |
| `file/write` | Write content to a file | `file` |
| `file/list` | List files in a directory | `file` |

### Web Tools

| Tool | Description | Required Permission |
|------|-------------|---------------------|
| `web/search` | Search the web for information | `network` |

## Configuring Tools in Sentinelfiles

To enable tools for an agent, specify them in your Sentinelfile:

```yaml
name: ResearchAssistant
description: An agent that can research topics and save results
baseModel: claude-3-sonnet-20240229

tools:
  - file/read
  - file/write
  - web/search

toolSettings:
  web/search:
    default_results: 5
    safe_search: true
  file/write:
    allowed_extensions: [".txt", ".md", ".json"]
```

### Tool Settings

Each tool can have specific settings that control its behavior:

#### web/search

- `default_results`: Number of results to return (default: 3)
- `safe_search`: Enable safe search filtering (default: true)
- `timeout`: Maximum time in seconds for the search (default: 10)

#### file/read

- `max_size`: Maximum file size in bytes that can be read (default: 1048576)
- `allowed_extensions`: List of allowed file extensions (default: all)

#### file/write

- `max_size`: Maximum file size in bytes that can be written (default: 1048576)
- `allowed_extensions`: List of allowed file extensions (default: all)
- `overwrite`: Allow overwriting existing files (default: false)

#### file/list

- `recursive`: List files recursively (default: false)
- `include_hidden`: Include hidden files in the listing (default: false)

## Managing Tool Permissions

Tools require specific permissions that must be granted to agents for security reasons. The available permissions are:

- `none`: No special permissions required
- `file`: File system access permissions
- `network`: Network access permissions
- `shell`: Shell access permissions
- `api`: API access permissions
- `all`: All permissions (use with caution)

### Permission Management

You can manage tool permissions using the CLI:

```bash
# List permissions for an agent
sentinel tools perms research-agent

# Grant a permission to an agent
sentinel tools grant research-agent file

# Revoke a permission from an agent
sentinel tools revoke research-agent network
```

## Using Tools in Agents

When an agent has tools enabled, it can use them during conversations. The agent will automatically learn about available tools from their descriptions.

### Example Conversation

```
User: I need to research quantum computing and save the results to a file.

Agent: I'll help you research quantum computing and save the information. Let me search for some information first.

[Tool: web/search]
Parameters: {"query": "quantum computing basics recent developments", "num_results": 5}
Result: [
  {
    "title": "Quantum Computing: An Overview of Principles and Recent Advances",
    "url": "https://example.com/quantum-overview",
    "snippet": "Quantum computing leverages quantum mechanical phenomena like superposition and entanglement..."
  },
  ...
]

Agent: I've found some information about quantum computing. Here are the key points:
1. Quantum computers use qubits instead of classical bits
2. They leverage superposition and entanglement
3. Recent advances include error correction improvements
4. IBM and Google have made significant progress in the field

Would you like me to save this information to a file?

User: Yes, please save it to quantum_research.md.

Agent: I'll save the information to quantum_research.md.

[Tool: file/write]
Parameters: {"path": "quantum_research.md", "content": "# Quantum Computing Research\n\n## Basic Principles\n\n- Quantum computers use qubits instead of classical bits\n- They leverage superposition and entanglement\n\n## Recent Developments\n\n- Error correction improvements\n- IBM and Google have made significant progress\n- Quantum advantage demonstrated for specific problems\n"}
Result: {"success": true, "bytes_written": 312}

Agent: I've saved the quantum computing research to quantum_research.md. The file includes information about basic principles and recent developments. Is there anything else you'd like me to add to the file or any other research topics you're interested in?
```

## Creating Custom Tools

For advanced users, SentinelStacks supports creating custom tools. This requires:

1. Implementing the `Tool` interface
2. Registering the tool with the global registry
3. Ensuring proper permission checking

For detailed information on developing custom tools, refer to the [Developer Guide](/development/custom_tools.md).

## Security Considerations

When using tools, keep these security practices in mind:

1. **Grant minimal permissions**: Only grant permissions that an agent absolutely needs
2. **Restrict file access**: Use `allowed_extensions` and `max_size` to limit file operations
3. **Monitor tool usage**: Keep track of how agents use tools, especially for production deployments
4. **Review agent outputs**: Always review agent-generated content before executing critical operations
5. **Use sandboxing**: For sensitive environments, consider additional sandboxing mechanisms

## Troubleshooting

### Agent Cannot Access Tools

If your agent cannot access tools, check:

1. That the tools are correctly specified in the Sentinelfile
2. That the required permissions have been granted
3. That the agent has been restarted after permission changes

### Tool Execution Fails

If a tool execution fails:

1. Check the error message for specific issues
2. Verify parameter formats and values
3. Ensure any resource the tool needs to access exists and is accessible
4. Confirm that network connectivity is available for network-based tools

## Best Practices

1. **Design purpose-specific agents**: Create agents with specific purposes and grant only the tools they need
2. **Structure tool parameters carefully**: Provide clear and specific parameters to tools
3. **Validate agent outputs**: Especially for file writing or other persistent actions
4. **Use tool settings to restrict scope**: Limit search results, file sizes, and allowed operations
5. **Combine tools effectively**: Design agents that combine multiple tools for complex workflows 