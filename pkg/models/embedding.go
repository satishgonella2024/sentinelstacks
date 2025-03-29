package models

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/internal/vector"
)

// EmbeddingProvider represents a provider for creating embeddings
type EmbeddingProvider interface {
	GetEmbedding(text string) (vector.Embedding, error)
	GetDimensions() int
}

// GetEmbeddingProvider returns an embedding provider for the given model
func GetEmbeddingProvider(model string) (EmbeddingProvider, error) {
	switch {
	case strings.HasPrefix(model, "openai:"):
		return NewOpenAIEmbeddingProvider(strings.TrimPrefix(model, "openai:")), nil
	case strings.HasPrefix(model, "ollama:"):
		return NewOllamaEmbeddingProvider(strings.TrimPrefix(model, "ollama:")), nil
	default:
		// Default to OpenAI's text-embedding-3-small if no prefix
		return NewOpenAIEmbeddingProvider("text-embedding-3-small"), nil
	}
}

// OpenAIEmbeddingProvider provides embeddings using OpenAI API
type OpenAIEmbeddingProvider struct {
	Model string
}

// NewOpenAIEmbeddingProvider creates a new OpenAI embedding provider
func NewOpenAIEmbeddingProvider(model string) *OpenAIEmbeddingProvider {
	if model == "" {
		model = "text-embedding-3-small"
	}
	return &OpenAIEmbeddingProvider{Model: model}
}

// GetEmbedding gets an embedding from OpenAI
func (p *OpenAIEmbeddingProvider) GetEmbedding(text string) (vector.Embedding, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := NewOpenAIClient(apiKey)

	resp, err := client.CreateEmbedding(context.Background(), p.Model, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	// Convert to our embedding format
	embedding := make(vector.Embedding, len(resp))
	for i, v := range resp {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

// GetDimensions returns the dimensions of the embedding
func (p *OpenAIEmbeddingProvider) GetDimensions() int {
	switch p.Model {
	case "text-embedding-3-small":
		return 1536
	case "text-embedding-3-large":
		return 3072
	case "text-embedding-ada-002":
		return 1536
	default:
		return 1536 // Default to text-embedding-3-small dimensions
	}
}

// OllamaEmbeddingProvider provides embeddings using Ollama
type OllamaEmbeddingProvider struct {
	Model string
}

// NewOllamaEmbeddingProvider creates a new Ollama embedding provider
func NewOllamaEmbeddingProvider(model string) *OllamaEmbeddingProvider {
	if model == "" {
		model = "llama3"
	}
	return &OllamaEmbeddingProvider{Model: model}
}

// GetEmbedding gets an embedding from Ollama
func (p *OllamaEmbeddingProvider) GetEmbedding(text string) (vector.Embedding, error) {
	client := NewOllamaClient()

	resp, err := client.CreateEmbedding(p.Model, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	// Convert to our embedding format
	embedding := make(vector.Embedding, len(resp))
	for i, v := range resp {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

// GetDimensions returns the dimensions of the embedding
func (p *OllamaEmbeddingProvider) GetDimensions() int {
	// Ollama embedding dimensions vary by model
	switch {
	case strings.Contains(p.Model, "llama3"):
		return 4096
	case strings.Contains(p.Model, "mixtral"):
		return 4096
	default:
		return 4096 // Default assumption
	}
}

// Simple clients for embedding API calls
// In a real implementation, these would be more robust

// OpenAIClient is a simple client for OpenAI API
type OpenAIClient struct {
	APIKey string
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{APIKey: apiKey}
}

// CreateEmbedding creates an embedding using OpenAI API
func (c *OpenAIClient) CreateEmbedding(ctx context.Context, model string, text string) ([]float64, error) {
	// This is a simplified implementation
	// In a real implementation, you would make an HTTP request to the OpenAI API
	
	// For now, we'll return a mock embedding to avoid making actual API calls
	// This should be replaced with a real API call
	mockEmbedding := make([]float64, 1536)
	for i := range mockEmbedding {
		mockEmbedding[i] = float64(i % 100) * 0.01
	}
	
	return mockEmbedding, nil
}

// OllamaClient is a simple client for Ollama API
type OllamaClient struct{}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient() *OllamaClient {
	return &OllamaClient{}
}

// CreateEmbedding creates an embedding using Ollama API
func (c *OllamaClient) CreateEmbedding(model string, text string) ([]float64, error) {
	// This is a simplified implementation
	// In a real implementation, you would make an HTTP request to the Ollama API
	
	// For now, we'll return a mock embedding to avoid making actual API calls
	// This should be replaced with a real API call
	mockEmbedding := make([]float64, 4096)
	for i := range mockEmbedding {
		mockEmbedding[i] = float64(i % 100) * 0.01
	}
	
	return mockEmbedding, nil
}
