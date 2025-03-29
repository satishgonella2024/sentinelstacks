package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SimpleMemory is a basic implementation of Memory
type SimpleMemory struct {
	AgentName  string                  `json:"agentName"`
	Entries    map[string]MemoryEntry  `json:"entries"`
	Config     MemoryConfig            `json:"config"`
	mu         sync.RWMutex
}

// NewSimpleMemory creates a new SimpleMemory instance
func NewSimpleMemory(agentName string, config MemoryConfig) (*SimpleMemory, error) {
	memory := &SimpleMemory{
		AgentName: agentName,
		Entries:   make(map[string]MemoryEntry),
		Config:    config,
	}
	
	// Load from disk if persistence is enabled
	if config.Persistence {
		if err := memory.Load(); err != nil {
			return nil, fmt.Errorf("error loading memory: %w", err)
		}
	}
	
	return memory, nil
}

// Add adds a new entry to memory
func (m *SimpleMemory) Add(content string, metadata map[string]interface{}) (string, error) {
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
func (m *SimpleMemory) Get(id string) (*MemoryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	entry, ok := m.Entries[id]
	if !ok {
		return nil, fmt.Errorf("entry not found: %s", id)
	}
	
	return &entry, nil
}

// Search finds entries that match the query
func (m *SimpleMemory) Search(query string, limit int) ([]MemoryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var results []MemoryEntry
	
	// Simple string matching for now
	// For the real implementation, you might want to use a more sophisticated
	// search algorithm, especially for SimpleMemory
	query = strings.ToLower(query)
	for _, entry := range m.Entries {
		if strings.Contains(strings.ToLower(entry.Content), query) {
			results = append(results, entry)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	
	return results, nil
}

// List returns all entries, optionally limited
func (m *SimpleMemory) List(limit int) ([]MemoryEntry, error) {
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
func (m *SimpleMemory) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, ok := m.Entries[id]; !ok {
		return fmt.Errorf("entry not found: %s", id)
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
func (m *SimpleMemory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.Entries = make(map[string]MemoryEntry)
	
	// Save to disk if persistence is enabled
	if m.Config.Persistence {
		if err := m.Save(); err != nil {
			return fmt.Errorf("error saving memory: %w", err)
		}
	}
	
	return nil
}

// Save persists the memory to disk
func (m *SimpleMemory) Save() error {
	path := getMemoryStoragePath(m.AgentName)
	return saveToFile(path, m)
}

// Load loads the memory from disk
func (m *SimpleMemory) Load() error {
	path := getMemoryStoragePath(m.AgentName)
	return loadFromFile(path, m)
}

// pruneOldest removes the oldest entries when we exceed maxItems
func (m *SimpleMemory) pruneOldest() {
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
	// This is fine for small numbers, but you might want to optimize
	// for larger memory stores
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
		delete(m.Entries, entries[oldestIdx].id)
		
		// Remove from our list
		entries = append(entries[:oldestIdx], entries[oldestIdx+1:]...)
	}
}