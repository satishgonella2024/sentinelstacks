package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/satishgonella2024/sentinelstacks/internal/vector"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
)

// VectorMemory implements Memory using vector embeddings
type VectorMemory struct {
	AgentName         string                 `json:"agentName"`
	Entries           map[string]MemoryEntry `json:"entries"`
	Config            MemoryConfig           `json:"config"`
	VectorDB          *vector.VectorIndex    `json:"-"`
	EmbeddingProvider models.EmbeddingProvider `json:"-"`
	mu                sync.RWMutex
}

// NewVectorMemory creates a new VectorMemory instance
func NewVectorMemory(agentName string, config MemoryConfig) (*VectorMemory, error) {
	memory := &VectorMemory{
		AgentName: agentName,
		Entries:   make(map[string]MemoryEntry),
		Config:    config,
	}

	// Get embedding provider based on configuration
	embeddingModel := "openai:text-embedding-3-small" // Default
	if config.EmbeddingModel != "" {
		embeddingModel = config.EmbeddingModel
	}

	provider, err := models.GetEmbeddingProvider(embeddingModel)
	if err != nil {
		return nil, fmt.Errorf("error creating embedding provider: %w", err)
	}
	memory.EmbeddingProvider = provider

	// Initialize vector database
	vectorPath := getVectorStoragePath(agentName)
	vectorDB, err := vector.NewVectorIndex(vectorPath)
	if err != nil {
		return nil, fmt.Errorf("error initializing vector database: %w", err)
	}
	memory.VectorDB = vectorDB

	// Load from disk if persistence is enabled
	if config.Persistence {
		if err := memory.Load(); err != nil {
			return nil, fmt.Errorf("error loading memory: %w", err)
		}
	}

	return memory, nil
}

// Add adds a new entry to memory
func (m *VectorMemory) Add(content string, metadata map[string]interface{}) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	entry := MemoryEntry{
		ID:        id,
		Content:   content,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	m.Entries[id] = entry

	// Create embedding using the provider
	embedding, err := m.EmbeddingProvider.GetEmbedding(content)
	if err != nil {
		// Fall back to a dummy embedding if the provider fails
		embedding = createDummyEmbedding(content)
	}

	// Add to vector database
	vectorID, err := m.VectorDB.Add(embedding, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return "", fmt.Errorf("error adding to vector database: %w", err)
	}

	// Store vector ID in metadata
	if entry.Metadata == nil {
		entry.Metadata = make(map[string]interface{})
	}
	entry.Metadata["vector_id"] = vectorID
	m.Entries[id] = entry

	// Enforce max items limit
	if m.Config.MaxItems > 0 && len(m.Entries) > m.Config.MaxItems {
		m.pruneOldest()
	}

	// Save to disk if persistence is enabled
	if m.Config.Persistence {
		if err := m.Save(); err != nil {
			return id, fmt.Errorf("error saving memory: %w", err)
		}
	}

	return id, nil
}

// Get retrieves an entry by ID
func (m *VectorMemory) Get(id string) (*MemoryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.Entries[id]
	if !ok {
		return nil, fmt.Errorf("entry not found: %s", id)
	}

	return &entry, nil
}

// Search finds entries that match the query
func (m *VectorMemory) Search(query string, limit int) ([]MemoryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create embedding for query using the provider
	queryEmbedding, err := m.EmbeddingProvider.GetEmbedding(query)
	if err != nil {
		// Fall back to a dummy embedding if the provider fails
		queryEmbedding = createDummyEmbedding(query)
	}

	// Search vector database
	ids, scores, err := m.VectorDB.Search(queryEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("error searching vector database: %w", err)
	}

	// Retrieve entries
	var results []MemoryEntry
	for i, id := range ids {
		// Get the memory entry ID from the vector metadata
		_, metadata, err := m.VectorDB.Get(id)
		if err != nil {
			continue
		}

		memoryID, ok := metadata["id"].(string)
		if !ok {
			continue
		}

		entry, ok := m.Entries[memoryID]
		if !ok {
			continue
		}

		// Add similarity score to metadata
		if entry.Metadata == nil {
			entry.Metadata = make(map[string]interface{})
		}
		entry.Metadata["similarity"] = scores[i]

		results = append(results, entry)
	}

	return results, nil
}

// List returns all entries, optionally limited
func (m *VectorMemory) List(limit int) ([]MemoryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]MemoryEntry, 0, len(m.Entries))
	for _, entry := range m.Entries {
		results = append(results, entry)
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, nil
}

// Delete removes an entry by ID
func (m *VectorMemory) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.Entries[id]
	if !ok {
		return fmt.Errorf("entry not found: %s", id)
	}

	// Delete from vector database if we have a vector ID
	if entry.Metadata != nil {
		if vectorID, ok := entry.Metadata["vector_id"].(string); ok {
			if err := m.VectorDB.Delete(vectorID); err != nil {
				return fmt.Errorf("error deleting from vector database: %w", err)
			}
		}
	}

	delete(m.Entries, id)

	// Save to disk if persistence is enabled
	if m.Config.Persistence {
		if err := m.Save(); err != nil {
			return fmt.Errorf("error saving memory: %w", err)
		}
	}

	return nil
}

// Clear removes all entries
func (m *VectorMemory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Entries = make(map[string]MemoryEntry)

	// Clear vector database
	if err := m.VectorDB.Clear(); err != nil {
		return fmt.Errorf("error clearing vector database: %w", err)
	}

	// Save to disk if persistence is enabled
	if m.Config.Persistence {
		if err := m.Save(); err != nil {
			return fmt.Errorf("error saving memory: %w", err)
		}
	}

	return nil
}

// Save persists the memory to disk
func (m *VectorMemory) Save() error {
	path := getMemoryStoragePath(m.AgentName)
	return saveToFile(path, m)
}

// Load loads the memory from disk
func (m *VectorMemory) Load() error {
	path := getMemoryStoragePath(m.AgentName)
	return loadFromFile(path, m)
}

// pruneOldest removes the oldest entries when we exceed maxItems
func (m *VectorMemory) pruneOldest() {
	// Find the oldest entries
	type timestampedID struct {
		id        string
		timestamp time.Time
	}

	var entries []timestampedID
	for id, entry := range m.Entries {
		entries = append(entries, timestampedID{id, entry.Timestamp})
	}

	// Sort by timestamp (oldest first)
	numToRemove := len(m.Entries) - m.Config.MaxItems
	if numToRemove <= 0 {
		return
	}

	// Simple n^2 algorithm to find the oldest entries
	for i := 0; i < numToRemove; i++ {
		oldestIdx := 0
		oldestTime := entries[0].timestamp

		for j := 1; j < len(entries); j++ {
			if entries[j].timestamp.Before(oldestTime) {
				oldestIdx = j
				oldestTime = entries[j].timestamp
			}
		}

		// Remove the oldest entry
		entry := m.Entries[entries[oldestIdx].id]
		if entry.Metadata != nil {
			if vectorID, ok := entry.Metadata["vector_id"].(string); ok {
				_ = m.VectorDB.Delete(vectorID) // Ignore errors
			}
		}
		delete(m.Entries, entries[oldestIdx].id)

		// Remove from our list
		entries = append(entries[:oldestIdx], entries[oldestIdx+1:]...)
	}
}

// createDummyEmbedding creates a simple embedding for testing
// This is used as a fallback when the embedding provider fails
func createDummyEmbedding(text string) vector.Embedding {
	// Create a very simple embedding based on character frequencies
	// This is NOT a real embedding, just for demonstration
	embedding := make(vector.Embedding, 1536) // Use OpenAI's text-embedding-3-small dimensions

	// Fill with small random values
	for i := range embedding {
		embedding[i] = float32(i % 10) * 0.1
	}

	// Use some characteristics of the text
	for i, c := range text {
		if i >= len(embedding) {
			break
		}
		embedding[i] = float32(c % 10) * 0.1
	}

	return embedding
}
