// Package api provides a unified API for the Sentinel Stacks system
package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/satishgonella2024/sentinelstacks/pkg/memory"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// MemoryServiceConfig contains configuration for the memory service
type MemoryServiceConfig struct {
	// StoragePath is where memory data is persisted
	StoragePath string

	// EmbeddingProvider specifies the embedding provider to use
	EmbeddingProvider string

	// EmbeddingModel specifies the embedding model to use
	EmbeddingModel string

	// EmbeddingDimensions specifies the dimensions of embeddings
	EmbeddingDimensions int
}

// MemoryService implements types.MemoryService
type MemoryService struct {
	config       MemoryServiceConfig
	memoryStores map[string]types.MemoryStore
	vectorStores map[string]types.VectorStore
	factory      types.MemoryStoreFactory
	mu           sync.RWMutex
}

// NewMemoryService creates a new memory service
func NewMemoryService(config MemoryServiceConfig) (*MemoryService, error) {
	// Create a memory factory
	factoryConfig := memory.FactoryConfig{
		BasePath:                config.StoragePath,
		PreferredStoreType:      config.EmbeddingProvider,
		DefaultVectorDimensions: config.EmbeddingDimensions,
	}

	factory := memory.NewFactory(factoryConfig)

	return &MemoryService{
		config:       config,
		memoryStores: make(map[string]types.MemoryStore),
		vectorStores: make(map[string]types.VectorStore),
		factory:      factory,
	}, nil
}

// getOrCreateMemoryStore gets or creates a memory store for a collection
func (s *MemoryService) getOrCreateMemoryStore(ctx context.Context, collection string) (types.MemoryStore, error) {
	s.mu.RLock()
	store, exists := s.memoryStores[collection]
	s.mu.RUnlock()

	if exists {
		return store, nil
	}

	// Use the factory to create a memory store
	store, err := s.factory.CreateMemoryStore(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory store: %w", err)
	}

	// Cache the store
	s.mu.Lock()
	s.memoryStores[collection] = store
	s.mu.Unlock()

	return store, nil
}

// getOrCreateVectorStore gets or creates a vector store for a collection
func (s *MemoryService) getOrCreateVectorStore(ctx context.Context, collection string) (types.VectorStore, error) {
	s.mu.RLock()
	store, exists := s.vectorStores[collection]
	s.mu.RUnlock()

	if exists {
		return store, nil
	}

	// Use the factory to create a vector store
	store, err := s.factory.CreateVectorStore(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	// Cache the store
	s.mu.Lock()
	s.vectorStores[collection] = store
	s.mu.Unlock()

	return store, nil
}

// StoreValue stores a value in memory
func (s *MemoryService) StoreValue(ctx context.Context, collection string, key string, value interface{}) error {
	store, err := s.getOrCreateMemoryStore(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get memory store: %w", err)
	}

	return store.Save(ctx, key, value)
}

// RetrieveValue retrieves a value from memory
func (s *MemoryService) RetrieveValue(ctx context.Context, collection string, key string) (interface{}, error) {
	store, err := s.getOrCreateMemoryStore(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory store: %w", err)
	}

	return store.Load(ctx, key)
}

// StoreEmbedding stores text with vector embedding
func (s *MemoryService) StoreEmbedding(ctx context.Context, collection string, key string, text string, metadata map[string]interface{}) error {
	store, err := s.getOrCreateVectorStore(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get vector store: %w", err)
	}

	return store.StoreVector(ctx, key, text, metadata)
}

// SearchSimilar finds similar texts using vector similarity
func (s *MemoryService) SearchSimilar(ctx context.Context, collection string, text string, limit int) ([]types.MemoryMatch, error) {
	store, err := s.getOrCreateVectorStore(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector store: %w", err)
	}

	// Empty filter for now
	filter := make(map[string]interface{})

	return store.SearchVector(ctx, text, limit, filter)
}
