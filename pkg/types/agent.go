// Package types contains common agent interfaces and types
package types

import (
	"context"
)

// AgentDefinition represents the definition of an agent
type AgentDefinition struct {
	Name           string
	Description    string
	Version        string
	BaseModel      string
	SystemPrompt   string
	PromptTemplate string
	MaxTokens      int
	Temperature    float64
	OutputFormat   string
	Config         map[string]interface{}
}

// Agent represents an executable agent
type Agent interface {
	// Execute executes the agent with the given inputs
	Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)

	// GetDefinition returns the agent's definition
	GetDefinition() *AgentDefinition
}

// RuntimeInterface defines the interface for agent runtimes
type RuntimeInterface interface {
	// Execute executes an agent with inputs and returns outputs
	Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)

	// GetName returns the name of the runtime
	GetName() string

	// Stop stops the execution of an agent
	Stop() error
}

// Provider defines the interface for LLM providers
type Provider interface {
	// Complete generates a completion given a prompt
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

	// GetName returns the name of the provider
	GetName() string

	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string
}
