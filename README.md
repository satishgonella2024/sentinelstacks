# Sentinel Stacks

An agent orchestration platform for AI systems

## Overview

Sentinel Stacks is a platform for orchestrating AI agents and creating complex AI workflows. It allows you to define stacks of agents that can be executed together to accomplish specific tasks.

## Architecture

The project is organized into several key components:

### Core Components

- **API Layer** (`pkg/api`): Provides a unified interface for interacting with all services
- **Stack Engine** (`pkg/stack`): Executes stacks of agents according to their dependencies
- **Memory System** (`pkg/memory`): Stores and retrieves data for agents and stacks
- **Registry** (`pkg/registry`): Manages packages and dependencies

### Type Definitions

- **Common Types** (`pkg/types`): Defines interfaces and types used across the codebase
- **Adapters** (`pkg/adapter`): Handles conversion between internal and public types

### Command Line Interfaces

- **Sentinel CLI** (`cmd/sentinelcli`): Main command line interface
- **API Example** (`cmd/api_example`): Example usage of the API

## Getting Started

### Prerequisites

- Go 1.20 or later
- SQLite (for persistent storage)

### Installation

```bash
# Clone the repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build the application
make build
```

### Running Examples

```bash
# Run the comprehensive example
make run-example-comprehensive

# Run specific service examples
make run-example-stack     # Stack service example
make run-example-memory    # Memory service example
make run-example-registry  # Registry service example
```

### Starting the CLI

```bash
# Start the Sentinel CLI
make run

# Start with custom data directory
./bin/sentinelcli -data /path/to/data
```

## Creating Stacks

A stack is defined as a collection of agents with dependencies and connections between them. Here's a simple example:

```go
// Create a stack specification
spec := types.StackSpec{
    Name:        "example-stack",
    Description: "A simple example stack",
    Version:     "1.0.0",
    Type:        types.StackTypeDefault,
    Agents: []types.StackAgentSpec{
        {
            ID:   "input-agent",
            Uses: "echo",
            With: map[string]interface{}{
                "message": "Hello from input agent",
            },
        },
        {
            ID:        "process-agent",
            Uses:      "transform",
            InputFrom: []string{"input-agent"},
            With: map[string]interface{}{
                "operation": "uppercase",
            },
        },
        {
            ID:        "output-agent",
            Uses:      "output",
            InputFrom: []string{"process-agent"},
            With: map[string]interface{}{
                "format": "json",
            },
        },
    },
}
```

## API Usage

The API provides a unified interface for accessing all services:

```go
// Initialize the API
config := api.APIConfig{
    StackConfig: api.StackServiceConfig{
        StoragePath: "data/stacks",
        Verbose:     true,
    },
    MemoryConfig: api.MemoryServiceConfig{
        StoragePath:         "data/memory",
        EmbeddingProvider:   "local",
        EmbeddingModel:      "local",
        EmbeddingDimensions: 1536,
    },
    RegistryConfig: api.RegistryServiceConfig{
        RegistryURL: "https://registry.example.com",
        CachePath:   "data/registry-cache",
    },
}

sentinel, err := api.NewAPI(config)
if err != nil {
    log.Fatalf("Failed to initialize API: %v", err)
}
defer sentinel.Close()

// Access services
stackService := sentinel.Stack()
memoryService := sentinel.Memory()
registryService := sentinel.Registry()

// Create a stack
stackID, err := stackService.CreateStack(ctx, spec)

// Execute a stack
inputs := map[string]interface{}{"message": "Hello, world!"}
results, err := stackService.ExecuteStack(ctx, stackID, inputs)
```

## Storage

Sentinel Stacks provides persistent storage for:

- **Stacks**: Definitions and execution history
- **Memory**: Key-value data and vector embeddings
- **Registry**: Cached packages and artifacts

By default, data is stored in the `data` directory:

```
data/
  ├── stacks/       # Stack definitions and state
  ├── memory/       # Memory storage
  ├── registry-cache/ # Downloaded packages
  └── exports/      # Exported stack definitions
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
