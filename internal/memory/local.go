package memory

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// LocalMemoryStore is an in-memory implementation of MemoryStore
type LocalMemoryStore struct {
	name      string
	data      map[string]types.MemoryEntry
	namespace string
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewLocalMemoryStore creates a new local memory store
func NewLocalMemoryStore(config types.MemoryConfig) (*LocalMemoryStore, error) {
	name := "local-memory"
	if config.CollectionName != "" {
		name = config.CollectionName
	}

	return &LocalMemoryStore{
		name:      name,
		data:      make(map[string]types.MemoryEntry),
		namespace: config.Namespace,
		ttl:       config.TTL,
	}, nil
}

// Save stores a value with the given key
func (s *LocalMemoryStore) Save(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add namespace prefix if specified
	if s.namespace != "" {
		key = s.namespace + ":" + key
	}

	// Create entry
	now := time.Now()
	entry := types.MemoryEntry{
		Key:       key,
		Value:     value,
		Metadata:  map[string]interface{}{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store entry
	s.data[key] = entry

	return nil
}

// Load retrieves a value by key
func (s *LocalMemoryStore) Load(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Add namespace prefix if specified
	if s.namespace != "" {
		key = s.namespace + ":" + key
	}

	// Check if entry exists
	entry, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Check if entry has expired
	if s.ttl > 0 && time.Since(entry.UpdatedAt) > s.ttl {
		// Remove expired entry
		go func() {
			s.mu.Lock()
			delete(s.data, key)
			s.mu.Unlock()
		}()
		return nil, fmt.Errorf("key expired: %s", key)
	}

	return entry.Value, nil
}

// Delete removes a key-value pair
func (s *LocalMemoryStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add namespace prefix if specified
	if s.namespace != "" {
		key = s.namespace + ":" + key
	}

	// Delete entry
	delete(s.data, key)

	return nil
}

// List returns all keys with optional prefix filtering
func (s *LocalMemoryStore) List(ctx context.Context, prefix string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Prepare namespace prefix
	namespacePrefix := ""
	if s.namespace != "" {
		namespacePrefix = s.namespace + ":"
	}

	// Combine namespace prefix with query prefix
	queryPrefix := namespacePrefix + prefix

	// Collect matching keys
	keys := make([]string, 0)
	for k := range s.data {
		// Check if key has prefix
		if strings.HasPrefix(k, queryPrefix) {
			// Check if entry has expired
			if s.ttl > 0 && time.Since(s.data[k].UpdatedAt) > s.ttl {
				// Skip expired entries
				continue
			}

			// Remove namespace prefix for returned keys
			key := k
			if namespacePrefix != "" && strings.HasPrefix(k, namespacePrefix) {
				key = k[len(namespacePrefix):]
			}

			keys = append(keys, key)
		}
	}

	// Sort keys for deterministic output
	sort.Strings(keys)

	return keys, nil
}

// Clear removes all data for this namespace
func (s *LocalMemoryStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If namespace is specified, only clear data for that namespace
	if s.namespace != "" {
		prefix := s.namespace + ":"
		for k := range s.data {
			if strings.HasPrefix(k, prefix) {
				delete(s.data, k)
			}
		}
	} else {
		// Otherwise clear all data
		s.data = make(map[string]types.MemoryEntry)
	}

	return nil
}

// Close releases all resources
func (s *LocalMemoryStore) Close() error {
	// Nothing to close for in-memory store
	return nil
}

// LocalVectorStore is an in-memory implementation of VectorStore
type LocalVectorStore struct {
	namespace string
	vectors   map[string][]float32
	metadata  map[string]map[string]interface{}
	maxDim    int
	mu        sync.RWMutex
}

// NewLocalVectorStore creates a new local vector store
func NewLocalVectorStore(config types.MemoryConfig) (*LocalVectorStore, error) {
	// Default vector dimensions
	maxDim := 1536
	if config.VectorDimensions > 0 {
		maxDim = config.VectorDimensions
	}

	return &LocalVectorStore{
		namespace: config.Namespace,
		vectors:   make(map[string][]float32),
		metadata:  make(map[string]map[string]interface{}),
		maxDim:    maxDim,
	}, nil
}

// StoreVector saves a vector with metadata
func (s *LocalVectorStore) StoreVector(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate vector dimensions
	if len(vector) > s.maxDim {
		return fmt.Errorf("vector dimensions exceed maximum (%d > %d)", len(vector), s.maxDim)
	}

	// Add namespace prefix if specified
	storeID := id
	if s.namespace != "" {
		storeID = s.namespace + ":" + id
	}

	// Create metadata if nil
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	// Store vector and metadata
	s.vectors[storeID] = vector
	s.metadata[storeID] = metadata

	return nil
}

// FindSimilar finds similar vectors
func (s *LocalVectorStore) FindSimilar(ctx context.Context, vector []float32, limit int) ([]types.SimilarityResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Validate vector dimensions
	if len(vector) > s.maxDim {
		return nil, fmt.Errorf("vector dimensions exceed maximum (%d > %d)", len(vector), s.maxDim)
	}

	// Check if there are any vectors
	if len(s.vectors) == 0 {
		return []types.SimilarityResult{}, nil
	}

	// Compute similarity scores for all vectors
	scores := make([]types.SimilarityResult, 0, len(s.vectors))

	prefix := ""
	if s.namespace != "" {
		prefix = s.namespace + ":"
	}

	for key, vec := range s.vectors {
		// Skip if key doesn't match namespace
		if s.namespace != "" && !strings.HasPrefix(key, prefix) {
			continue
		}

		// Compute cosine similarity
		score := cosineSimilarity(vector, vec)

		// Get metadata
		var metadata map[string]interface{}
		if metaMap, ok := s.metadata[key]; ok {
			metadata = metaMap
		} else {
			metadata = map[string]interface{}{}
		}

		// Remove namespace prefix from returned key
		returnKey := key
		if s.namespace != "" && strings.HasPrefix(key, prefix) {
			returnKey = key[len(prefix):]
		}

		// Add to scores
		scores = append(scores, types.SimilarityResult{
			ID:       returnKey,
			Score:    score,
			Metadata: metadata,
		})
	}

	// Sort by score (descending)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Limit to topK
	if limit > 0 && limit < len(scores) {
		scores = scores[:limit]
	}

	return scores, nil
}

// GetVector gets a vector by ID
func (s *LocalVectorStore) GetVector(ctx context.Context, id string) ([]float32, map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Add namespace prefix if specified
	storeID := id
	if s.namespace != "" {
		storeID = s.namespace + ":" + id
	}

	// Get vector
	vector, ok := s.vectors[storeID]
	if !ok {
		return nil, nil, fmt.Errorf("vector not found: %s", id)
	}

	// Get metadata
	metadata, ok := s.metadata[storeID]
	if !ok {
		metadata = make(map[string]interface{})
	}

	return vector, metadata, nil
}

// DeleteVector removes a vector
func (s *LocalVectorStore) DeleteVector(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add namespace prefix if specified
	storeID := id
	if s.namespace != "" {
		storeID = s.namespace + ":" + id
	}

	// Check if vector exists
	if _, ok := s.vectors[storeID]; !ok {
		return fmt.Errorf("vector not found: %s", id)
	}

	// Delete vector and metadata
	delete(s.vectors, storeID)
	delete(s.metadata, storeID)

	return nil
}

// Clear removes all vectors
func (s *LocalVectorStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If namespace is specified, only clear vectors for that namespace
	if s.namespace != "" {
		prefix := s.namespace + ":"
		for key := range s.vectors {
			if strings.HasPrefix(key, prefix) {
				delete(s.vectors, key)
				delete(s.metadata, key)
			}
		}
	} else {
		// Otherwise clear all vectors
		s.vectors = make(map[string][]float32)
		s.metadata = make(map[string]map[string]interface{})
	}

	return nil
}

// Close releases all resources
func (s *LocalVectorStore) Close() error {
	// Nothing to close for in-memory store
	return nil
}

// Helper function for cosine similarity calculation
func cosineSimilarity(a, b []float32) float32 {
	// Use the minimum length
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	// Calculate dot product and magnitudes
	var dotProduct, magnitudeA, magnitudeB float32

	for i := 0; i < length; i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	// Handle zero magnitude
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	// Return cosine similarity
	return dotProduct / (float32(math.Sqrt(float64(magnitudeA))) * float32(math.Sqrt(float64(magnitudeB))))
}
