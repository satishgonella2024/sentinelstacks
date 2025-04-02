package stack

// StackSpec defines the structure of a multi-agent stack
type StackSpec struct {
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Version     string                 `json:"version" yaml:"version"`
	Agents      []StackAgentSpec       `json:"agents" yaml:"agents"`
	Networks    []string               `json:"networks" yaml:"networks"`
	Volumes     []string               `json:"volumes" yaml:"volumes"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
}

// StackAgentSpec defines an individual agent within a stack
type StackAgentSpec struct {
	ID        string                 `json:"id" yaml:"id"`
	Uses      string                 `json:"uses" yaml:"uses"`
	InputFrom []string               `json:"inputFrom" yaml:"inputFrom"`
	InputKey  string                 `json:"inputKey" yaml:"inputKey"`
	OutputKey string                 `json:"outputKey" yaml:"outputKey"`
	Params    map[string]interface{} `json:"params" yaml:"params"`
	Depends   []string               `json:"depends" yaml:"depends"`
}

// AgentState represents the current state of an agent in the execution flow
type AgentState struct {
	ID           string
	Status       AgentStatus
	Inputs       map[string]interface{}
	Outputs      map[string]interface{}
	ErrorMessage string
	StartTime    int64
	EndTime      int64
}

// AgentStatus represents the execution status of an agent
type AgentStatus string

const (
	// AgentStatusPending indicates the agent is waiting to be executed
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

// StackExecutionSummary provides a summary of the stack execution
type StackExecutionSummary struct {
	StackName      string
	StartTime      int64
	EndTime        int64
	TotalAgents    int
	CompletedCount int
	FailedCount    int
	BlockedCount   int
	AgentStates    map[string]AgentState
}
