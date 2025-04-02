package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

// PersistentStateManager implements StateManager using the memory system
type PersistentStateManager struct {
	memoryManager *MemoryManager
	ctx           context.Context
	stackName     string
	agentIDs      []string
	agentStates   map[string]*stack.AgentState
	summary       *stack.StackExecutionSummary
	mu            sync.RWMutex
}

// NewPersistentStateManager creates a new persistent state manager
func NewPersistentStateManager(ctx context.Context, stackName string, factory memory.MemoryStoreFactory, stackID, executionID string) (*PersistentStateManager, error) {
	// Create memory manager
	memoryManager := NewMemoryManager(factory, stackID, executionID)
	
	// Create summary
	summary := &stack.StackExecutionSummary{
		StackName:     stackName,
		ExecutionID:   executionID,
		StartTime:     time.Now(),
		TotalAgents:   0,
		CompletedCount: 0,
		FailedCount:   0,
		BlockedCount:  0,
		AgentStates:   make(map[string]*stack.AgentState),
	}
	
	return &PersistentStateManager{
		memoryManager: memoryManager,
		ctx:           ctx,
		stackName:     stackName,
		agentIDs:      []string{},
		agentStates:   make(map[string]*stack.AgentState),
		summary:       summary,
	}, nil
}

// InitializeAgents initializes the state for a list of agents
func (m *PersistentStateManager) InitializeAgents(agentIDs []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.agentIDs = agentIDs
	m.summary.TotalAgents = len(agentIDs)
	
	// Initialize agent states
	for _, agentID := range agentIDs {
		m.agentStates[agentID] = &stack.AgentState{
			ID:           agentID,
			Status:       stack.AgentStatusPending,
			StartTime:    time.Time{},
			EndTime:      time.Time{},
			Dependencies: []string{},
		}
		
		m.summary.AgentStates[agentID] = m.agentStates[agentID]
	}
	
	// Save initial state
	m.saveState()
}

// UpdateAgentStatus updates the status of an agent
func (m *PersistentStateManager) UpdateAgentStatus(agentID string, status stack.AgentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if agent exists
	agentState, ok := m.agentStates[agentID]
	if !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Update status
	prevStatus := agentState.Status
	agentState.Status = status
	
	// Update timestamps
	now := time.Now()
	if status == stack.AgentStatusRunning && agentState.StartTime.IsZero() {
		agentState.StartTime = now
	} else if (status == stack.AgentStatusCompleted || status == stack.AgentStatusFailed) && agentState.EndTime.IsZero() {
		agentState.EndTime = now
	}
	
	// Update summary counts
	if prevStatus == stack.AgentStatusPending && status == stack.AgentStatusRunning {
		// No change in counts
	} else if prevStatus == stack.AgentStatusRunning && status == stack.AgentStatusCompleted {
		m.summary.CompletedCount++
	} else if prevStatus == stack.AgentStatusRunning && status == stack.AgentStatusFailed {
		m.summary.FailedCount++
	} else if prevStatus == stack.AgentStatusPending && status == stack.AgentStatusBlocked {
		m.summary.BlockedCount++
	} else if prevStatus == stack.AgentStatusBlocked && status == stack.AgentStatusRunning {
		m.summary.BlockedCount--
	}
	
	// Save updated state
	m.saveState()
	
	return nil
}

// UpdateAgentErrorMessage updates the error message for an agent
func (m *PersistentStateManager) UpdateAgentErrorMessage(agentID string, errorMessage string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if agent exists
	agentState, ok := m.agentStates[agentID]
	if !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Update error message
	agentState.ErrorMessage = errorMessage
	
	// Save updated state
	m.saveState()
	
	return nil
}

// UpdateAgentDependencies updates the dependencies for an agent
func (m *PersistentStateManager) UpdateAgentDependencies(agentID string, dependencies []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if agent exists
	agentState, ok := m.agentStates[agentID]
	if !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Update dependencies
	agentState.Dependencies = dependencies
	
	// Save updated state
	m.saveState()
	
	return nil
}

// Get gets a value from an agent's state
func (m *PersistentStateManager) Get(agentID, key string) (interface{}, error) {
	// Check if agent exists
	m.mu.RLock()
	_, ok := m.agentStates[agentID]
	m.mu.RUnlock()
	
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Get all values
	allValues, err := m.GetAll(agentID)
	if err != nil {
		return nil, err
	}
	
	// Get specific value
	value, ok := allValues[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	return value, nil
}

// Set sets a value in an agent's state
func (m *PersistentStateManager) Set(agentID, key string, value interface{}) error {
	// Check if agent exists
	m.mu.RLock()
	_, ok := m.agentStates[agentID]
	m.mu.RUnlock()
	
	if !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Get all values
	allValues, err := m.GetAll(agentID)
	if err != nil {
		// If agent has no state yet, create empty map
		if err.Error() == fmt.Sprintf("failed to load agent state: key not found: %s", "state") {
			allValues = make(map[string]interface{})
		} else {
			return err
		}
	}
	
	// Set value
	allValues[key] = value
	
	// Save state
	err = m.memoryManager.SaveAgentState(m.ctx, agentID, allValues)
	if err != nil {
		return fmt.Errorf("failed to save agent state: %w", err)
	}
	
	// If key is "input", also save as agent input
	if key == "input" {
		err = m.memoryManager.SaveAgentInput(m.ctx, agentID, value)
		if err != nil {
			return fmt.Errorf("failed to save agent input: %w", err)
		}
	}
	
	// If key is "output", also save as agent output
	if key == "output" {
		err = m.memoryManager.SaveAgentOutput(m.ctx, agentID, value)
		if err != nil {
			return fmt.Errorf("failed to save agent output: %w", err)
		}
	}
	
	return nil
}

// GetAll gets all values from an agent's state
func (m *PersistentStateManager) GetAll(agentID string) (map[string]interface{}, error) {
	// Check if agent exists
	m.mu.RLock()
	_, ok := m.agentStates[agentID]
	m.mu.RUnlock()
	
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}
	
	// Get agent state
	state, err := m.memoryManager.LoadAgentState(m.ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent state: %w", err)
	}
	
	return state, nil
}

// GetStackSummary gets the current execution summary
func (m *PersistentStateManager) GetStackSummary() *stack.StackExecutionSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create a copy to avoid concurrent access issues
	summary := &stack.StackExecutionSummary{
		StackName:      m.summary.StackName,
		ExecutionID:    m.summary.ExecutionID,
		StartTime:      m.summary.StartTime,
		EndTime:        m.summary.EndTime,
		TotalAgents:    m.summary.TotalAgents,
		CompletedCount: m.summary.CompletedCount,
		FailedCount:    m.summary.FailedCount,
		BlockedCount:   m.summary.BlockedCount,
		AgentStates:    make(map[string]*stack.AgentState),
	}
	
	// Copy agent states
	for agentID, state := range m.summary.AgentStates {
		stateCopy := *state
		summary.AgentStates[agentID] = &stateCopy
	}
	
	return summary
}

// Close releases all resources
func (m *PersistentStateManager) Close() error {
	// Set end time if not already set
	if m.summary.EndTime.IsZero() {
		m.summary.EndTime = time.Now()
	}
	
	// Save final state
	m.saveState()
	
	// Close memory manager
	return m.memoryManager.Close()
}

// saveState saves the current execution summary to memory
func (m *PersistentStateManager) saveState() {
	// Get stack store
	stackStore, err := m.memoryManager.GetStackStore(m.ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to get stack store: %v\n", err)
		return
	}
	
	// Save summary
	err = stackStore.Save(m.ctx, "summary", m.summary)
	if err != nil {
		fmt.Printf("Warning: Failed to save stack summary: %v\n", err)
	}
}

// loadState loads the execution summary from memory
func (m *PersistentStateManager) loadState() error {
	// Get stack store
	stackStore, err := m.memoryManager.GetStackStore(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to get stack store: %w", err)
	}
	
	// Load summary
	summaryValue, err := stackStore.Load(m.ctx, "summary")
	if err != nil {
		// If summary not found, return nil
		return nil
	}
	
	// Convert to summary
	summary, ok := summaryValue.(*stack.StackExecutionSummary)
	if !ok {
		return fmt.Errorf("invalid summary format: %T", summaryValue)
	}
	
	// Update summary
	m.summary = summary
	
	// Update agent states
	for agentID, state := range summary.AgentStates {
		m.agentStates[agentID] = state
	}
	
	return nil
}
