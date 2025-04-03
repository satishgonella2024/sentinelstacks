package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// PersistentStateManager implements types.StateManager using memory stores
type PersistentStateManager struct {
	memoryManager *MemoryManager
	ctx           context.Context
	stackName     string
	agentIDs      []string
	agentStates   map[string]*types.AgentState
	summary       *types.StackExecutionSummary
	mu            sync.RWMutex
}

// NewPersistentStateManager creates a new persistent state manager
func NewPersistentStateManager(ctx context.Context, stackName string, factory types.MemoryStoreFactory, stackID, executionID string) (*PersistentStateManager, error) {
	// Create memory manager
	memManager, err := NewMemoryManager(ctx, factory, stackID, executionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory manager: %w", err)
	}

	// Create state manager
	manager := &PersistentStateManager{
		memoryManager: memManager,
		ctx:           ctx,
		stackName:     stackName,
		agentIDs:      []string{},
		agentStates:   make(map[string]*types.AgentState),
		summary: &types.StackExecutionSummary{
			StackName:   stackName,
			ExecutionID: executionID,
			StartTime:   time.Now(),
			AgentStates: make(map[string]*types.AgentState),
		},
	}

	// Try to load existing state
	err = manager.loadState()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	return manager, nil
}

// InitializeAgents initializes the state for a list of agents
func (m *PersistentStateManager) InitializeAgents(agentIDs []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.agentIDs = agentIDs
	m.summary.TotalAgents = len(agentIDs)

	// Initialize agent states
	for _, agentID := range agentIDs {
		if _, exists := m.agentStates[agentID]; !exists {
			m.agentStates[agentID] = &types.AgentState{
				ID:           agentID,
				Status:       types.AgentStatusPending,
				StartTime:    time.Time{},
				EndTime:      time.Time{},
				Dependencies: []string{},
				Inputs:       make(map[string]interface{}),
				Outputs:      make(map[string]interface{}),
			}
		}
		
		m.summary.AgentStates[agentID] = m.agentStates[agentID]
	}

	// Save state
	m.saveState()
}

// UpdateAgentStatus updates the status of an agent
func (m *PersistentStateManager) UpdateAgentStatus(agentID string, status types.AgentStatus) error {
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
	if status == types.AgentStatusRunning && agentState.StartTime.IsZero() {
		agentState.StartTime = now
	} else if (status == types.AgentStatusCompleted || status == types.AgentStatusFailed) && agentState.EndTime.IsZero() {
		agentState.EndTime = now
	}
	
	// Update summary counts
	if prevStatus == types.AgentStatusPending && status == types.AgentStatusRunning {
		// No change in counts
	} else if prevStatus == types.AgentStatusRunning && status == types.AgentStatusCompleted {
		m.summary.CompletedCount++
	} else if prevStatus == types.AgentStatusRunning && status == types.AgentStatusFailed {
		m.summary.FailedCount++
	} else if prevStatus == types.AgentStatusPending && status == types.AgentStatusBlocked {
		m.summary.BlockedCount++
	} else if prevStatus == types.AgentStatusBlocked && status == types.AgentStatusRunning {
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
func (m *PersistentStateManager) GetStackSummary() *types.StackExecutionSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create a copy to avoid concurrent access issues
	summary := &types.StackExecutionSummary{
		StackName:      m.summary.StackName,
		ExecutionID:    m.summary.ExecutionID,
		StartTime:      m.summary.StartTime,
		EndTime:        m.summary.EndTime,
		TotalAgents:    m.summary.TotalAgents,
		CompletedCount: m.summary.CompletedCount,
		FailedCount:    m.summary.FailedCount,
		BlockedCount:   m.summary.BlockedCount,
		AgentStates:    make(map[string]*types.AgentState),
	}
	
	// Copy agent states
	for agentID, state := range m.summary.AgentStates {
		stateCopy := *state
		summary.AgentStates[agentID] = &stateCopy
	}
	
	return summary
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
	summary, ok := summaryValue.(*types.StackExecutionSummary)
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
