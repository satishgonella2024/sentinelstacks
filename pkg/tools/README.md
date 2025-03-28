# SentinelStacks Tools Package

This package provides the infrastructure for adding tool capabilities to SentinelStacks agents. Tools allow agents to perform actions such as calculations, fetching weather information, or any other custom functionality.

## Architecture

The tools package is built around the following core components:

1. **Tool Interface**: Defines the contract that all tools must implement
2. **ToolManager**: Manages the instantiation and execution of tools
3. **ToolRegistry**: Global registry for tool registration and creation
4. **Built-in Tools**: Pre-defined tools that are available out of the box

## Tool Interface

The `Tool` interface defines the contract for all tools:

```go
type Tool interface {
    // ID returns the unique identifier for the tool
    ID() string

    // Name returns a user-friendly name for the tool
    Name() string

    // Description returns a detailed description of what the tool does
    Description() string

    // Version returns the semantic version of the tool
    Version() string

    // ParameterSchema returns the JSON schema for the tool's parameters
    ParameterSchema() map[string]interface{}

    // Execute runs the tool with the provided parameters and returns the result
    Execute(params map[string]interface{}) (interface{}, error)
}
```

## Built-in Tools

### Calculator Tool

A simple calculator that can perform basic arithmetic operations:
- Addition
- Subtraction
- Multiplication
- Division
- Power
- Square root

Example usage in agent response:
```
To calculate 5 + 3, I'll use the calculator tool: {{tool:calculator,operation:add,a:5,b:3}}
```

### Weather Tool

A tool for fetching current weather information for a location:
- Temperature
- Weather conditions
- Humidity
- Wind speed

Example usage in agent response:
```
Let me check the weather in London: {{tool:weather,location:London,units:metric}}
```

## Creating Custom Tools

To create a custom tool:

1. Implement the `Tool` interface
2. Register the tool with the global registry
3. Update the agent's Agentfile to include the tool

Example:

```go
// Define your custom tool
type MyCustomTool struct{}

// Implement the Tool interface methods
func (t *MyCustomTool) ID() string { return "my-custom-tool" }
func (t *MyCustomTool) Name() string { return "My Custom Tool" }
func (t *MyCustomTool) Description() string { return "Description of my tool" }
func (t *MyCustomTool) Version() string { return "0.1.0" }
func (t *MyCustomTool) ParameterSchema() map[string]interface{} { ... }
func (t *MyCustomTool) Execute(params map[string]interface{}) (interface{}, error) { ... }

// Register your tool with the registry
func init() {
    registry := tools.GetToolRegistry()
    registry.RegisterFactory("my-custom-tool", func() tools.Tool {
        return &MyCustomTool{}
    })
}
```

Then add the tool to your agent's configuration:

```yaml
tools:
  - id: my-custom-tool
    version: "0.1.0"
```

## Future Development

Planned enhancements for the tools package:

1. **Tool Permissions**: Fine-grained control over what tools can do
2. **Tool Chaining**: Allowing tools to call other tools
3. **Tool Versioning**: Better support for tool versioning and compatibility
4. **Web Tools**: Tools for web searches, API calls, etc.
5. **Document Tools**: Tools for document analysis, summarization, etc.
