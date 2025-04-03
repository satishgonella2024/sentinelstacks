// Package types defines common type definitions used across packages
package types

import (
	"context"
	"time"
)

// MemoryEntry represents a single entry in a memory store
type MemoryEntry struct {
	Key       string
	Value     interface{}
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemoryStoreType represents the type of memory store
type MemoryStoreType string

const (
	// MemoryStoreTypeLocal is an in-memory store
	MemoryStoreTypeLocal MemoryStoreType = "local"
	// MemoryStoreTypeSQLite is a SQLite-backed store
	MemoryStoreTypeSQLite MemoryStoreType = "sqlite"
	// MemoryStoreTypeChroma is a Chroma vector store
	MemoryStoreTypeChroma MemoryStoreType = "chroma"
)

// MemoryConfig defines configuration for a memory store
type MemoryConfig struct {
	// StorePath is the path to the memory store
	StorePath string

	// Embeddings defines the embedding configuration
	Embeddings EmbeddingConfig

	// MetadataFields are fields to index for metadata-based search
	MetadataFields []string

	// AdditionalOptions contains store-specific configuration
	AdditionalOptions map[string]interface{}
}

// EmbeddingConfig defines configuration for embeddings
type EmbeddingConfig struct {
	// Provider is the name of the embedding provider
	Provider string

	// Model is the name of the embedding model
	Model string

	// Dimensions is the dimensionality of embeddings
	Dimensions int

	// AdditionalOptions contains provider-specific configuration
	AdditionalOptions map[string]interface{}
}

// MemoryStore provides a generic interface for memory storage
type MemoryStore interface {
	// Save stores a value with the given key
	Save(ctx context.Context, key string, value interface{}) error

	// Load retrieves a value by key
	Load(ctx context.Context, key string) (interface{}, error)

	// Delete removes a value by key
	Delete(ctx context.Context, key string) error

	// List returns all keys in the store
	List(ctx context.Context) ([]string, error)

	// Close releases resources
	Close() error
}

// MemoryMatch represents a match from a similarity search
type MemoryMatch struct {
	// Key is the unique identifier of the match
	Key string

	// Content is the matched text content
	Content string

	// Metadata contains additional information about the match
	Metadata map[string]interface{}

	// Score is the similarity score (0-1)
	Score float64

	// Distance is the vector distance (lower means more similar)
	Distance float64

	// Timestamp is when the content was stored
	Timestamp time.Time
}

// VectorStore provides a generic interface for vector storage
type VectorStore interface {
	// StoreVector stores a text with its vector embedding
	StoreVector(ctx context.Context, key string, text string, metadata map[string]interface{}) error

	// SearchVector finds similar texts using vector similarity
	SearchVector(ctx context.Context, text string, limit int, filter map[string]interface{}) ([]MemoryMatch, error)

	// GetVector retrieves a stored vector by key
	GetVector(ctx context.Context, key string) (*MemoryMatch, error)

	// DeleteVector removes a vector by key
	DeleteVector(ctx context.Context, key string) error

	// ListVectors returns all keys in the store
	ListVectors(ctx context.Context) ([]string, error)

	// Close releases resources
	Close() error
}

// MemoryStoreFactory creates memory stores
type MemoryStoreFactory interface {
	// CreateMemoryStore creates a basic key-value memory store
	CreateMemoryStore(ctx context.Context, name string) (MemoryStore, error)

	// CreateVectorStore creates a vector-based memory store
	CreateVectorStore(ctx context.Context, name string) (VectorStore, error)
}
