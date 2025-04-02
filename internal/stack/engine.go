package stack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
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
}

// NewStackEngine creates a new stack engine with the given specification
func NewStackEngine(spec StackSpec) (*StackEngine, error) {
	// Create DAG from spec
	dag, err := NewDAG(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build execution graph: %w", err)
	}

	// Create state manager
	stateManager := NewInMemoryStateManager(spec.Name)
	
	// Initialize agents in state manager
	agentIDs := make([]string, 0, len(spec.Agents))
	for _, agent := range spec.Agents {
		agentIDs = append(agentIDs, agent.ID)
	}
	stateManager.InitializeAgents(agentIDs)

	// Create engine
	ctx, cancel := context.WithCancel(context.Background())
	
	return &StackEngine{
		spec:         spec,
		stateManager: stateManager,
		dag:          dag,
		ctx:          ctx,
		cancel:       cancel,
		runID:        fmt.Sprintf("run-%d", time.Now().Unix()),
		isRunning:    false,
	}, nil
}

// BuildExecutionGraph validates and returns the DAG for the stack
func (e *StackEngine) BuildExecutionGraph() (*DAG, error) {
	return e.dag, nil
}

// Execute runs the stack in topological order
func (e *StackEngine) Execute(ctx context.Context) error {
	e.mu.Lock()
	if e.isRunning {
		e.mu.Unlock()
		return errors.New("stack is already running")
	}
	e.isRunning = true
	e.mu.Unlock()

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel
	e.ctx = ctx

	// Get execution order
	executionOrder, err := e.dag.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Create a map to track executed nodes
	executedNodes := make(map[string]bool)
	
	log.Printf("Starting stack execution: %s (Run ID: %s)", e.spec.Name, e.runID)
	log.Printf("Execution order: %v", executionOrder)

	// Execute agents in order
	for _, agentID := range executionOrder {
		select {
		case <-ctx.Done():
			log.Printf("Stack execution cancelled: %s", e.runID)
			return ctx.Err()
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
			log.Printf("Error updating agent status: %v", err)
		}

		log.Printf("Executing agent: %s", agentID)
		
		// Collect inputs from dependencies
		inputs, err := e.collectInputs(agentSpec)
		if err != nil {
			log.Printf("Error collecting inputs for agent %s: %v", agentID, err)
			e.stateManager.UpdateAgentStatus(agentID, AgentStatusFailed)
			continue
		}

		// Set agent inputs
		if err := e.stateManager.Set(agentID, "input", inputs); err != nil {
			log.Printf("Error setting agent inputs: %v", err)
		}

		// Execute agent
		outputs, err := e.executeAgent(ctx, agentSpec, inputs)
		if err != nil {
			log.Printf("Agent %s execution failed: %v", agentID, err)
			e.stateManager.UpdateAgentStatus(agentID, AgentStatusFailed)
			continue
		}

		// Set agent outputs
		if err := e.stateManager.Set(agentID, "output", outputs); err != nil {
			log.Printf("Error setting agent outputs: %v", err)
		}

		// Mark agent as completed
		e.stateManager.UpdateAgentStatus(agentID, AgentStatusCompleted)
		executedNodes[agentID] = true
		
		log.Printf("Agent completed: %s", agentID)
	}

	// Check for any agents that didn't execute
	summary := e.stateManager.GetStackSummary()
	if summary.CompletedCount != summary.TotalAgents {
		log.Printf("Stack execution completed with errors: %d/%d agents completed", 
			summary.CompletedCount, summary.TotalAgents)
		return fmt.Errorf("stack execution completed with errors: %d/%d agents completed", 
			summary.CompletedCount, summary.TotalAgents)
	}

	log.Printf("Stack execution completed successfully: %s (Run ID: %s)", e.spec.Name, e.runID)
	e.isRunning = false
	
	return nil
}

// collectInputs gathers inputs for an agent from its dependencies
func (e *StackEngine) collectInputs(agentSpec StackAgentSpec) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})
	
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
			
			output = agentState["outputs"]
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
func (e *StackEngine) executeAgent(ctx context.Context, agentSpec StackAgentSpec, inputs map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("Executing agent %s (uses: %s)", agentSpec.ID, agentSpec.Uses)
	
	// Create a real agent runtime
	runtime, err := NewRealAgentRuntime()
	if err != nil {
		return nil, fmt.Errorf("failed to create agent runtime: %w", err)
	}
	defer runtime.Cleanup()
	
	// Execute the agent using the runtime
	outputs, err := runtime.Execute(ctx, agentSpec, inputs)
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
