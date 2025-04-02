package memory

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// LocalMemoryStore is an in-memory implementation of MemoryStore
type LocalMemoryStore struct {
	name      string
	data      map[string]MemoryEntry
	namespace string
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewLocalMemoryStore creates a new local memory store
func NewLocalMemoryStore(config MemoryConfig) (*LocalMemoryStore, error) {
	name := "local-memory"
	if config.CollectionName != "" {
		name = config.CollectionName
	}
	
	return &LocalMemoryStore{
		name:      name,
		data:      make(map[string]MemoryEntry),
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
	entry := MemoryEntry{
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

// Clear removes all keys and values
func (s *LocalMemoryStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// If namespace is specified, only clear entries with that namespace
	if s.namespace != "" {
		prefix := s.namespace + ":"
		for key := range s.data {
			if strings.HasPrefix(key, prefix) {
				delete(s.data, key)
			}
		}
	} else {
		// Otherwise clear all entries
		s.data = make(map[string]MemoryEntry)
	}
	
	return nil
}

// Keys returns all keys in the store
func (s *LocalMemoryStore) Keys(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var keys []string
	
	// If namespace is specified, only return keys with that namespace
	if s.namespace != "" {
		prefix := s.namespace + ":"
		prefixLen := len(prefix)
		
		for key := range s.data {
			if strings.HasPrefix(key, prefix) {
				// Remove namespace prefix from returned keys
				keys = append(keys, key[prefixLen:])
			}
		}
	} else {
		// Otherwise return all keys
		keys = make([]string, 0, len(s.data))
		for key := range s.data {
			keys = append(keys, key)
		}
	}
	
	return keys, nil
}

// Close closes the memory store
func (s *LocalMemoryStore) Close() error {
	// No resources to release for in-memory store
	return nil
}

// LocalVectorStore is an in-memory implementation of VectorStore
type LocalVectorStore struct {
	*LocalMemoryStore
	vectors map[string][]float32
	maxDim  int
}

// NewLocalVectorStore creates a new local vector store
func NewLocalVectorStore(config MemoryConfig) (*LocalVectorStore, error) {
	base, err := NewLocalMemoryStore(config)
	if err != nil {
		return nil, err
	}
	
	// Default vector dimensions
	maxDim := 1536
	if config.VectorDimensions > 0 {
		maxDim = config.VectorDimensions
	}
	
	return &LocalVectorStore{
		LocalMemoryStore: base,
		vectors:          make(map[string][]float32),
		maxDim:           maxDim,
	}, nil
}

// SaveEmbedding stores a vector embedding
func (s *LocalVectorStore) SaveEmbedding(ctx context.Context, key string, vector []float32, metadata map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Validate vector dimensions
	if len(vector) > s.maxDim {
		return fmt.Errorf("vector dimensions exceed maximum (%d > %d)", len(vector), s.maxDim)
	}
	
	// Add namespace prefix if specified
	if s.namespace != "" {
		key = s.namespace + ":" + key
	}
	
	// Create metadata if nil
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	
	// Store vector
	s.vectors[key] = vector
	
	// Store metadata
	now := time.Now()
	entry := MemoryEntry{
		Key:       key,
		Value:     metadata,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.data[key] = entry
	
	return nil
}

// Query performs a similarity search on stored embeddings
func (s *LocalVectorStore) Query(ctx context.Context, vector []float32, topK int) ([]SimilarityMatch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Validate vector dimensions
	if len(vector) > s.maxDim {
		return nil, fmt.Errorf("vector dimensions exceed maximum (%d > %d)", len(vector), s.maxDim)
	}
	
	// Check if there are any vectors
	if len(s.vectors) == 0 {
		return []SimilarityMatch{}, nil
	}
	
	// Compute similarity scores for all vectors
	scores := make([]SimilarityMatch, 0, len(s.vectors))
	
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
		entry, ok := s.data[key]
		var metadata map[string]interface{}
		if ok {
			metadata = entry.Metadata
		} else {
			metadata = map[string]interface{}{}
		}
		
		// Remove namespace prefix from returned key
		returnKey := key
		if s.namespace != "" && strings.HasPrefix(key, prefix) {
			returnKey = key[len(prefix):]
		}
		
		// Add to scores
		scores = append(scores, SimilarityMatch{
			Key:      returnKey,
			Score:    score,
			Metadata: metadata,
		})
	}
	
	// Sort by score (descending)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	
	// Limit to topK
	if topK > 0 && topK < len(scores) {
		scores = scores[:topK]
	}
	
	return scores, nil
}

// DeleteEmbedding removes an embedding
func (s *LocalVectorStore) DeleteEmbedding(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Add namespace prefix if specified
	if s.namespace != "" {
		key = s.namespace + ":" + key
	}
	
	// Delete vector
	delete(s.vectors, key)
	
	// Delete metadata
	delete(s.data, key)
	
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
