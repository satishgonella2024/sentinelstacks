// Package types defines common type definitions used across packages
package types

import (
	"context"
	"time"
)

// StackType defines the type of stack
type StackType string

const (
	// StackTypeDefault is the default stack type
	StackTypeDefault StackType = "default"

	// StackTypeAgent is a stack that represents a single agent
	StackTypeAgent StackType = "agent"

	// StackTypeWorkflow is a stack that represents a workflow
	StackTypeWorkflow StackType = "workflow"
)

// StackAgentSpec defines an agent in a stack
type StackAgentSpec struct {
	// ID is the unique identifier for the agent
	ID string

	// Uses specifies the agent implementation to use
	Uses string

	// Depends specifies agent IDs this agent depends on
	Depends []string

	// InputFrom specifies agent IDs to take input from
	InputFrom []string

	// With specifies configuration parameters for the agent
	With map[string]interface{}
}

// StackSpec defines a stack of agents
type StackSpec struct {
	// Name is the name of the stack
	Name string

	// Description describes the purpose of the stack
	Description string

	// Version is the version of the stack
	Version string

	// Type is the type of stack
	Type StackType

	// Agents are the agents in the stack
	Agents []StackAgentSpec
}

// AgentStatus represents the status of an agent during execution
type AgentStatus string

const (
	// AgentStatusPending indicates the agent is waiting to execute
	AgentStatusPending AgentStatus = "pending"

	// AgentStatusRunning indicates the agent is currently executing
	AgentStatusRunning AgentStatus = "running"

	// AgentStatusCompleted indicates the agent has completed successfully
	AgentStatusCompleted AgentStatus = "completed"

	// AgentStatusFailed indicates the agent execution failed
	AgentStatusFailed AgentStatus = "failed"

	// AgentStatusBlocked indicates the agent is blocked on dependencies
	AgentStatusBlocked AgentStatus = "blocked"
)

// AgentState represents the state of an agent during execution
type AgentState struct {
	// ID is the agent identifier
	ID string

	// Status is the current execution status
	Status AgentStatus

	// ErrorMessage contains error details if status is failed
	ErrorMessage string

	// StartTime is when the agent started execution
	StartTime time.Time

	// EndTime is when the agent completed execution
	EndTime time.Time

	// Inputs are the inputs provided to the agent
	Inputs map[string]interface{}

	// Outputs are the outputs produced by the agent
	Outputs map[string]interface{}
}

// StackExecutionSummary provides a summary of stack execution
type StackExecutionSummary struct {
	// StackName is the name of the stack
	StackName string

	// ExecutionID is a unique identifier for this execution
	ExecutionID string

	// StartTime is when execution started
	StartTime time.Time

	// EndTime is when execution completed
	EndTime time.Time

	// TotalAgents is the total number of agents in the stack
	TotalAgents int

	// CompletedCount is the number of completed agents
	CompletedCount int

	// FailedCount is the number of failed agents
	FailedCount int

	// BlockedCount is the number of blocked agents
	BlockedCount int

	// AgentStates contains the state of each agent
	AgentStates map[string]*AgentState
}

// StackStatus represents the current status of a stack
type StackStatus string

const (
	// StackStatusReady indicates the stack is ready to be executed
	StackStatusReady StackStatus = "ready"

	// StackStatusRunning indicates the stack is currently running
	StackStatusRunning StackStatus = "running"

	// StackStatusSucceeded indicates the stack executed successfully
	StackStatusSucceeded StackStatus = "succeeded"

	// StackStatusFailed indicates the stack execution failed
	StackStatusFailed StackStatus = "failed"

	// StackStatusCancelled indicates the stack execution was cancelled
	StackStatusCancelled StackStatus = "cancelled"
)

// ExecutionSummary provides a summary of a stack execution
type ExecutionSummary struct {
	// ExecutionID is the unique identifier for this execution
	ExecutionID string

	// StartTime is when the execution started
	StartTime time.Time

	// EndTime is when the execution completed
	EndTime time.Time

	// Status is the final status of the execution
	Status StackStatus

	// CompletedCount is the number of agents that completed successfully
	CompletedCount int

	// FailedCount is the number of agents that failed
	FailedCount int

	// BlockedCount is the number of agents that were blocked
	BlockedCount int
}

// StateManager manages the state of agents during stack execution
type StateManager interface {
	// InitializeAgents sets up initial state for all agents
	InitializeAgents(agentIDs []string)

	// Get retrieves a value from an agent's state by key
	Get(agentID string, key string) (interface{}, error)

	// Set stores a value in an agent's state
	Set(agentID string, key string, value interface{}) error

	// GetAll returns all state values for an agent
	GetAll(agentID string) (map[string]interface{}, error)

	// UpdateAgentStatus updates the execution status of an agent
	UpdateAgentStatus(agentID string, status AgentStatus) error

	// UpdateAgentErrorMessage updates the error message for an agent
	UpdateAgentErrorMessage(agentID string, errorMessage string) error

	// GetStackSummary returns a summary of the current stack execution
	GetStackSummary() *StackExecutionSummary
}

// DAGNode represents a node in the DAG
type DAGNode struct {
	ID           string
	Dependencies []string
	Dependents   []string
}

// DAG represents a directed acyclic graph for execution order
type DAG interface {
	// GetNode returns the node with the given ID
	GetNode(id string) (*DAGNode, bool)

	// GetNodes returns all nodes in the DAG
	GetNodes() []*DAGNode

	// GetRoots returns all root nodes (nodes with no dependencies)
	GetRoots() []*DAGNode

	// GetLeaves returns all leaf nodes (nodes with no dependents)
	GetLeaves() []*DAGNode

	// TopologicalSort returns nodes in topological order
	TopologicalSort() ([]string, error)

	// HasCycle checks if the DAG has a cycle
	HasCycle() bool
}

// StackEngine executes a stack of agents
type StackEngine interface {
	// Execute executes the stack with the given inputs
	Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)

	// Stop stops the execution of the stack
	Stop(ctx context.Context) error

	// GetState returns the current state of the stack
	GetState(ctx context.Context) (StackExecutionSummary, error)
}

// StackExecutor defines the interface for executing stacks
type StackExecutor interface {
	// Execute runs the stack with the given context
	Execute(ctx context.Context) error

	// Stop cancels the execution of the stack
	Stop()

	// GetState returns the current state of the stack execution
	GetState() map[string]interface{}

	// GetAgentState returns the current state of an agent
	GetAgentState(agentID string) (map[string]interface{}, error)
}

// StackManager provides operations for managing stacks
type StackManager interface {
	// CreateStack creates a new stack from a specification
	CreateStack(spec StackSpec) (string, error)

	// GetStack gets a stack by ID
	GetStack(id string) (*StackSpec, error)

	// ListStacks lists all stacks
	ListStacks() ([]StackSpec, error)

	// DeleteStack deletes a stack
	DeleteStack(id string) error

	// UpdateStack updates a stack
	UpdateStack(id string, spec StackSpec) error

	// ExecuteStack executes a stack and returns an executor
	ExecuteStack(id string, options map[string]interface{}) (StackExecutor, error)
}
