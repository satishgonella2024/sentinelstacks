package stack

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// StateManager manages the state of agents during stack execution
type StateManager interface {
	// Get retrieves a value from an agent's state by key
	Get(agentID string, key string) (interface{}, error)
	
	// Set stores a value in an agent's state
	Set(agentID string, key string, value interface{}) error
	
	// GetAll returns all state values for an agent
	GetAll(agentID string) (map[string]interface{}, error)
	
	// UpdateAgentStatus updates the execution status of an agent
	UpdateAgentStatus(agentID string, status AgentStatus) error
	
	// GetAgentStatus returns the current status of an agent
	GetAgentStatus(agentID string) (AgentStatus, error)
	
	// Clear removes all state for an agent
	Clear(agentID string) error
	
	// GetStackSummary returns a summary of the current stack execution
	GetStackSummary() *StackExecutionSummary
}

// InMemoryStateManager is an in-memory implementation of StateManager
type InMemoryStateManager struct {
	mu          sync.RWMutex
	agentStates map[string]*AgentState
	stackName   string
	startTime   int64
}

// NewInMemoryStateManager creates a new in-memory state manager
func NewInMemoryStateManager(stackName string) *InMemoryStateManager {
	return &InMemoryStateManager{
		agentStates: make(map[string]*AgentState),
		stackName:   stackName,
		startTime:   time.Now().Unix(),
	}
}

// InitializeAgents sets up initial state for all agents
func (m *InMemoryStateManager) InitializeAgents(agentIDs []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, id := range agentIDs {
		m.agentStates[id] = &AgentState{
			ID:      id,
			Status:  AgentStatusPending,
			Inputs:  make(map[string]interface{}),
			Outputs: make(map[string]interface{}),
		}
	}
}

// Get retrieves a value from an agent's state
func (m *InMemoryStateManager) Get(agentID, key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	state, exists := m.agentStates[agentID]
	if !exists {
		return nil, fmt.Errorf("no state found for agent %s", agentID)
	}
	
	if key == "input" {
		return state.Inputs, nil
	} else if key == "output" {
		return state.Outputs, nil
	}
	
	value, exists := state.Outputs[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in agent %s state", key, agentID)
	}
	
	return value, nil
}

// Set stores a value in an agent's state
func (m *InMemoryStateManager) Set(agentID, key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	state, exists := m.agentStates[agentID]
	if !exists {
		return fmt.Errorf("no state found for agent %s", agentID)
	}
	
	// Special cases for complete input/output objects
	if key == "input" {
		inputMap, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("input must be a map[string]interface{}")
		}
		state.Inputs = inputMap
		return nil
	} else if key == "output" {
		outputMap, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("output must be a map[string]interface{}")
		}
		state.Outputs = outputMap
		return nil
	}
	
	// Regular key-value setting
	state.Outputs[key] = value
	return nil
}

// GetAll returns all state values for an agent
func (m *InMemoryStateManager) GetAll(agentID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	state, exists := m.agentStates[agentID]
	if !exists {
		return nil, fmt.Errorf("no state found for agent %s", agentID)
	}
	
	// Create a copy to avoid data races
	allState := map[string]interface{}{
		"id":      state.ID,
		"status":  state.Status,
		"inputs":  make(map[string]interface{}),
		"outputs": make(map[string]interface{}),
	}
	
	// Copy inputs
	for k, v := range state.Inputs {
		allState["inputs"].(map[string]interface{})[k] = v
	}
	
	// Copy outputs
	for k, v := range state.Outputs {
		allState["outputs"].(map[string]interface{})[k] = v
	}
	
	return allState, nil
}

// UpdateAgentStatus updates the execution status of an agent
func (m *InMemoryStateManager) UpdateAgentStatus(agentID string, status AgentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	state, exists := m.agentStates[agentID]
	if !exists {
		return fmt.Errorf("no state found for agent %s", agentID)
	}
	
	// Update status and timestamps
	state.Status = status
	
	// Set timestamps based on status
	now := time.Now().Unix()
	if status == AgentStatusRunning && state.StartTime == 0 {
		state.StartTime = now
	} else if (status == AgentStatusCompleted || status == AgentStatusFailed) && state.EndTime == 0 {
		state.EndTime = now
	}
	
	return nil
}

// GetAgentStatus returns the current status of an agent
func (m *InMemoryStateManager) GetAgentStatus(agentID string) (AgentStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	state, exists := m.agentStates[agentID]
	if !exists {
		return "", fmt.Errorf("no state found for agent %s", agentID)
	}
	
	return state.Status, nil
}

// Clear removes all state for an agent
func (m *InMemoryStateManager) Clear(agentID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.agentStates[agentID]; !exists {
		return fmt.Errorf("no state found for agent %s", agentID)
	}
	
	// Reset state
	m.agentStates[agentID] = &AgentState{
		ID:      agentID,
		Status:  AgentStatusPending,
		Inputs:  make(map[string]interface{}),
		Outputs: make(map[string]interface{}),
	}
	
	return nil
}

// GetStackSummary returns a summary of the current stack execution
func (m *InMemoryStateManager) GetStackSummary() *StackExecutionSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	now := time.Now().Unix()
	summary := &StackExecutionSummary{
		StackName:   m.stackName,
		StartTime:   m.startTime,
		EndTime:     now,
		TotalAgents: len(m.agentStates),
		AgentStates: make(map[string]AgentState),
	}
	
	// Count agents by status
	for id, state := range m.agentStates {
		// Create a copy of agent state to avoid race conditions
		stateCopy := AgentState{
			ID:           state.ID,
			Status:       state.Status,
			ErrorMessage: state.ErrorMessage,
			StartTime:    state.StartTime,
			EndTime:      state.EndTime,
			Inputs:       make(map[string]interface{}),
			Outputs:      make(map[string]interface{}),
		}
		
		// Copy inputs and outputs
		for k, v := range state.Inputs {
			stateCopy.Inputs[k] = v
		}
		for k, v := range state.Outputs {
			stateCopy.Outputs[k] = v
		}
		
		summary.AgentStates[id] = stateCopy
		
		// Update counts
		switch state.Status {
		case AgentStatusCompleted:
			summary.CompletedCount++
		case AgentStatusFailed:
			summary.FailedCount++
		case AgentStatusBlocked:
			summary.BlockedCount++
		}
	}
	
	return summary
}

// JSONStateManager is a StateManager that persists state to JSON files
type JSONStateManager struct {
	InMemoryStateManager // Embed in-memory manager for core functionality
	stateFilePath        string
}

// SerializeState converts the current state to JSON
func (m *InMemoryStateManager) SerializeState() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create a serializable structure
	state := struct {
		StackName   string                  `json:"stackName"`
		StartTime   int64                   `json:"startTime"`
		AgentStates map[string]*AgentState  `json:"agentStates"`
	}{
		StackName:   m.stackName,
		StartTime:   m.startTime,
		AgentStates: m.agentStates,
	}
	
	return json.Marshal(state)
}

// DeserializeState loads state from JSON
func (m *InMemoryStateManager) DeserializeState(data []byte) error {
	var state struct {
		StackName   string                  `json:"stackName"`
		StartTime   int64                   `json:"startTime"`
		AgentStates map[string]*AgentState  `json:"agentStates"`
	}
	
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.stackName = state.StackName
	m.startTime = state.StartTime
	m.agentStates = state.AgentStates
	
	return nil
}
