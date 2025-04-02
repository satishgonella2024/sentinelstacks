package stack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
	stackmemory "github.com/satishgonella2024/sentinelstacks/internal/stack/memory"
	"github.com/satishgonella2024/sentinelstacks/internal/stack/runtime"
)

// StackEngine is responsible for executing a stack of agents
type StackEngine struct {
	spec         StackSpec
	stateManager StateManager
	dag          *DAG
	mu           sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
	runID        string
	isRunning    bool
	verbose      bool
	memoryFactory memory.MemoryStoreFactory
}

// NewStackEngine creates a new stack engine with the given specification
func NewStackEngine(spec StackSpec, options ...EngineOption) (*StackEngine, error) {
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
	engine := &StackEngine{
		spec:         spec,
		dag:          dag,
		ctx:          ctx,
		cancel:       cancel,
		runID:        runID,
		isRunning:    false,
		verbose:      false,
		memoryFactory: nil,
	}
	
	// Apply options
	for _, option := range options {
		option(engine)
	}
	
	// Create default memory factory if none provided
	if engine.memoryFactory == nil {
		factory, err := memory.NewMemoryStoreFactory("")
		if err != nil {
			return nil, fmt.Errorf("failed to create memory factory: %w", err)
		}
		engine.memoryFactory = factory
	}

	// Initialize agents in state manager
	agentIDs := make([]string, 0, len(spec.Agents))
	for _, agent := range spec.Agents {
		agentIDs = append(agentIDs, agent.ID)
	}
	
	// Create new persistent state manager using memory system
	stateManager, err := stackmemory.NewPersistentStateManager(
		ctx,
		spec.Name,
		engine.memoryFactory,
		spec.Name,
		runID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create state manager: %w", err)
	}
	
	// Initialize state
	stateManager.InitializeAgents(agentIDs)
	
	// Set state manager
	engine.stateManager = stateManager
	
	return engine, nil
}

// BuildExecutionGraph validates and returns the DAG for the stack
func (e *StackEngine) BuildExecutionGraph() (*DAG, error) {
	return e.dag, nil
}

// Execute runs the stack with provided options
func (e *StackEngine) Execute(ctx context.Context, options ...ExecuteOption) error {
	e.mu.Lock()
	if e.isRunning {
		e.mu.Unlock()
		return errors.New("stack is already running")
	}
	e.isRunning = true
	e.mu.Unlock()

	// Apply execution options
	execOptions := &ExecuteOptions{
		Timeout:        0,
		Input:          make(map[string]interface{}),
		RuntimeOptions: make(map[string]interface{}),
		RuntimeType:    "direct",
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

	// Execute agents in order
	for _, agentID := range executionOrder {
		select {
		case <-execCtx.Done():
			if e.verbose {
				log.Printf("Stack execution cancelled: %s", e.runID)
			}
			return execCtx.Err()
		default:
			// Continue execution
		}

		// Get agent spec
		var agentSpec StackAgentSpec
		for _, a := range e.spec.Agents {
			if a.ID == agentID {
				agentSpec = a
				break
			}
		}

		// Set agent status to running
		if err := e.stateManager.UpdateAgentStatus(agentID, AgentStatusRunning); err != nil {
			if e.verbose {
				log.Printf("Error updating agent status: %v", err)
			}
		}

		if e.verbose {
			log.Printf("Executing agent: %s", agentID)
		}
		
		// Collect inputs from dependencies
		inputs, err := e.collectInputs(agentSpec, execOptions.Input, executedNodes)
		if err != nil {
			if e.verbose {
				log.Printf("Error collecting inputs for agent %s: %v", agentID, err)
			}
			e.stateManager.UpdateAgentStatus(agentID, AgentStatusFailed)
			e.stateManager.UpdateAgentErrorMessage(agentID, fmt.Sprintf("Failed to collect inputs: %v", err))
			continue
		}

		// Set agent inputs
		if err := e.stateManager.Set(agentID, "input", inputs); err != nil {
			if e.verbose {
				log.Printf("Error setting agent inputs: %v", err)
			}
		}

		// Execute agent
		outputs, err := e.executeAgent(execCtx, agentSpec, inputs, execOptions)
		if err != nil {
			if e.verbose {
				log.Printf("Agent %s execution failed: %v", agentID, err)
			}
			e.stateManager.UpdateAgentStatus(agentID, AgentStatusFailed)
			e.stateManager.UpdateAgentErrorMessage(agentID, fmt.Sprintf("Execution failed: %v", err))
			continue
		}

		// Set agent outputs
		if err := e.stateManager.Set(agentID, "output", outputs); err != nil {
			if e.verbose {
				log.Printf("Error setting agent outputs: %v", err)
			}
		}

		// Mark agent as completed
		e.stateManager.UpdateAgentStatus(agentID, AgentStatusCompleted)
		executedNodes[agentID] = true
		
		if e.verbose {
			log.Printf("Agent completed: %s", agentID)
		}
	}

	// Check for any agents that didn't execute
	summary := e.stateManager.GetStackSummary()
	if summary.CompletedCount != summary.TotalAgents {
		if e.verbose {
			log.Printf("Stack execution completed with errors: %d/%d agents completed", 
				summary.CompletedCount, summary.TotalAgents)
		}
		return fmt.Errorf("stack execution completed with errors: %d/%d agents completed", 
			summary.CompletedCount, summary.TotalAgents)
	}

	if e.verbose {
		log.Printf("Stack execution completed successfully: %s (Run ID: %s)", e.spec.Name, e.runID)
	}
	e.isRunning = false
	
	return nil
}

// collectInputs gathers inputs for an agent from its dependencies
func (e *StackEngine) collectInputs(agentSpec StackAgentSpec, initialInput map[string]interface{}, executedNodes map[string]bool) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})
	
	// Add initial inputs
	if initialInput != nil {
		for k, v := range initialInput {
			inputs[k] = v
		}
	}
	
	// Add parameters to inputs
	if agentSpec.Params != nil {
		for k, v := range agentSpec.Params {
			inputs[k] = v
		}
	}
	
	// Collect inputs from inputFrom dependencies
	for _, inputFrom := range agentSpec.InputFrom {
		if inputFrom == "" {
			continue
		}
		
		// Skip dependencies that haven't executed
		if !executedNodes[inputFrom] {
			return nil, fmt.Errorf("dependency %s has not executed yet", inputFrom)
		}
		
		// Get output from dependency
		var output interface{}
		var err error
		
		if agentSpec.InputKey != "" {
			// If inputKey is specified, get that specific key
			output, err = e.stateManager.Get(inputFrom, agentSpec.InputKey)
		} else {
			// Otherwise get all outputs
			agentState, err := e.stateManager.GetAll(inputFrom)
			if err != nil {
				return nil, fmt.Errorf("failed to get state for agent %s: %w", inputFrom, err)
			}
			
			output = agentState["output"]
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to get output from agent %s: %w", inputFrom, err)
		}
		
		// Add to inputs based on source agent ID
		inputs[inputFrom] = output
	}
	
	return inputs, nil
}

// executeAgent runs a single agent with the given inputs
func (e *StackEngine) executeAgent(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}, options *ExecuteOptions) (map[string]interface{}, error) {
	if e.verbose {
		log.Printf("Executing agent %s (uses: %s)", agentSpec.ID, agentSpec.Uses)
	}
	
	// Create agent runtime factory
	factory := runtime.NewRuntimeFactory(e.verbose)
	
	// Create the runtime based on options
	var agentRuntime runtime.AgentRuntime
	var err error
	
	switch options.RuntimeType {
	case "direct":
		agentRuntime, err = factory.CreateRuntime(runtime.RuntimeTypeDirect)
	case "cli":
		agentRuntime, err = factory.CreateRuntime(runtime.RuntimeTypeCli)
	default:
		agentRuntime, err = factory.DefaultRuntime()
	}
	
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
func (e *StackEngine) Stop() {
	if e.cancel != nil {
		e.cancel()
	}
	e.isRunning = false
	log.Printf("Stack execution stopped: %s (Run ID: %s)", e.spec.Name, e.runID)
}

// GetState returns the current state of the stack execution
func (e *StackEngine) GetState() *StackExecutionSummary {
	return e.stateManager.GetStackSummary()
}

// GetAgentState returns the current state of an agent
func (e *StackEngine) GetAgentState(agentID string) (map[string]interface{}, error) {
	return e.stateManager.GetAll(agentID)
}

// ExportStackState exports the current state of the stack as JSON
func (e *StackEngine) ExportStackState() ([]byte, error) {
	summary := e.stateManager.GetStackSummary()
	return json.Marshal(summary)
}
