package stack

import (
	"context"
)

// StackExecutor defines the interface for executing stacks
type StackExecutor interface {
	// BuildExecutionGraph builds and validates the execution graph for the stack
	BuildExecutionGraph() (*DAG, error)
	
	// Execute runs the stack with the given context
	Execute(ctx context.Context) error
	
	// Stop cancels the execution of the stack
	Stop()
	
	// GetState returns the current state of the stack execution
	GetState() *StackExecutionSummary
	
	// GetAgentState returns the current state of an agent
	GetAgentState(agentID string) (map[string]interface{}, error)
	
	// ExportStackState exports the current state of the stack as JSON
	ExportStackState() ([]byte, error)
}

// StackManager provides operations for managing stacks
type StackManager interface {
	// CreateStack creates a new stack from a specification
	CreateStack(spec StackSpec) (string, error)
	
	// GetStack gets a stack by ID
	GetStack(id string) (*Stack, error)
	
	// GetStackByName gets a stack by name
	GetStackByName(name string) (*Stack, error)
	
	// ListStacks lists all stacks
	ListStacks() ([]*Stack, error)
	
	// DeleteStack deletes a stack
	DeleteStack(id string) error
	
	// UpdateStack updates a stack
	UpdateStack(id string, spec StackSpec) error
}

// Stack represents a stack of agents
type Stack struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Spec        StackSpec              `json:"spec"`
	CreatedAt   int64                  `json:"createdAt"`
	UpdatedAt   int64                  `json:"updatedAt"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// StackExecutionOptions defines options for stack execution
type StackExecutionOptions struct {
	// Input data to provide to the stack
	Input map[string]interface{}
	
	// Whether to enable verbose logging
	Verbose bool
	
	// Timeout for the execution
	Timeout int
	
	// RuntimeOptions contains options for the runtime
	RuntimeOptions map[string]interface{}
}
