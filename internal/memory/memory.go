package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MemoryType defines the type of memory to use
type MemoryType string

const (
	// SimpleMemoryType is a basic key-value memory
	SimpleMemoryType MemoryType = "simple"

	// VectorMemoryType uses vector embeddings for semantic search
	VectorMemoryType MemoryType = "vector"
)

// MemoryConfig defines how memory should be configured
type MemoryConfig struct {
	Type           MemoryType `json:"type" yaml:"type"`
	Persistence    bool       `json:"persistence" yaml:"persistence"`
	MaxItems       int        `json:"maxItems" yaml:"maxItems"`
	EmbeddingModel string     `json:"embeddingModel,omitempty" yaml:"embeddingModel,omitempty"`
}

// DefaultConfig returns a default memory configuration
func DefaultConfig() MemoryConfig {
	return MemoryConfig{
		Type:        SimpleMemoryType,
		Persistence: true,
		MaxItems:    1000,
	}
}

// DefaultVectorConfig returns a default vector memory configuration
func DefaultVectorConfig() MemoryConfig {
	return MemoryConfig{
		Type:           VectorMemoryType,
		Persistence:    true,
		MaxItems:       1000,
		EmbeddingModel: "openai:text-embedding-3-small",
	}
}

// MemoryEntry represents a single memory entry
type MemoryEntry struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Memory defines the interface for all memory implementations
type Memory interface {
	// Add adds a new entry to memory
	Add(content string, metadata map[string]interface{}) (string, error)

	// Get retrieves an entry by ID
	Get(id string) (*MemoryEntry, error)

	// Search finds entries that match the query
	Search(query string, limit int) ([]MemoryEntry, error)

	// List returns all entries, optionally limited
	List(limit int) ([]MemoryEntry, error)

	// Delete removes an entry by ID
	Delete(id string) error

	// Clear removes all entries
	Clear() error

	// Save persists the memory to disk
	Save() error

	// Load loads the memory from disk
	Load() error
}

// NewMemory creates a new memory instance based on configuration
func NewMemory(agentName string, config MemoryConfig) (Memory, error) {
	switch config.Type {
	case SimpleMemoryType:
		return NewSimpleMemory(agentName, config)
	case VectorMemoryType:
		return NewVectorMemory(agentName, config)
	default:
		return nil, fmt.Errorf("unsupported memory type: %s", config.Type)
	}
}

// getMemoryStoragePath returns the path where memory files are stored
func getMemoryStoragePath(agentName string) string {
	base := filepath.Join(os.Getenv("HOME"), ".sentinel", "memory")
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(base, 0755); err != nil {
		fmt.Printf("Error creating memory directory: %v\n", err)
	}
	return filepath.Join(base, fmt.Sprintf("%s.json", agentName))
}

// getVectorStoragePath returns the path where vector index files are stored
func getVectorStoragePath(agentName string) string {
	base := filepath.Join(os.Getenv("HOME"), ".sentinel", "vectors")
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(base, 0755); err != nil {
		fmt.Printf("Error creating vectors directory: %v\n", err)
	}
	return filepath.Join(base, agentName)
}

// saveToFile saves data to a file
func saveToFile(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// loadFromFile loads data from a file
func loadFromFile(path string, target interface{}) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, that's fine for a new agent
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	return nil
}
