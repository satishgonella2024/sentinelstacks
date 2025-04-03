// Package stack provides core stack execution functionality
package stack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/runtime"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// Engine is a clean implementation of the stack execution engine
type Engine struct {
	spec      types.StackSpec
	dag       *DAG
	verbose   bool
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	runID     string
	isRunning bool
}

// NewEngine creates a new stack execution engine
func NewEngine(spec types.StackSpec, options ...EngineOption) (*Engine, error) {
	// Create DAG from spec
	dag, err := NewDAG(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build execution graph: %w", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Generate run ID
	runID := fmt.Sprintf("run-%d", time.Now().Unix())

	// Create engine with defaults
	engine := &Engine{
		spec:      spec,
		dag:       dag,
		ctx:       ctx,
		cancel:    cancel,
		runID:     runID,
		isRunning: false,
		verbose:   false,
	}

	// Apply options
	for _, option := range options {
		option(engine)
	}

	return engine, nil
}

// EngineOption defines a function that configures an Engine
type EngineOption func(*Engine)

// WithVerbose enables verbose logging for the engine
func WithVerbose(verbose bool) EngineOption {
	return func(e *Engine) {
		e.verbose = verbose
	}
}

// WithRunID sets a custom run ID for the engine
func WithRunID(runID string) EngineOption {
	return func(e *Engine) {
		e.runID = runID
	}
}

// ExecuteOptions defines options for executing a stack
type ExecuteOptions struct {
	Timeout     int
	Input       map[string]interface{}
	RuntimeType types.RuntimeType
}

// ExecuteOption defines a function that configures execution options
type ExecuteOption func(*ExecuteOptions)

// WithTimeout sets a timeout for execution
func WithTimeout(timeout int) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Timeout = timeout
	}
}

// WithInput sets input data for execution
func WithInput(input map[string]interface{}) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Input = input
	}
}

// WithRuntimeType sets the runtime type for execution
func WithRuntimeType(runtimeType types.RuntimeType) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.RuntimeType = runtimeType
	}
}

// Execute runs the stack with the provided options
func (e *Engine) Execute(ctx context.Context, options ...ExecuteOption) error {
	e.mu.Lock()
	if e.isRunning {
		e.mu.Unlock()
		return errors.New("stack is already running")
	}
	e.isRunning = true
	e.mu.Unlock()

	// Apply execution options
	execOptions := &ExecuteOptions{
		Timeout:     0,
		Input:       make(map[string]interface{}),
		RuntimeType: types.RuntimeTypeDirect,
	}

	for _, option := range options {
		option(execOptions)
	}

	// Set up context with cancellation
	var execCtx context.Context
	var cancel context.CancelFunc

	if execOptions.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, time.Duration(execOptions.Timeout)*time.Second)
	} else {
		execCtx, cancel = context.WithCancel(ctx)
	}
	e.cancel = cancel
	e.ctx = execCtx

	// Get execution order
	executionOrder, err := e.dag.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Create a map to track executed nodes
	executedNodes := make(map[string]bool)

	if e.verbose {
		log.Printf("Starting stack execution: %s (Run ID: %s)", e.spec.Name, e.runID)
		log.Printf("Execution order: %v", executionOrder)
	}

	// Create map to store agent outputs
	agentOutputs := make(map[string]map[string]interface{})

	// Execute agents in order
	for _, agentID := range executionOrder {
		// Check if execution was cancelled
		select {
		case <-execCtx.Done():
			return execCtx.Err()
		default:
			// Continue execution
		}

		// Get agent specification
		var agentSpec types.StackAgentSpec
		for _, agent := range e.spec.Agents {
			if agent.ID == agentID {
				agentSpec = agent
				break
			}
		}

		if e.verbose {
			log.Printf("Preparing to execute agent: %s", agentID)
		}

		// Collect inputs from dependencies
		inputs := make(map[string]interface{})

		// Add global inputs
		for k, v := range execOptions.Input {
			inputs[k] = v
		}

		// Add inputs from dependencies
		for _, inputFrom := range agentSpec.InputFrom {
			if outputs, ok := agentOutputs[inputFrom]; ok {
				for k, v := range outputs {
					inputs[k] = v
				}
			}
		}

		// Execute the agent
		outputs, err := e.executeAgent(execCtx, agentSpec, inputs, execOptions.RuntimeType)
		if err != nil {
			return fmt.Errorf("failed to execute agent %s: %w", agentID, err)
		}

		// Store outputs
		agentOutputs[agentID] = outputs

		// Mark as executed
		executedNodes[agentID] = true

		if e.verbose {
			log.Printf("Agent executed successfully: %s", agentID)
		}
	}

	if e.verbose {
		log.Printf("Stack execution completed: %s (Run ID: %s)", e.spec.Name, e.runID)
	}

	// Reset running state
	e.mu.Lock()
	e.isRunning = false
	e.mu.Unlock()

	return nil
}

// executeAgent runs a single agent and returns its outputs
func (e *Engine) executeAgent(ctx context.Context, agentSpec types.StackAgentSpec, inputs map[string]interface{}, runtimeType types.RuntimeType) (map[string]interface{}, error) {
	if e.verbose {
		log.Printf("Executing agent %s (uses: %s)", agentSpec.ID, agentSpec.Uses)
	}

	// Create agent runtime factory
	factory := runtime.NewRuntimeFactory(e.verbose)

	// Create the runtime
	agentRuntime, err := factory.CreateRuntime(runtimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent runtime: %w", err)
	}
	defer agentRuntime.Cleanup()

	// Execute the agent using the runtime
	outputs, err := agentRuntime.Execute(ctx, agentSpec, inputs)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	return outputs, nil
}

// Stop cancels the execution of the stack
func (e *Engine) Stop() {
	if e.cancel != nil {
		e.cancel()
	}
	e.isRunning = false
	if e.verbose {
		log.Printf("Stack execution stopped: %s (Run ID: %s)", e.spec.Name, e.runID)
	}
}
