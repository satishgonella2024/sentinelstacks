package memory

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultMemoryStoreFactory is the default implementation of MemoryStoreFactory
type DefaultMemoryStoreFactory struct {
	// basePath is the base path for file-backed stores
	basePath string
}

// NewMemoryStoreFactory creates a new memory store factory
func NewMemoryStoreFactory(basePath string) (*DefaultMemoryStoreFactory, error) {
	// If base path is not provided, use a default
	if basePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(home, ".sentinel", "memory")
	}

	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	return &DefaultMemoryStoreFactory{
		basePath: basePath,
	}, nil
}

// Create creates a memory store of the specified type
func (f *DefaultMemoryStoreFactory) Create(storeType MemoryStoreType, config MemoryConfig) (MemoryStore, error) {
	// Prepare config
	if config.ConnectionString == "" && storeType != MemoryStoreTypeLocal {
		// Set default connection string based on type
		switch storeType {
		case MemoryStoreTypeSQLite:
			// Use a file in the base path
			dbName := "sentinel_memory.db"
			if config.CollectionName != "" {
				dbName = fmt.Sprintf("%s.db", config.CollectionName)
			}
			config.ConnectionString = filepath.Join(f.basePath, dbName)
		case MemoryStoreTypeChroma:
			config.ConnectionString = "http://localhost:8000"
		case MemoryStoreTypeRedis:
			config.ConnectionString = "localhost:6379"
		case MemoryStoreTypePostgres:
			config.ConnectionString = "postgres://localhost/sentinel?sslmode=disable"
		}
	}

	// Create store based on type
	switch storeType {
	case MemoryStoreTypeLocal:
		return NewLocalMemoryStore(config)
	case MemoryStoreTypeSQLite:
		return NewSQLiteMemoryStore(config)
	case MemoryStoreTypeChroma:
		return NewChromaMemoryStore(config)
	default:
		return nil, fmt.Errorf("unsupported memory store type: %s", storeType)
	}
}

// CreateVector creates a vector store of the specified type
func (f *DefaultMemoryStoreFactory) CreateVector(storeType MemoryStoreType, config MemoryConfig) (VectorStore, error) {
	// Prepare config
	if config.ConnectionString == "" && storeType != MemoryStoreTypeLocal {
		// Set default connection string based on type
		switch storeType {
		case MemoryStoreTypeSQLite:
			// Use a file in the base path
			dbName := "sentinel_vectors.db"
			if config.CollectionName != "" {
				dbName = fmt.Sprintf("%s_vectors.db", config.CollectionName)
			}
			config.ConnectionString = filepath.Join(f.basePath, dbName)
		case MemoryStoreTypeChroma:
			config.ConnectionString = "http://localhost:8000"
		case MemoryStoreTypePostgres:
			config.ConnectionString = "postgres://localhost/sentinel?sslmode=disable"
		}
	}

	// Create vector store based on type
	switch storeType {
	case MemoryStoreTypeLocal:
		return NewLocalVectorStore(config)
	case MemoryStoreTypeSQLite:
		return NewSQLiteVectorStore(config)
	case MemoryStoreTypeChroma:
		// Chroma implements both interfaces
		return NewChromaMemoryStore(config)
	default:
		return nil, fmt.Errorf("unsupported vector store type: %s", storeType)
	}
}

// MemoryManager manages memory stores for agent stacks
type MemoryManager struct {
	factory   MemoryStoreFactory
	stores    map[string]MemoryStore
	vectorStores map[string]VectorStore
}

// NewMemoryManager creates a new memory manager
func NewMemoryManager(factory MemoryStoreFactory) *MemoryManager {
	return &MemoryManager{
		factory:      factory,
		stores:       make(map[string]MemoryStore),
		vectorStores: make(map[string]VectorStore),
	}
}

// GetStore gets or creates a memory store
func (m *MemoryManager) GetStore(name string, storeType MemoryStoreType, config MemoryConfig) (MemoryStore, error) {
	// Check if store already exists
	storeKey := fmt.Sprintf("%s-%s", name, storeType)
	if store, ok := m.stores[storeKey]; ok {
		return store, nil
	}

	// Set default collection name
	if config.CollectionName == "" {
		config.CollectionName = name
	}

	// Create new store
	store, err := m.factory.Create(storeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	// Save for future use
	m.stores[storeKey] = store

	return store, nil
}

// GetVectorStore gets or creates a vector store
func (m *MemoryManager) GetVectorStore(name string, storeType MemoryStoreType, config MemoryConfig) (VectorStore, error) {
	// Check if store already exists
	storeKey := fmt.Sprintf("%s-%s-vector", name, storeType)
	if store, ok := m.vectorStores[storeKey]; ok {
		return store, nil
	}

	// Set default collection name
	if config.CollectionName == "" {
		config.CollectionName = name
	}

	// Create new store
	store, err := m.factory.CreateVector(storeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	// Save for future use
	m.vectorStores[storeKey] = store

	return store, nil
}

// CloseAll closes all stores
func (m *MemoryManager) CloseAll() error {
	var lastErr error

	// Close all regular stores
	for key, store := range m.stores {
		if err := store.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close store %s: %w", key, err)
		}
		delete(m.stores, key)
	}

	// Close all vector stores
	for key, store := range m.vectorStores {
		if err := store.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close vector store %s: %w", key, err)
		}
		delete(m.vectorStores, key)
	}

	return lastErr
}
