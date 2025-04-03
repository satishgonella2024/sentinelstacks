package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// MemoryManager manages memory for a stack execution
type MemoryManager struct {
	factory       types.MemoryStoreFactory
	stores        map[string]types.MemoryStore
	stackID       string
	executionID   string
	defaultConfig types.MemoryConfig
	mu            sync.Mutex
}

// NewMemoryManager creates a new memory manager for a stack execution
func NewMemoryManager(ctx context.Context, factory types.MemoryStoreFactory, stackID, executionID string) (*MemoryManager, error) {
	return &MemoryManager{
		factory:     factory,
		stores:      make(map[string]types.MemoryStore),
		stackID:     stackID,
		executionID: executionID,
		defaultConfig: types.MemoryConfig{
			TTL:              24 * time.Hour, // Default TTL for execution data
			VectorDimensions: 1536,
		},
	}, nil
}

// GetAgentStore gets or creates a memory store for an agent in the stack
func (m *MemoryManager) GetAgentStore(ctx context.Context, agentID string) (types.MemoryStore, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create store key
	storeKey := fmt.Sprintf("%s-%s", agentID, m.executionID)
	
	// Check if store already exists
	if store, ok := m.stores[storeKey]; ok {
		return store, nil
	}
	
	// Create config for this store
	config := m.defaultConfig
	config.CollectionName = fmt.Sprintf("stack_%s_agent_%s", m.stackID, agentID)
	config.Namespace = m.executionID
	
	// Create new store
	store, err := m.factory.Create(types.MemoryStoreTypeLocal, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent memory store: %w", err)
	}
	
	// Save for future use
	m.stores[storeKey] = store
	
	return store, nil
}

// GetStackStore gets or creates a memory store for the stack itself
func (m *MemoryManager) GetStackStore(ctx context.Context) (types.MemoryStore, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create store key
	storeKey := fmt.Sprintf("stack-%s", m.executionID)
	
	// Check if store already exists
	if store, ok := m.stores[storeKey]; ok {
		return store, nil
	}
	
	// Create config for this store
	config := m.defaultConfig
	config.CollectionName = fmt.Sprintf("stack_%s", m.stackID)
	config.Namespace = m.executionID
	
	// Create new store
	store, err := m.factory.Create(types.MemoryStoreTypeLocal, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack memory store: %w", err)
	}
	
	// Save for future use
	m.stores[storeKey] = store
	
	return store, nil
}

// SaveAgentState saves an agent's state
func (m *MemoryManager) SaveAgentState(ctx context.Context, agentID string, state map[string]interface{}) error {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Save state
	err = store.Save(ctx, "state", state)
	if err != nil {
		return fmt.Errorf("failed to save agent state: %w", err)
	}
	
	return nil
}

// LoadAgentState loads an agent's state
func (m *MemoryManager) LoadAgentState(ctx context.Context, agentID string) (map[string]interface{}, error) {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Load state
	stateValue, err := store.Load(ctx, "state")
	if err != nil {
		// Return empty state if not found
		if fmt.Sprintf("%v", err) == fmt.Sprintf("key not found: %s", "state") {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to load agent state: %w", err)
	}
	
	// Convert to map
	state, ok := stateValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid state format: %T", stateValue)
	}
	
	return state, nil
}

// SaveAgentInput saves an agent's input
func (m *MemoryManager) SaveAgentInput(ctx context.Context, agentID string, input interface{}) error {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Save input
	err = store.Save(ctx, "input", input)
	if err != nil {
		return fmt.Errorf("failed to save agent input: %w", err)
	}
	
	return nil
}

// LoadAgentInput loads an agent's input
func (m *MemoryManager) LoadAgentInput(ctx context.Context, agentID string) (interface{}, error) {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Load input
	input, err := store.Load(ctx, "input")
	if err != nil {
		return nil, fmt.Errorf("failed to load agent input: %w", err)
	}
	
	return input, nil
}

// SaveAgentOutput saves an agent's output
func (m *MemoryManager) SaveAgentOutput(ctx context.Context, agentID string, output interface{}) error {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Save output
	err = store.Save(ctx, "output", output)
	if err != nil {
		return fmt.Errorf("failed to save agent output: %w", err)
	}
	
	// Also save to stack store for access by other agents
	stackStore, err := m.GetStackStore(ctx)
	if err != nil {
		return fmt.Errorf("failed to get stack store: %w", err)
	}
	
	// Save output to stack store with agent ID as key
	err = stackStore.Save(ctx, fmt.Sprintf("agent_%s_output", agentID), output)
	if err != nil {
		return fmt.Errorf("failed to save agent output to stack store: %w", err)
	}
	
	return nil
}

// LoadAgentOutput loads an agent's output
func (m *MemoryManager) LoadAgentOutput(ctx context.Context, agentID string) (interface{}, error) {
	// Get agent store
	store, err := m.GetAgentStore(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent store: %w", err)
	}
	
	// Load output
	output, err := store.Load(ctx, "output")
	if err != nil {
		return nil, fmt.Errorf("failed to load agent output: %w", err)
	}
	
	return output, nil
}

// Close releases all resources
func (m *MemoryManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var lastErr error
	
	// Close all stores
	for key, store := range m.stores {
		if err := store.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close store %s: %w", key, err)
		}
		delete(m.stores, key)
	}
	
	return lastErr
}
