// Package memory provides memory store implementations
package memory

import (
	"context"
	"fmt"
	"sync"
)

// LocalStore is an in-memory implementation of MemoryStore
type LocalStore struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewLocalStore creates a new local in-memory store
func NewLocalStore() *LocalStore {
	return &LocalStore{
		data: make(map[string]interface{}),
	}
}

// Save stores a value by key
func (s *LocalStore) Save(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

// Load retrieves a value by key
func (s *LocalStore) Load(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

// Delete removes a value by key
func (s *LocalStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return nil
}

// List returns all keys
func (s *LocalStore) List(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for key := range s.data {
		keys = append(keys, key)
	}

	return keys, nil
}

// hasPrefix checks if a string has the given prefix
func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

// Close releases resources
func (s *LocalStore) Close() error {
	// Nothing to close for an in-memory store
	return nil
}
