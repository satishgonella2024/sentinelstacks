# Add Tool Support to SentinelStacks

This PR adds support for tools in SentinelStacks agents, allowing them to perform actions such as calculations and fetching weather information.

## Features

- **Tool Interface**: Created a standard interface that all tools must implement
- **Tool Manager**: Added a manager for executing tools and handling their results
- **Tool Registry**: Implemented a global registry for tool registration and creation
- **Built-in Tools**:
  - Calculator: For performing arithmetic operations
  - Weather: For fetching current weather information
- **Agent Runtime Integration**: Updated the agent runtime to support tool execution
- **Example Agent**: Created an example agent that demonstrates tool usage

## Implementation Details

- Tools use a simple command pattern in agent responses: `{{tool:tool_name,param1:value1,param2:value2}}`
- The agent runtime detects and executes tool calls, replacing them with the results
- Tools are registered globally but instantiated per agent
- Tool parameters are converted to appropriate types (number, boolean, string)

## Testing

- Created unit tests for tool components
- Manually tested the example agent with various tool calls
- Verified that tool results are correctly integrated into agent responses

## Usage Example

```yaml
# Agent configuration with tools
name: ToolHelper
version: "0.1.0"
description: "An AI assistant that can perform calculations and get weather information"
tools:
  - id: calculator
    version: "0.1.0"
  - id: weather
    version: "0.1.0"
```

Then in the agent's responses:

```
To calculate 5 + 3, I'll use the calculator: {{tool:calculator,operation:add,a:5,b:3}}

Let me check the weather in London: {{tool:weather,location:London,units:metric}}
```

## Future Enhancements

- Add more built-in tools (web search, document analysis, etc.)
- Implement tool permissions for security
- Support function calling in compatible models
- Add tool chaining and composition

## Related Issues

Closes #XX - Add tool support for agents

## Checklist

- [x] Code follows the project's coding style
- [x] Documentation has been updated
- [x] Tests have been added for new functionality
- [x] All tests pass
- [x] No new warnings were introduced
