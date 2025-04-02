package memory

import (
	"context"
	"time"
)

// MemoryStore defines the interface for memory storage backends
type MemoryStore interface {
	// Save stores a value with the given key
	Save(ctx context.Context, key string, value interface{}) error
	
	// Load retrieves a value by key
	Load(ctx context.Context, key string) (interface{}, error)
	
	// Delete removes a key-value pair
	Delete(ctx context.Context, key string) error
	
	// Clear removes all keys and values
	Clear(ctx context.Context) error
	
	// Keys returns all keys in the store
	Keys(ctx context.Context) ([]string, error)
	
	// Close closes the memory store
	Close() error
}

// VectorStore extends MemoryStore with vector operations
type VectorStore interface {
	MemoryStore
	
	// SaveEmbedding stores a vector embedding
	SaveEmbedding(ctx context.Context, key string, vector []float32, metadata map[string]interface{}) error
	
	// Query performs a similarity search on stored embeddings
	Query(ctx context.Context, vector []float32, topK int) ([]SimilarityMatch, error)
	
	// DeleteEmbedding removes an embedding
	DeleteEmbedding(ctx context.Context, key string) error
}

// SimilarityMatch represents a match from a vector similarity search
type SimilarityMatch struct {
	Key       string                 `json:"key"`
	Score     float32                `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryEntry represents a stored memory item
type MemoryEntry struct {
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// MemoryConfig contains configuration for memory stores
type MemoryConfig struct {
	// ConnectionString is the connection string for the backend
	ConnectionString string `json:"connectionString"`
	
	// CollectionName is the name of the collection or table
	CollectionName string `json:"collectionName"`
	
	// Namespace is an optional prefix for keys
	Namespace string `json:"namespace"`
	
	// TTL is the time-to-live for entries (0 = no expiration)
	TTL time.Duration `json:"ttl"`
	
	// VectorDimensions is the size of vector embeddings
	VectorDimensions int `json:"vectorDimensions"`
	
	// AdditionalOptions contains backend-specific options
	AdditionalOptions map[string]interface{} `json:"additionalOptions,omitempty"`
}

// MemoryStoreType identifies the type of memory store
type MemoryStoreType string

const (
	// MemoryStoreTypeLocal is an in-memory store
	MemoryStoreTypeLocal MemoryStoreType = "local"
	
	// MemoryStoreTypeSQLite is a SQLite-backed store
	MemoryStoreTypeSQLite MemoryStoreType = "sqlite"
	
	// MemoryStoreTypeChroma is a Chroma vector database store
	MemoryStoreTypeChroma MemoryStoreType = "chroma"
	
	// MemoryStoreTypeRedis is a Redis-backed store
	MemoryStoreTypeRedis MemoryStoreType = "redis"
	
	// MemoryStoreTypePostgres is a PostgreSQL-backed store
	MemoryStoreTypePostgres MemoryStoreType = "postgres"
)

// MemoryStoreFactory creates memory stores
type MemoryStoreFactory interface {
	// Create creates a memory store of the specified type
	Create(storeType MemoryStoreType, config MemoryConfig) (MemoryStore, error)
	
	// CreateVector creates a vector store of the specified type
	CreateVector(storeType MemoryStoreType, config MemoryConfig) (VectorStore, error)
}
