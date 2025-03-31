# Sentinel CLI Implementation Plan

This document outlines the implementation plan for the Sentinel CLI, which is the primary interface for users to interact with the SentinelStacks system.

## Core Architecture

The CLI will be built using the following structure:

```
sentinel/
├── cmd/            # Command implementations
│   ├── build/      # Build command
│   ├── run/        # Run command
│   ├── push/       # Push command
│   ├── pull/       # Pull command
│   └── ...         # Other commands
├── internal/       # Internal packages
│   ├── config/     # Configuration management
│   ├── parser/     # Sentinelfile parser
│   ├── runtime/    # Agent runtime
│   ├── registry/   # Registry client
│   └── shim/       # LLM provider abstraction
├── pkg/            # Public packages for extending functionality
│   ├── agent/      # Agent definition models
│   ├── tools/      # Tool integration framework
│   └── api/        # API client for services
└── main.go         # Entry point
```

## Command Structure

The CLI will support the following commands:

### Basic Commands

- `sentinel init`: Initialize a new Sentinelfile
- `sentinel build`: Build a Sentinel Image from a Sentinelfile
- `sentinel run`: Run a Sentinel Agent from an image
- `sentinel ps`: List running agents
- `sentinel stop`: Stop a running agent
- `sentinel logs`: View logs from a running agent

### Registry Commands

- `sentinel login`: Authenticate with a Sentinel Registry
- `sentinel logout`: Log out from a Sentinel Registry
- `sentinel push`: Push a Sentinel Image to a registry
- `sentinel pull`: Pull a Sentinel Image from a registry
- `sentinel search`: Search for agents in a registry

### Management Commands

- `sentinel images`: List local Sentinel Images
- `sentinel prune`: Remove unused images and stopped agents
- `sentinel inspect`: Display detailed information about an image
- `sentinel validate`: Validate a Sentinelfile

### Configuration Commands

- `sentinel config set`: Set configuration values
- `sentinel config get`: Get configuration values
- `sentinel config list`: List all configuration values

## Implementation Phases

### Phase 1: Core Framework (Weeks 1-2)

1. Set up the basic CLI structure using Cobra
2. Create command placeholders
3. Implement configuration management
4. Add logging and error handling

### Phase 2: Sentinelfile & Building (Weeks 3-4)

1. Implement basic Sentinelfile parser
2. Create agent definition model
3. Build the `init` command with templates
4. Implement the `build` command for local images

### Phase 3: Runtime (Weeks 5-6)

1. Develop the agent runtime
2. Implement the `run` command
3. Add state management (basic)
4. Create process management for agents

### Phase 4: Registry (Weeks 7-8)

1. Design registry API communication
2. Implement authentication
3. Build push/pull functionality
4. Add registry search

### Phase 5: Management & Monitoring (Weeks 9-10)

1. Implement logging system
2. Add agent monitoring
3. Create inspection and validation tools
4. Build prune and cleanup functionality

## CLI Design Principles

1. **Consistent Patterns**: Follow Docker's command patterns for familiarity
2. **Progressive Disclosure**: Basic commands should be simple, with advanced options available
3. **Helpful Errors**: Provide clear error messages with suggestions for resolution
4. **Fast Startup**: Minimize initial load time for the CLI
5. **Extensibility**: Allow for plugins and extensions

## Code Examples

### Main Entry Point

```go
// main.go
package main

import (
	"os"

	"github.com/sentinelstacks/cli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

### Root Command

```go
// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "SentinelStacks - Agent Management System",
	Long: `SentinelStacks is a complete system for creating, managing,
and distributing AI agents using natural language definitions.`,
}

func init() {
	// Add global flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	
	// Add commands
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(pushCmd)
	RootCmd.AddCommand(pullCmd)
	// Add other commands...
}
```

### Build Command

```go
// cmd/build.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sentinelstacks/cli/internal/parser"
	"github.com/sentinelstacks/cli/internal/builder"
)

var buildCmd = &cobra.Command{
	Use:   "build [options] -t name:tag",
	Short: "Build a Sentinel Image from a Sentinelfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse flags
		tag, _ := cmd.Flags().GetString("tag")
		file, _ := cmd.Flags().GetString("file")
		
		// Process Sentinelfile
		parser := parser.NewSentinelfileParser()
		def, err := parser.ParseFile(file)
		if err != nil {
			return err
		}
		
		// Build image
		builder := builder.NewImageBuilder()
		return builder.Build(def, tag)
	},
}

func init() {
	buildCmd.Flags().StringP("tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringP("file", "f", "Sentinelfile", "Path to Sentinelfile")
	buildCmd.MarkFlagRequired("tag")
}
```

## Testing Strategy

1. **Unit Tests**: For parsing, building, and internal functions
2. **Integration Tests**: For command execution and API communication
3. **End-to-End Tests**: For complete workflows
4. **Mocks**: For LLM API and registry communication

## Performance Considerations

1. Cache parsed Sentinelfiles when appropriate
2. Use efficient serialization for agent definitions
3. Implement parallel operations where possible
4. Minimize API calls to LLM providers
