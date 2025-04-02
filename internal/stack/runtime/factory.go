package runtime

import (
	"fmt"
)

// RuntimeType defines the type of agent runtime to use
type RuntimeType string

const (
	// RuntimeTypeDirect executes agents directly using the LLM provider
	RuntimeTypeDirect RuntimeType = "direct"
	
	// RuntimeTypeCli executes agents using the sentinel CLI
	RuntimeTypeCli RuntimeType = "cli"
)

// RuntimeFactory creates and configures agent runtimes
type RuntimeFactory struct {
	logToConsole bool
}

// NewRuntimeFactory creates a new runtime factory
func NewRuntimeFactory(logToConsole bool) *RuntimeFactory {
	return &RuntimeFactory{
		logToConsole: logToConsole,
	}
}

// CreateRuntime creates an agent runtime of the specified type
func (f *RuntimeFactory) CreateRuntime(runtimeType RuntimeType) (AgentRuntime, error) {
	switch runtimeType {
	case RuntimeTypeDirect:
		return NewSimpleAgentRuntime(f.logToConsole)
	case RuntimeTypeCli:
		return NewSimpleAgentRuntime(f.logToConsole)
	default:
		return nil, fmt.Errorf("unsupported runtime type: %s", runtimeType)
	}
}

// DefaultRuntime creates the default agent runtime based on environment
func (f *RuntimeFactory) DefaultRuntime() (AgentRuntime, error) {
	// Use simple runtime for all cases in test mode
	return NewSimpleAgentRuntime(f.logToConsole)
}
