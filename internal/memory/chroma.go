package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ChromaMemoryStore is a memory store implementation using Chroma vector database
type ChromaMemoryStore struct {
	baseURL    string
	collectionID string
	collectionName string
	namespace  string
	ttl        time.Duration
	httpClient *http.Client
}

// ChromaCollection represents a Chroma collection
type ChromaCollection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ChromaEmbedding represents a vector embedding in Chroma
type ChromaEmbedding struct {
	ID       string                 `json:"id"`
	Embedding []float32              `json:"embedding"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ChromaQueryResult represents a query result from Chroma
type ChromaQueryResult struct {
	IDs       []string                `json:"ids"`
	Embeddings [][]float32             `json:"embeddings"`
	Metadatas  []map[string]interface{} `json:"metadatas"`
	Distances  []float32                `json:"distances"`
}

// NewChromaMemoryStore creates a new Chroma memory store
func NewChromaMemoryStore(config MemoryConfig) (*ChromaMemoryStore, error) {
	// Set default connection string
	baseURL := "http://localhost:8000"
	if config.ConnectionString != "" {
		baseURL = config.ConnectionString
	}
	
	// Ensure base URL ends with /api
	if !strings.HasSuffix(baseURL, "/api") {
		baseURL = strings.TrimSuffix(baseURL, "/") + "/api"
	}
	
	// Set default collection name
	collectionName := "sentinel"
	if config.CollectionName != "" {
		collectionName = config.CollectionName
	}
	
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Create store
	store := &ChromaMemoryStore{
		baseURL:    baseURL,
		collectionName: collectionName,
		namespace:  config.Namespace,
		ttl:        config.TTL,
		httpClient: httpClient,
	}
	
	// Initialize collection
	if err := store.initCollection(); err != nil {
		return nil, fmt.Errorf("failed to initialize collection: %w", err)
	}
	
	return store, nil
}

// initCollection creates or gets the collection
func (s *ChromaMemoryStore) initCollection() error {
	// First check if collection already exists
	collections, err := s.listCollections()
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}
	
	// Check if collection exists
	for _, collection := range collections {
		if collection.Name == s.collectionName {
			s.collectionID = collection.ID
			return nil
		}
	}
	
	// Collection doesn't exist, create it
	metadata := map[string]interface{}{
		"namespace": s.namespace,
		"created_at": time.Now().Format(time.RFC3339),
	}
	
	err = s.createCollection(s.collectionName, metadata)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	
	// Get the collection ID
	collections, err = s.listCollections()
	if err != nil {
		return fmt.Errorf("failed to list collections after creation: %w", err)
	}
	
	for _, collection := range collections {
		if collection.Name == s.collectionName {
			s.collectionID = collection.ID
			return nil
		}
	}
	
	return fmt.Errorf("failed to find created collection")
}

// listCollections lists all collections in Chroma
func (s *ChromaMemoryStore) listCollections() ([]ChromaCollection, error) {
	// Create request
	url := fmt.Sprintf("%s/collections", s.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var result struct {
		Collections []ChromaCollection `json:"collections"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return result.Collections, nil
}

// createCollection creates a new collection in Chroma
func (s *ChromaMemoryStore) createCollection(name string, metadata map[string]interface{}) error {
	// Create request payload
	payload := map[string]interface{}{
		"name": name,
		"metadata": metadata,
	}
	
	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	// Create request
	url := fmt.Sprintf("%s/collections", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// collectionURL returns the URL for the current collection
func (s *ChromaMemoryStore) collectionURL() string {
	return fmt.Sprintf("%s/collections/%s", s.baseURL, s.collectionID)
}

// Save stores a value with the given key
func (s *ChromaMemoryStore) Save(ctx context.Context, key string, value interface{}) error {
	// Serialize value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	// Create metadata
	metadata := map[string]interface{}{
		"value": string(valueBytes),
		"updated_at": time.Now().Format(time.RFC3339),
	}
	
	// Add namespace to metadata if specified
	if s.namespace != "" {
		metadata["namespace"] = s.namespace
	}
	
	// Add TTL to metadata if specified
	if s.ttl > 0 {
		metadata["expires_at"] = time.Now().Add(s.ttl).Format(time.RFC3339)
	}
	
	// Create placeholder embedding (just store metadata)
	embeddings := []ChromaEmbedding{
		{
			ID:        key,
			Embedding: make([]float32, 1), // Placeholder
			Metadata:  metadata,
		},
	}
	
	// Upsert embedding
	url := fmt.Sprintf("%s/upsert", s.collectionURL())
	
	// Create request payload
	payload := map[string]interface{}{
		"ids": []string{key},
		"embeddings": [][]float32{embeddings[0].Embedding},
		"metadatas": []map[string]interface{}{embeddings[0].Metadata},
	}
	
	// Execute request
	err = s.executeRequest("POST", url, payload, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert embedding: %w", err)
	}
	
	return nil
}

// Load retrieves a value by key
func (s *ChromaMemoryStore) Load(ctx context.Context, key string) (interface{}, error) {
	// Get embedding by ID
	url := fmt.Sprintf("%s/get", s.collectionURL())
	
	// Create request payload
	payload := map[string]interface{}{
		"ids": []string{key},
	}
	
	// Execute request
	var result struct {
		IDs      []string                `json:"ids"`
		Metadatas []map[string]interface{} `json:"metadatas"`
	}
	
	err := s.executeRequest("POST", url, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %w", err)
	}
	
	// Check if embedding was found
	if len(result.IDs) == 0 {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	// Get metadata
	metadata := result.Metadatas[0]
	
	// Check if entry has expired
	if expiresAt, ok := metadata["expires_at"].(string); ok && expiresAt != "" {
		expiry, err := time.Parse(time.RFC3339, expiresAt)
		if err == nil && time.Now().After(expiry) {
			// Remove expired entry
			go func() {
				s.Delete(context.Background(), key)
			}()
			return nil, fmt.Errorf("key expired: %s", key)
		}
	}
	
	// Get value from metadata
	valueStr, ok := metadata["value"].(string)
	if !ok {
		return nil, fmt.Errorf("value not found in metadata for key: %s", key)
	}
	
	// Unmarshal value
	var value interface{}
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}
	
	return value, nil
}

// Delete removes a key-value pair
func (s *ChromaMemoryStore) Delete(ctx context.Context, key string) error {
	// Delete embedding by ID
	url := fmt.Sprintf("%s/delete", s.collectionURL())
	
	// Create request payload
	payload := map[string]interface{}{
		"ids": []string{key},
	}
	
	// Execute request
	err := s.executeRequest("POST", url, payload, nil)
	if err != nil {
		return fmt.Errorf("failed to delete embedding: %w", err)
	}
	
	return nil
}

// Clear removes all keys and values
func (s *ChromaMemoryStore) Clear(ctx context.Context) error {
	// If namespace is specified, only clear entries with that namespace
	if s.namespace != "" {
		// Get all embeddings
		url := fmt.Sprintf("%s/get", s.collectionURL())
		
		// Create request payload for getting all
		payload := map[string]interface{}{}
		
		// Execute request
		var result struct {
			IDs      []string                `json:"ids"`
			Metadatas []map[string]interface{} `json:"metadatas"`
		}
		
		err := s.executeRequest("POST", url, payload, &result)
		if err != nil {
			return fmt.Errorf("failed to get embeddings: %w", err)
		}
		
		// Filter IDs by namespace
		var idsToDelete []string
		for i, metadata := range result.Metadatas {
			if namespace, ok := metadata["namespace"].(string); ok && namespace == s.namespace {
				idsToDelete = append(idsToDelete, result.IDs[i])
			}
		}
		
		// Delete filtered IDs
		if len(idsToDelete) > 0 {
			deleteURL := fmt.Sprintf("%s/delete", s.collectionURL())
			deletePayload := map[string]interface{}{
				"ids": idsToDelete,
			}
			
			err := s.executeRequest("POST", deleteURL, deletePayload, nil)
			if err != nil {
				return fmt.Errorf("failed to delete embeddings: %w", err)
			}
		}
	} else {
		// Delete all embeddings
		url := fmt.Sprintf("%s/delete", s.collectionURL())
		
		// Create request payload for deleting all
		payload := map[string]interface{}{
			"where": map[string]interface{}{}, // Empty where clause matches all
		}
		
		// Execute request
		err := s.executeRequest("POST", url, payload, nil)
		if err != nil {
			return fmt.Errorf("failed to delete all embeddings: %w", err)
		}
	}
	
	return nil
}

// Keys returns all keys in the store
func (s *ChromaMemoryStore) Keys(ctx context.Context) ([]string, error) {
	// Get all embeddings
	url := fmt.Sprintf("%s/get", s.collectionURL())
	
	// Create request payload for getting all
	payload := map[string]interface{}{}
	
	// If namespace is specified, add where clause
	if s.namespace != "" {
		payload["where"] = map[string]interface{}{
			"namespace": s.namespace,
		}
	}
	
	// Execute request
	var result struct {
		IDs []string `json:"ids"`
	}
	
	err := s.executeRequest("POST", url, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings: %w", err)
	}
	
	return result.IDs, nil
}

// Close closes the memory store
func (s *ChromaMemoryStore) Close() error {
	// No resources to release
	return nil
}

// executeRequest executes an HTTP request to the Chroma API
func (s *ChromaMemoryStore) executeRequest(method, url string, payload interface{}, result interface{}) error {
	// Marshal payload
	var reqBody *bytes.Buffer
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		reqBody = bytes.NewBuffer(payloadBytes)
	}
	
	// Create request
	var req *http.Request
	var err error
	if reqBody != nil {
		req, err = http.NewRequest(method, url, reqBody)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	// Parse response if result is provided
	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}
	
	return nil
}

// SaveEmbedding stores a vector embedding
func (s *ChromaMemoryStore) SaveEmbedding(ctx context.Context, key string, vector []float32, metadata map[string]interface{}) error {
	// Create copy of metadata
	metadataCopy := make(map[string]interface{})
	if metadata != nil {
		for k, v := range metadata {
			metadataCopy[k] = v
		}
	}
	
	// Add namespace to metadata if specified
	if s.namespace != "" {
		metadataCopy["namespace"] = s.namespace
	}
	
	// Add updated_at to metadata
	metadataCopy["updated_at"] = time.Now().Format(time.RFC3339)
	
	// Add TTL to metadata if specified
	if s.ttl > 0 {
		metadataCopy["expires_at"] = time.Now().Add(s.ttl).Format(time.RFC3339)
	}
	
	// Create request payload
	payload := map[string]interface{}{
		"ids": []string{key},
		"embeddings": [][]float32{vector},
		"metadatas": []map[string]interface{}{metadataCopy},
	}
	
	// Upsert embedding
	url := fmt.Sprintf("%s/upsert", s.collectionURL())
	
	// Execute request
	err := s.executeRequest("POST", url, payload, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert embedding: %w", err)
	}
	
	return nil
}

// Query performs a similarity search on stored embeddings
func (s *ChromaMemoryStore) Query(ctx context.Context, vector []float32, topK int) ([]SimilarityMatch, error) {
	// Create request payload
	payload := map[string]interface{}{
		"query_embeddings": [][]float32{vector},
		"n_results": topK,
	}
	
	// Add namespace filter if specified
	if s.namespace != "" {
		payload["where"] = map[string]interface{}{
			"namespace": s.namespace,
		}
	}
	
	// Query embeddings
	url := fmt.Sprintf("%s/query", s.collectionURL())
	
	// Execute request
	var result struct {
		IDs       [][]string                `json:"ids"`
		Distances  [][]float32                `json:"distances"`
		Metadatas  [][]map[string]interface{} `json:"metadatas"`
	}
	
	err := s.executeRequest("POST", url, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query embeddings: %w", err)
	}
	
	// Check if results were found
	if len(result.IDs) == 0 || len(result.IDs[0]) == 0 {
		return []SimilarityMatch{}, nil
	}
	
	// Convert to SimilarityMatch
	matches := make([]SimilarityMatch, len(result.IDs[0]))
	for i := range result.IDs[0] {
		matches[i] = SimilarityMatch{
			Key:      result.IDs[0][i],
			Score:    1.0 - result.Distances[0][i], // Convert distance to similarity
			Metadata: result.Metadatas[0][i],
		}
	}
	
	return matches, nil
}

// DeleteEmbedding removes an embedding
func (s *ChromaMemoryStore) DeleteEmbedding(ctx context.Context, key string) error {
	// Delete is the same as Delete for regular key-value pairs
	return s.Delete(ctx, key)
}
