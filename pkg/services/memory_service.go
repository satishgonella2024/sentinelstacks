package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
)

// MemoryService provides memory storage operations for agents and stacks
type MemoryService struct {
	factory       memory.MemoryStoreFactory
	stores        map[string]memory.MemoryStore
	vectorStores  map[string]memory.VectorStore
	mu            sync.Mutex
	defaultConfig memory.MemoryConfig
}

// NewMemoryService creates a new memory service
func NewMemoryService(factory memory.MemoryStoreFactory) *MemoryService {
	return &MemoryService{
		factory:      factory,
		stores:       make(map[string]memory.MemoryStore),
		vectorStores: make(map[string]memory.VectorStore),
		defaultConfig: memory.MemoryConfig{
			TTL:             24 * time.Hour,
			VectorDimensions: 1536,
		},
	}
}

// GetStore gets or creates a memory store for the given agent
func (s *MemoryService) GetStore(ctx context.Context, agentID string, storeType memory.MemoryStoreType) (memory.MemoryStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create store key
	storeKey := fmt.Sprintf("%s-%s", agentID, storeType)
	
	// Check if store already exists
	if store, ok := s.stores[storeKey]; ok {
		return store, nil
	}
	
	// Create config for this store
	config := s.defaultConfig
	config.CollectionName = agentID
	config.Namespace = agentID
	
	// Create new store
	store, err := s.factory.Create(storeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory store: %w", err)
	}
	
	// Save for future use
	s.stores[storeKey] = store
	
	return store, nil
}

// GetVectorStore gets or creates a vector store for the given agent
func (s *MemoryService) GetVectorStore(ctx context.Context, agentID string, storeType memory.MemoryStoreType) (memory.VectorStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create store key
	storeKey := fmt.Sprintf("%s-%s-vector", agentID, storeType)
	
	// Check if store already exists
	if store, ok := s.vectorStores[storeKey]; ok {
		return store, nil
	}
	
	// Create config for this store
	config := s.defaultConfig
	config.CollectionName = agentID
	config.Namespace = agentID
	
	// Create new vector store
	store, err := s.factory.CreateVector(storeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}
	
	// Save for future use
	s.vectorStores[storeKey] = store
	
	return store, nil
}

// GetStackStore gets or creates a memory store for the given stack
func (s *MemoryService) GetStackStore(ctx context.Context, stackID string, storeType memory.MemoryStoreType) (memory.MemoryStore, error) {
	// Stack stores use a different namespace prefix
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create store key
	storeKey := fmt.Sprintf("stack-%s-%s", stackID, storeType)
	
	// Check if store already exists
	if store, ok := s.stores[storeKey]; ok {
		return store, nil
	}
	
	// Create config for this store
	config := s.defaultConfig
	config.CollectionName = fmt.Sprintf("stack_%s", stackID)
	config.Namespace = fmt.Sprintf("stack_%s", stackID)
	
	// Create new store
	store, err := s.factory.Create(storeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack memory store: %w", err)
	}
	
	// Save for future use
	s.stores[storeKey] = store
	
	return store, nil
}

// SaveAgentState saves an agent's state
func (s *MemoryService) SaveAgentState(ctx context.Context, agentID string, state map[string]interface{}) error {
	// Get store
	store, err := s.GetStore(ctx, agentID, memory.MemoryStoreTypeLocal)
	if err != nil {
		return fmt.Errorf("failed to get store: %w", err)
	}
	
	// Save state
	err = store.Save(ctx, "state", state)
	if err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}
	
	return nil
}

// LoadAgentState loads an agent's state
func (s *MemoryService) LoadAgentState(ctx context.Context, agentID string) (map[string]interface{}, error) {
	// Get store
	store, err := s.GetStore(ctx, agentID, memory.MemoryStoreTypeLocal)
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	
	// Load state
	stateValue, err := store.Load(ctx, "state")
	if err != nil {
		// Return empty state if not found
		if err.Error() == "key not found: state" {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to load state: %w", err)
	}
	
	// Convert to map
	state, ok := stateValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid state format")
	}
	
	return state, nil
}

// SaveEmbedding saves an embedding
func (s *MemoryService) SaveEmbedding(ctx context.Context, agentID, documentID string, embedding []float32, metadata map[string]interface{}) error {
	// Get vector store
	store, err := s.GetVectorStore(ctx, agentID, memory.MemoryStoreTypeChroma)
	if err != nil {
		return fmt.Errorf("failed to get vector store: %w", err)
	}
	
	// Save embedding
	err = store.SaveEmbedding(ctx, documentID, embedding, metadata)
	if err != nil {
		return fmt.Errorf("failed to save embedding: %w", err)
	}
	
	return nil
}

// SearchSimilar searches for similar documents
func (s *MemoryService) SearchSimilar(ctx context.Context, agentID string, queryEmbedding []float32, topK int) ([]memory.SimilarityMatch, error) {
	// Get vector store
	store, err := s.GetVectorStore(ctx, agentID, memory.MemoryStoreTypeChroma)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector store: %w", err)
	}
	
	// Query
	matches, err := store.Query(ctx, queryEmbedding, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to query vector store: %w", err)
	}
	
	return matches, nil
}

// SaveStackData saves data for a stack
func (s *MemoryService) SaveStackData(ctx context.Context, stackID, key string, value interface{}) error {
	// Get stack store
	store, err := s.GetStackStore(ctx, stackID, memory.MemoryStoreTypeLocal)
	if err != nil {
		return fmt.Errorf("failed to get stack store: %w", err)
	}
	
	// Save data
	err = store.Save(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to save stack data: %w", err)
	}
	
	return nil
}

// LoadStackData loads data for a stack
func (s *MemoryService) LoadStackData(ctx context.Context, stackID, key string) (interface{}, error) {
	// Get stack store
	store, err := s.GetStackStore(ctx, stackID, memory.MemoryStoreTypeLocal)
	if err != nil {
		return nil, fmt.Errorf("failed to get stack store: %w", err)
	}
	
	// Load data
	value, err := store.Load(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load stack data: %w", err)
	}
	
	return value, nil
}

// Close closes all stores
func (s *MemoryService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var lastErr error
	
	// Close all stores
	for key, store := range s.stores {
		if err := store.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close store %s: %w", key, err)
		}
		delete(s.stores, key)
	}
	
	// Close all vector stores
	for key, store := range s.vectorStores {
		if err := store.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close vector store %s: %w", key, err)
		}
		delete(s.vectorStores, key)
	}
	
	return lastErr
}
