// Package memory provides memory store implementations
package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// ChromaConfig contains configuration for a Chroma vector store
type ChromaConfig struct {
	// Server is the URL of the Chroma server
	Server string

	// CollectionName is the name of the collection
	CollectionName string

	// Dimensions is the dimensionality of embeddings
	Dimensions int

	// APIKey is the API key for the Chroma server
	APIKey string
}

// ChromaStore is a vector store implementation using the Chroma DB API
type ChromaStore struct {
	config ChromaConfig
	client *http.Client
}

// NewChromaStore creates a new Chroma vector store
func NewChromaStore(config ChromaConfig) (*ChromaStore, error) {
	if config.Server == "" {
		config.Server = "http://localhost:8000"
	}
	if config.CollectionName == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	if config.Dimensions <= 0 {
		config.Dimensions = 1536 // Default for OpenAI embeddings
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Check if collection exists, create if not
	store := &ChromaStore{
		config: config,
		client: client,
	}

	// Initialize the collection
	if err := store.initCollection(); err != nil {
		return nil, fmt.Errorf("failed to initialize collection: %w", err)
	}

	return store, nil
}

// initCollection creates the collection if it doesn't exist
func (s *ChromaStore) initCollection() error {
	// Check if collection exists
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/collections", s.config.Server), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check collections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to check collections: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var collections struct {
		Collections []struct {
			Name string `json:"name"`
		} `json:"collections"`
	}
	if err := json.Unmarshal(body, &collections); err != nil {
		return fmt.Errorf("failed to parse collections: %w", err)
	}

	// Check if our collection exists
	collectionExists := false
	for _, coll := range collections.Collections {
		if coll.Name == s.config.CollectionName {
			collectionExists = true
			break
		}
	}

	// Create collection if it doesn't exist
	if !collectionExists {
		createReq := struct {
			Name     string   `json:"name"`
			Metadata struct{} `json:"metadata"`
		}{
			Name:     s.config.CollectionName,
			Metadata: struct{}{},
		}

		payload, err := json.Marshal(createReq)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections", s.config.Server), strings.NewReader(string(payload)))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		s.setHeaders(req)

		resp, err := s.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to create collection: %s", resp.Status)
		}
	}

	return nil
}

// setHeaders sets common headers for Chroma API requests
func (s *ChromaStore) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if s.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	}
}

// StoreVector stores a text with its vector embedding
func (s *ChromaStore) StoreVector(ctx context.Context, id string, text string, metadata map[string]interface{}) error {
	// In a real implementation, we would generate embeddings here
	// For now, we'll just create a mock vector
	vector := make([]float32, s.config.Dimensions)
	for i := 0; i < s.config.Dimensions && i < len(text); i++ {
		if i < len(text) {
			vector[i] = float32(text[i]) / 255.0
		}
	}

	// Prepare the request
	reqBody := struct {
		IDs        []string                 `json:"ids"`
		Embeddings [][]float32              `json:"embeddings"`
		Documents  []string                 `json:"documents"`
		Metadatas  []map[string]interface{} `json:"metadatas"`
	}{
		IDs:        []string{id},
		Embeddings: [][]float32{vector},
		Documents:  []string{text},
		Metadatas:  []map[string]interface{}{metadata},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections/%s/upsert", s.config.Server, s.config.CollectionName), strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upsert vector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upsert vector: %s - %s", resp.Status, string(body))
	}

	return nil
}

// SearchVector finds similar vectors using vector similarity
func (s *ChromaStore) SearchVector(ctx context.Context, text string, limit int, filter map[string]interface{}) ([]types.MemoryMatch, error) {
	if limit <= 0 {
		limit = 10
	}

	// In a real implementation, we would generate embeddings here
	// For now, we'll just create a mock vector
	queryVector := make([]float32, s.config.Dimensions)
	for i := 0; i < s.config.Dimensions && i < len(text); i++ {
		if i < len(text) {
			queryVector[i] = float32(text[i]) / 255.0
		}
	}

	// Prepare the request
	reqBody := struct {
		QueryEmbeddings [][]float32            `json:"query_embeddings"`
		NResults        int                    `json:"n_results"`
		Where           map[string]interface{} `json:"where"`
		Include         []string               `json:"include"`
	}{
		QueryEmbeddings: [][]float32{queryVector},
		NResults:        limit,
		Where:           filter,
		Include:         []string{"documents", "metadatas", "distances"},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections/%s/query", s.config.Server, s.config.CollectionName), strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query vectors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to query vectors: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		IDs       []string                 `json:"ids"`
		Documents []string                 `json:"documents"`
		Metadatas []map[string]interface{} `json:"metadatas"`
		Distances []float64                `json:"distances"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to MemoryMatch objects
	matches := make([]types.MemoryMatch, 0, len(result.IDs))
	for i := 0; i < len(result.IDs); i++ {
		// Convert similarity to score (higher is better)
		// Chroma distances are typically L2 distances, so we need to convert
		// Lower distance means more similar, so we invert it
		distance := result.Distances[i]
		score := 1.0 / (1.0 + distance)

		matches = append(matches, types.MemoryMatch{
			Key:       result.IDs[i],
			Content:   result.Documents[i],
			Metadata:  result.Metadatas[i],
			Score:     score,
			Distance:  distance,
			Timestamp: time.Now(), // We don't get timestamp from Chroma
		})
	}

	return matches, nil
}

// GetVector retrieves a vector by ID
func (s *ChromaStore) GetVector(ctx context.Context, id string) (*types.MemoryMatch, error) {
	// Prepare the request
	reqBody := struct {
		IDs     []string `json:"ids"`
		Include []string `json:"include"`
	}{
		IDs:     []string{id},
		Include: []string{"documents", "metadatas"},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections/%s/get", s.config.Server, s.config.CollectionName), strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get vector: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		IDs       []string                 `json:"ids"`
		Documents []string                 `json:"documents"`
		Metadatas []map[string]interface{} `json:"metadatas"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the vector was found
	if len(result.IDs) == 0 {
		return nil, fmt.Errorf("vector not found: %s", id)
	}

	// Return the first match
	return &types.MemoryMatch{
		Key:       result.IDs[0],
		Content:   result.Documents[0],
		Metadata:  result.Metadatas[0],
		Score:     1.0, // Perfect match for direct retrieval
		Distance:  0.0,
		Timestamp: time.Now(), // We don't get timestamp from Chroma
	}, nil
}

// DeleteVector removes a vector by ID
func (s *ChromaStore) DeleteVector(ctx context.Context, id string) error {
	// Prepare the request
	reqBody := struct {
		IDs []string `json:"ids"`
	}{
		IDs: []string{id},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/collections/%s", s.config.Server, s.config.CollectionName), strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete vector: %s - %s", resp.Status, string(body))
	}

	return nil
}

// ListVectors returns all IDs in the store
func (s *ChromaStore) ListVectors(ctx context.Context) ([]string, error) {
	// Prepare the request
	reqBody := struct {
		Include []string `json:"include"`
		Limit   int      `json:"limit"`
	}{
		Include: []string{"ids"},
		Limit:   10000, // A high limit to get all IDs
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/collections/%s/get", s.config.Server, s.config.CollectionName), strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	s.setHeaders(req)
	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list vectors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list vectors: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		IDs []string `json:"ids"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.IDs, nil
}

// Close releases resources
func (s *ChromaStore) Close() error {
	// No need to close anything, as the HTTP client will be garbage collected
	return nil
}
