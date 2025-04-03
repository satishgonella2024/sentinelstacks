// Package memory provides memory store implementations
package memory

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// Factory implements types.MemoryStoreFactory
type Factory struct {
	config       FactoryConfig
	memoryStores map[string]types.MemoryStore
	vectorStores map[string]types.VectorStore
	mu           sync.RWMutex
}

// FactoryConfig contains configuration for the memory factory
type FactoryConfig struct {
	// BasePath is the root directory for persistent stores
	BasePath string

	// DefaultVectorDimensions is the default dimensionality for vector stores
	DefaultVectorDimensions int

	// PreferredStoreType indicates the preferred memory store implementation
	PreferredStoreType string // "sqlite", "chroma", "local"
}

// NewFactory creates a new memory store factory
func NewFactory(config FactoryConfig) *Factory {
	return &Factory{
		config:       config,
		memoryStores: make(map[string]types.MemoryStore),
		vectorStores: make(map[string]types.VectorStore),
	}
}

// CreateMemoryStore creates a basic key-value memory store
func (f *Factory) CreateMemoryStore(ctx context.Context, name string) (types.MemoryStore, error) {
	f.mu.RLock()
	store, exists := f.memoryStores[name]
	f.mu.RUnlock()

	if exists {
		return store, nil
	}

	var newStore types.MemoryStore
	var err error

	switch f.config.PreferredStoreType {
	case "sqlite":
		storePath := filepath.Join(f.config.BasePath, "memory", fmt.Sprintf("%s.db", name))
		newStore, err = NewSQLiteStore(storePath)
	case "local", "":
		newStore = NewLocalStore()
		err = nil
	default:
		return nil, fmt.Errorf("unsupported store type: %s", f.config.PreferredStoreType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create memory store: %w", err)
	}

	f.mu.Lock()
	f.memoryStores[name] = newStore
	f.mu.Unlock()

	return newStore, nil
}

// CreateVectorStore creates a vector-based memory store
func (f *Factory) CreateVectorStore(ctx context.Context, name string) (types.VectorStore, error) {
	f.mu.RLock()
	store, exists := f.vectorStores[name]
	f.mu.RUnlock()

	if exists {
		return store, nil
	}

	var newStore types.VectorStore
	var err error

	dimensions := f.config.DefaultVectorDimensions
	if dimensions <= 0 {
		dimensions = 1536 // Default for OpenAI embeddings
	}

	switch f.config.PreferredStoreType {
	case "chroma":
		chromaConfig := ChromaConfig{
			CollectionName: name,
			Dimensions:     dimensions,
		}
		newStore, err = NewChromaStore(chromaConfig)
	case "sqlite":
		storePath := filepath.Join(f.config.BasePath, "vectors", fmt.Sprintf("%s.db", name))
		newStore, err = NewSQLiteVectorStore(storePath, dimensions)
	case "local", "":
		newStore = NewLocalVectorStore(dimensions)
		err = nil
	default:
		return nil, fmt.Errorf("unsupported vector store type: %s", f.config.PreferredStoreType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	f.mu.Lock()
	f.vectorStores[name] = newStore
	f.mu.Unlock()

	return newStore, nil
}
