package stack

import (
	"github.com/satishgonella2024/sentinelstacks/internal/memory"
)

// EngineOption defines a function that configures a StackEngine
type EngineOption func(*StackEngine)

// WithVerbose enables verbose logging for the stack engine
func WithVerbose(verbose bool) EngineOption {
	return func(e *StackEngine) {
		e.verbose = verbose
	}
}

// WithMemoryFactory sets the memory factory for the stack engine
func WithMemoryFactory(factory memory.MemoryStoreFactory) EngineOption {
	return func(e *StackEngine) {
		e.memoryFactory = factory
	}
}

// WithStateManager sets a custom state manager for the stack engine
func WithStateManager(stateManager StateManager) EngineOption {
	return func(e *StackEngine) {
		e.stateManager = stateManager
	}
}

// ExecuteOption defines a function that configures execution options
type ExecuteOption func(*ExecuteOptions)

// ExecuteOptions defines options for executing a stack
type ExecuteOptions struct {
	Timeout         int
	Input           map[string]interface{}
	RuntimeOptions  map[string]interface{}
	RuntimeType     string
}

// WithTimeout sets the execution timeout in seconds
func WithTimeout(timeout int) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Timeout = timeout
	}
}

// WithInput sets the input data for the stack execution
func WithInput(input map[string]interface{}) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Input = input
	}
}

// WithRuntimeOptions sets runtime options for agent execution
func WithRuntimeOptions(options map[string]interface{}) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.RuntimeOptions = options
	}
}

// WithRuntimeType sets the runtime type for agent execution
func WithRuntimeType(runtimeType string) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.RuntimeType = runtimeType
	}
}
