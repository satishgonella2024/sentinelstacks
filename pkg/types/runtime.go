// Package types defines common type definitions used across packages
package types

import (
	"context"
)

// RuntimeType defines the type of agent runtime
type RuntimeType string

const (
	// RuntimeTypeDirect executes agents directly using the LLM provider
	RuntimeTypeDirect RuntimeType = "direct"

	// RuntimeTypeCli executes agents using the sentinel CLI
	RuntimeTypeCli RuntimeType = "cli"
)

// AgentRuntime defines the interface for agent execution
type AgentRuntime interface {
	// Execute runs an agent with the provided inputs and returns its outputs
	Execute(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error)

	// Cleanup releases resources after execution
	Cleanup() error
}

// RuntimeFactory creates and configures agent runtimes
type RuntimeFactory interface {
	// CreateRuntime creates an agent runtime of the specified type
	CreateRuntime(runtimeType RuntimeType) (AgentRuntime, error)

	// DefaultRuntime creates the default agent runtime based on environment
	DefaultRuntime() (AgentRuntime, error)
}
