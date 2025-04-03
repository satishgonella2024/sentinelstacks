package memory

import (
	"fmt"
	"sync"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// DefaultFactory is the default implementation of MemoryStoreFactory
type DefaultFactory struct {
	mu               sync.Mutex
	memoryStores     map[string]types.MemoryStore
	vectorStores     map[string]types.VectorStore
	overrideCreators map[types.MemoryStoreType]func(config types.MemoryConfig) (types.MemoryStore, error)
}

// NewDefaultFactory creates a new default factory
func NewDefaultFactory() *DefaultFactory {
	return &DefaultFactory{
		memoryStores:     make(map[string]types.MemoryStore),
		vectorStores:     make(map[string]types.VectorStore),
		overrideCreators: make(map[types.MemoryStoreType]func(config types.MemoryConfig) (types.MemoryStore, error)),
	}
}

// Create creates a new memory store of the requested type
func (f *DefaultFactory) Create(storeType types.MemoryStoreType, config types.MemoryConfig) (types.MemoryStore, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Generate store key
	storeKey := string(storeType)
	if config.CollectionName != "" {
		storeKey += ":" + config.CollectionName
	}
	if config.Namespace != "" {
		storeKey += ":" + config.Namespace
	}

	// Check if store already exists
	if store, ok := f.memoryStores[storeKey]; ok {
		return store, nil
	}

	// Check if there's an override creator for this type
	if creator, ok := f.overrideCreators[storeType]; ok {
		store, err := creator(config)
		if err != nil {
			return nil, err
		}
		f.memoryStores[storeKey] = store
		return store, nil
	}

	// Create new store based on type
	var store types.MemoryStore
	var err error

	switch storeType {
	case types.MemoryStoreTypeLocal:
		store, err = NewLocalMemoryStore(config)
	case types.MemoryStoreTypeSQLite:
		store, err = NewSQLiteMemoryStore(config)
	case types.MemoryStoreTypeChroma:
		store, err = NewChromaMemoryStore(config)
	default:
		return nil, fmt.Errorf("unsupported memory store type: %s", storeType)
	}

	if err != nil {
		return nil, err
	}

	// Save for future use
	f.memoryStores[storeKey] = store

	return store, nil
}

// CreateVector creates a new vector store
func (f *DefaultFactory) CreateVector(config types.MemoryConfig) (types.VectorStore, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Generate store key
	storeKey := "vector"
	if config.CollectionName != "" {
		storeKey += ":" + config.CollectionName
	}
	if config.Namespace != "" {
		storeKey += ":" + config.Namespace
	}

	// Check if store already exists
	if store, ok := f.vectorStores[storeKey]; ok {
		return store, nil
	}

	// Create new store based on type
	var store types.VectorStore
	var err error

	// Determine which vector store to create based on configuration
	if config.StoragePath != "" && config.StoragePath != "memory" {
		// If storage path is provided, use SQLite vector store
		store, err = NewSQLiteVectorStore(config)
	} else if config.ConnectionString != "" {
		// If connection string is provided, use Chroma vector store
		store, err = NewChromaVectorStore(config)
	} else {
		// Default to local in-memory vector store
		store, err = NewLocalVectorStore(config)
	}

	if err != nil {
		return nil, err
	}

	// Save for future use
	f.vectorStores[storeKey] = store

	return store, nil
}

// RegisterOverride registers a custom creator for a memory store type
func (f *DefaultFactory) RegisterOverride(storeType types.MemoryStoreType, creator func(config types.MemoryConfig) (types.MemoryStore, error)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.overrideCreators[storeType] = creator
}

// NewMemoryStoreFactory creates a new memory store factory
func NewMemoryStoreFactory(storageBasePath string) (types.MemoryStoreFactory, error) {
	return NewDefaultFactory(), nil
}
