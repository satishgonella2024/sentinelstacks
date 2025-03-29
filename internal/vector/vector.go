package vector

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/google/uuid"
)

// Embedding represents a vector embedding
type Embedding []float32

// VectorMetadata contains metadata for vector storage
type VectorMetadata map[string]interface{}

// VectorItem represents a vector with its ID and metadata
type VectorItem struct {
	ID       string
	Vector   Embedding
	Metadata VectorMetadata
}

// ScoredVector represents a vector with a similarity score
type ScoredVector struct {
	ID       string
	Score    float32
	Metadata VectorMetadata
}

// SimilarityFunc is a function that compares two vectors and returns a similarity score
type SimilarityFunc func(Embedding, Embedding) float32

// VectorIndex is a simple vector index for searching embeddings
type VectorIndex struct {
	// Path to the index file
	Path string

	// Map of ID to embedding
	Vectors map[string]Embedding

	// Map of ID to metadata
	Metadata map[string]VectorMetadata

	// Similarity function to use
	SimilarityFunc SimilarityFunc

	// Mutex for thread safety
	mu sync.RWMutex
}

// VectorIndexConfig contains configuration for the vector index
type VectorIndexConfig struct {
	SimilarityType string // "cosine", "dot", "euclidean"
}

// DefaultVectorIndexConfig returns the default vector index configuration
func DefaultVectorIndexConfig() VectorIndexConfig {
	return VectorIndexConfig{
		SimilarityType: "cosine",
	}
}

// NewVectorIndex creates a new vector index
func NewVectorIndex(path string) (*VectorIndex, error) {
	return NewVectorIndexWithConfig(path, DefaultVectorIndexConfig())
}

// NewVectorIndexWithConfig creates a new vector index with the specified configuration
func NewVectorIndexWithConfig(path string, config VectorIndexConfig) (*VectorIndex, error) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Choose similarity function based on configuration
	var similarityFunc SimilarityFunc
	switch config.SimilarityType {
	case "cosine":
		similarityFunc = cosineSimilarity
	case "dot":
		similarityFunc = dotProduct
	case "euclidean":
		similarityFunc = euclideanDistance
	default:
		similarityFunc = cosineSimilarity
	}

	idx := &VectorIndex{
		Path:           path,
		Vectors:        make(map[string]Embedding),
		Metadata:       make(map[string]VectorMetadata),
		SimilarityFunc: similarityFunc,
	}

	// Load the index if it exists
	if _, err := os.Stat(path); err == nil {
		if err := idx.Load(); err != nil {
			return nil, fmt.Errorf("failed to load index: %w", err)
		}
	}

	return idx, nil
}

// Add adds a vector to the index
func (idx *VectorIndex) Add(vector Embedding, metadata VectorMetadata) (string, error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	id := uuid.New().String()
	idx.Vectors[id] = vector
	idx.Metadata[id] = metadata

	// Save the index
	if err := idx.Save(); err != nil {
		return "", fmt.Errorf("failed to save index: %w", err)
	}

	return id, nil
}

// Get retrieves a vector from the index
func (idx *VectorIndex) Get(id string) (Embedding, VectorMetadata, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	vector, ok := idx.Vectors[id]
	if !ok {
		return nil, nil, fmt.Errorf("vector not found: %s", id)
	}

	metadata := idx.Metadata[id]
	return vector, metadata, nil
}

// Search finds the nearest neighbors to the query vector
func (idx *VectorIndex) Search(query Embedding, limit int) ([]string, []float32, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Compute similarities to all vectors
	var scored []ScoredVector
	for id, vector := range idx.Vectors {
		score := idx.SimilarityFunc(query, vector)
		scored = append(scored, ScoredVector{
			ID:       id,
			Score:    score,
			Metadata: idx.Metadata[id],
		})
	}

	// Sort by score (higher is better for cosine and dot product, lower is better for euclidean)
	// The SimilarityFunc should return higher values for more similar items
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Limit results
	if limit > len(scored) {
		limit = len(scored)
	}

	// Extract IDs and scores
	ids := make([]string, limit)
	scores := make([]float32, limit)
	for i := 0; i < limit; i++ {
		ids[i] = scored[i].ID
		scores[i] = scored[i].Score
	}

	return ids, scores, nil
}

// Delete removes a vector from the index
func (idx *VectorIndex) Delete(id string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if _, ok := idx.Vectors[id]; !ok {
		return fmt.Errorf("vector not found: %s", id)
	}

	delete(idx.Vectors, id)
	delete(idx.Metadata, id)

	// Save the index
	if err := idx.Save(); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}

// Clear removes all vectors from the index
func (idx *VectorIndex) Clear() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.Vectors = make(map[string]Embedding)
	idx.Metadata = make(map[string]VectorMetadata)

	// Save the index
	if err := idx.Save(); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}

// Save persists the index to disk
func (idx *VectorIndex) Save() error {
	// Create temporary file
	tempFile := idx.Path + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary index file: %w", err)
	}
	defer file.Close()

	// Write header identifier
	if _, err := file.WriteString("SENTINEL_VECTOR_INDEX_V1\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write number of vectors
	numVectors := uint32(len(idx.Vectors))
	if err := binary.Write(file, binary.LittleEndian, numVectors); err != nil {
		return fmt.Errorf("failed to write vector count: %w", err)
	}

	// Write each vector
	for id, vector := range idx.Vectors {
		// Write ID length and ID
		idBytes := []byte(id)
		idLen := uint16(len(idBytes))
		if err := binary.Write(file, binary.LittleEndian, idLen); err != nil {
			return fmt.Errorf("failed to write ID length: %w", err)
		}
		if _, err := file.Write(idBytes); err != nil {
			return fmt.Errorf("failed to write ID: %w", err)
		}

		// Write vector length and vector
		vecLen := uint32(len(vector))
		if err := binary.Write(file, binary.LittleEndian, vecLen); err != nil {
			return fmt.Errorf("failed to write vector length: %w", err)
		}
		if err := binary.Write(file, binary.LittleEndian, vector); err != nil {
			return fmt.Errorf("failed to write vector: %w", err)
		}

		// Write metadata as JSON
		metadata := idx.Metadata[id]
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataLen := uint32(len(metadataBytes))
		if err := binary.Write(file, binary.LittleEndian, metadataLen); err != nil {
			return fmt.Errorf("failed to write metadata length: %w", err)
		}
		if _, err := file.Write(metadataBytes); err != nil {
			return fmt.Errorf("failed to write metadata: %w", err)
		}
	}

	// Close the file
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close temporary index file: %w", err)
	}

	// Rename the temporary file to the actual file
	if err := os.Rename(tempFile, idx.Path); err != nil {
		return fmt.Errorf("failed to rename temporary index file: %w", err)
	}

	return nil
}

// Load loads the index from disk
func (idx *VectorIndex) Load() error {
	// Open the file
	file, err := os.Open(idx.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// Index doesn't exist yet, that's fine
			return nil
		}
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	// Read and verify header
	header := make([]byte, 24)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	if string(header) != "SENTINEL_VECTOR_INDEX_V1\n" {
		return fmt.Errorf("invalid index file format")
	}

	// Read number of vectors
	var numVectors uint32
	if err := binary.Read(file, binary.LittleEndian, &numVectors); err != nil {
		return fmt.Errorf("failed to read vector count: %w", err)
	}

	// Clear existing data
	idx.Vectors = make(map[string]Embedding)
	idx.Metadata = make(map[string]VectorMetadata)

	// Read each vector
	for i := uint32(0); i < numVectors; i++ {
		// Read ID length and ID
		var idLen uint16
		if err := binary.Read(file, binary.LittleEndian, &idLen); err != nil {
			return fmt.Errorf("failed to read ID length: %w", err)
		}
		idBytes := make([]byte, idLen)
		if _, err := file.Read(idBytes); err != nil {
			return fmt.Errorf("failed to read ID: %w", err)
		}
		id := string(idBytes)

		// Read vector length and vector
		var vecLen uint32
		if err := binary.Read(file, binary.LittleEndian, &vecLen); err != nil {
			return fmt.Errorf("failed to read vector length: %w", err)
		}
		vector := make(Embedding, vecLen)
		if err := binary.Read(file, binary.LittleEndian, vector); err != nil {
			return fmt.Errorf("failed to read vector: %w", err)
		}

		// Read metadata
		var metadataLen uint32
		if err := binary.Read(file, binary.LittleEndian, &metadataLen); err != nil {
			return fmt.Errorf("failed to read metadata length: %w", err)
		}
		metadataBytes := make([]byte, metadataLen)
		if _, err := file.Read(metadataBytes); err != nil {
			return fmt.Errorf("failed to read metadata: %w", err)
		}
		var metadata VectorMetadata
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		// Store in memory
		idx.Vectors[id] = vector
		idx.Metadata[id] = metadata
	}

	return nil
}

// Normalize normalizes a vector to unit length
func Normalize(v Embedding) Embedding {
	norm := float32(0)
	for _, x := range v {
		norm += x * x
	}
	norm = float32(math.Sqrt(float64(norm)))

	if norm == 0 {
		return v
	}

	normalized := make(Embedding, len(v))
	for i, x := range v {
		normalized[i] = x / norm
	}
	return normalized
}

// Similarity functions

// cosineSimilarity computes the cosine similarity between two vectors
func cosineSimilarity(a, b Embedding) float32 {
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

// dotProduct computes the dot product between two vectors
func dotProduct(a, b Embedding) float32 {
	if len(a) != len(b) {
		return 0
	}

	var sum float32
	for i := 0; i < len(a); i++ {
		sum += a[i] * b[i]
	}

	return sum
}

// euclideanDistance computes the Euclidean distance between two vectors
// Returns a similarity score (higher means more similar)
func euclideanDistance(a, b Embedding) float32 {
	if len(a) != len(b) {
		return 0
	}

	var sum float32
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	// Convert distance to similarity (higher is more similar)
	// Using 1/(1+distance) to get a value between 0 and 1
	return 1.0 / (1.0 + float32(math.Sqrt(float64(sum))))
}
