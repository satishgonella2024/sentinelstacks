// Package memory provides memory store implementations
package memory

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// VectorEntry represents a stored vector with its text and metadata
type VectorEntry struct {
	ID       string
	Text     string
	Vector   []float32
	Metadata map[string]interface{}
}

// LocalVectorStore is an in-memory implementation of VectorStore
type LocalVectorStore struct {
	dimensions int
	vectors    map[string]*VectorEntry
	mu         sync.RWMutex
}

// NewLocalVectorStore creates a new local in-memory vector store
func NewLocalVectorStore(dimensions int) *LocalVectorStore {
	return &LocalVectorStore{
		dimensions: dimensions,
		vectors:    make(map[string]*VectorEntry),
	}
}

// StoreVector stores a text with its vector embedding
func (s *LocalVectorStore) StoreVector(ctx context.Context, id string, text string, metadata map[string]interface{}) error {
	// In a real implementation, we would generate embeddings here
	// For now, we'll just create a mock vector
	vector := make([]float32, s.dimensions)

	// Generate a simple mock vector based on the text content
	// This is just for demonstration purposes
	for i := 0; i < s.dimensions && i < len(text); i++ {
		if i < len(text) {
			vector[i] = float32(text[i]) / 255.0
		}
	}

	entry := &VectorEntry{
		ID:       id,
		Text:     text,
		Vector:   vector,
		Metadata: metadata,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors[id] = entry
	return nil
}

// SearchVector finds similar vectors using cosine similarity
func (s *LocalVectorStore) SearchVector(ctx context.Context, text string, limit int, filter map[string]interface{}) ([]types.MemoryMatch, error) {
	if limit <= 0 {
		limit = 10
	}

	// Again, in a real implementation, we would generate embeddings for the query text
	// For now, we'll create a mock query vector
	queryVector := make([]float32, s.dimensions)
	for i := 0; i < s.dimensions && i < len(text); i++ {
		if i < len(text) {
			queryVector[i] = float32(text[i]) / 255.0
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Calculate similarity for all vectors
	type scoredVector struct {
		id    string
		score float32
	}

	var results []scoredVector
	for id, entry := range s.vectors {
		// Check if entry matches filter
		if !matchesFilter(entry.Metadata, filter) {
			continue
		}

		// Calculate cosine similarity
		similarity := cosineSimilarity(queryVector, entry.Vector)

		results = append(results, scoredVector{
			id:    id,
			score: similarity,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// Convert to MemoryMatch objects
	matches := make([]types.MemoryMatch, 0, len(results))
	for i := 0; i < len(results) && i < limit; i++ {
		entry := s.vectors[results[i].id]
		matches = append(matches, types.MemoryMatch{
			Key:       entry.ID,
			Content:   entry.Text,
			Score:     float64(results[i].score),
			Metadata:  entry.Metadata,
			Timestamp: time.Now(),
			Distance:  1.0 - float64(results[i].score), // Convert similarity to distance
		})
	}

	return matches, nil
}

// DeleteVector removes a vector by ID
func (s *LocalVectorStore) DeleteVector(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.vectors[id]; !exists {
		return errors.New("vector not found")
	}

	delete(s.vectors, id)
	return nil
}

// GetVector retrieves a vector by ID
func (s *LocalVectorStore) GetVector(ctx context.Context, id string) (*types.MemoryMatch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.vectors[id]
	if !exists {
		return nil, fmt.Errorf("vector not found: %s", id)
	}

	return &types.MemoryMatch{
		Key:       entry.ID,
		Content:   entry.Text,
		Score:     1.0, // Perfect match for direct retrieval
		Metadata:  entry.Metadata,
		Timestamp: time.Now(),
		Distance:  0.0, // Zero distance for direct retrieval
	}, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	var magnitudeA float32
	var magnitudeB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	magnitudeA = float32(math.Sqrt(float64(magnitudeA)))
	magnitudeB = float32(math.Sqrt(float64(magnitudeB)))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}

// matchesFilter checks if metadata matches the filter criteria
func matchesFilter(metadata, filter map[string]interface{}) bool {
	// Empty filter matches everything
	if len(filter) == 0 {
		return true
	}

	// Check each filter condition
	for key, value := range filter {
		metaValue, exists := metadata[key]
		if !exists || metaValue != value {
			return false
		}
	}

	return true
}

// ListVectors returns all keys in the store
func (s *LocalVectorStore) ListVectors(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for key := range s.vectors {
		keys = append(keys, key)
	}

	return keys, nil
}

// Close releases resources
func (s *LocalVectorStore) Close() error {
	// Nothing to close for an in-memory store
	return nil
}
