package shim

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

// Config represents configuration for an LLM provider
type Config struct {
	Provider string
	Model    string
	APIKey   string
	Endpoint string
	Timeout  time.Duration
}

// LLMShim is an interface for interacting with different LLM providers
type LLMShim interface {
	// Text completion methods
	Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error)
	CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error)
	
	// Multimodal methods
	MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error)
	MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error)
	
	// Streaming methods
	StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error)
	StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error)
	
	// System prompts
	SetSystemPrompt(prompt string)
	
	// Utility methods
	ParseSentinelfile(content string) (map[string]interface{}, error)
	SupportsMultimodal() bool
	Close() error
}

// ShimFactory creates a new LLM shim based on the provider
func ShimFactory(provider, endpoint, apiKey, model string) (LLMShim, error) {
	// Validate provider
	if provider == "" {
		return nil, fmt.Errorf("provider cannot be empty")
	}
	
	// Create config
	config := Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
		Endpoint: endpoint,
		Timeout:  60 * time.Second, // Default timeout
	}
	
	// Create shim based on provider
	switch provider {
	case "claude":
		return NewClaudeShim(config), nil
	case "openai":
		return NewOpenAIShim(config), nil
	case "ollama":
		return NewOllamaShim(config), nil
	case "google":
		return NewGoogleShim(config), nil
	case "mock":
		return NewMockShim(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// MockShim is a placeholder implementation that mocks LLM responses
type MockShim struct {
	provider     string
	model        string
	systemPrompt string
}

// NewMockShim creates a new mock shim
func NewMockShim(config Config) *MockShim {
	return &MockShim{
		provider: config.Provider,
		model:    config.Model,
	}
}

// Completion mocks an LLM completion
func (s *MockShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	return fmt.Sprintf("This is a mock response from %s model %s.\nSystem: %s\nYou asked: %s", 
		s.provider, s.model, s.systemPrompt, prompt), nil
}

// CompletionWithContext mocks an LLM completion with context
func (s *MockShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	return s.Completion(prompt, maxTokens, temperature, 30*time.Second)
}

// MultimodalCompletion mocks a multimodal completion
func (s *MockShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	// Create a mock output
	output := multimodal.NewOutput()
	
	// Check if there are any images in the input
	hasImages := false
	for _, content := range input.Contents {
		if content.Type == multimodal.MediaTypeImage {
			hasImages = true
			break
		}
	}
	
	if hasImages {
		output.AddText("This is a mock multimodal response. I can see the image you provided.")
	} else {
		output.AddText("This is a mock text response from the multimodal API.")
	}
	
	return output, nil
}

// MultimodalCompletionWithContext mocks a multimodal completion with context
func (s *MockShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	return s.MultimodalCompletion(input, 30*time.Second)
}

// StreamCompletion mocks a streaming completion
func (s *MockShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	ch := make(chan string)
	
	go func() {
		defer close(ch)
		
		// Simulate chunks of response
		messages := []string{
			"This is a mock streaming response ",
			"from " + s.provider + " model " + s.model + ". ",
			"System: " + s.systemPrompt + "\n",
			"You asked: " + prompt,
		}
		
		for _, message := range messages {
			select {
			case <-ctx.Done():
				return
			case ch <- message:
				time.Sleep(200 * time.Millisecond) // Simulate delay between chunks
			}
		}
	}()
	
	return ch, nil
}

// StreamMultimodalCompletion mocks a streaming multimodal completion
func (s *MockShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	ch := make(chan *multimodal.Chunk)
	
	go func() {
		defer close(ch)
		
		// Simulate chunks of response
		messages := []string{
			"This is a mock streaming multimodal response. ",
			"I can process both text and images. ",
			"Let me analyze what you've provided.",
		}
		
		for i, message := range messages {
			chunk := multimodal.NewChunk(multimodal.NewTextContent(message), i == len(messages)-1)
			
			select {
			case <-ctx.Done():
				return
			case ch <- chunk:
				time.Sleep(300 * time.Millisecond) // Simulate delay between chunks
			}
		}
	}()
	
	return ch, nil
}

// SetSystemPrompt sets the system prompt for the model
func (s *MockShim) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// ParseSentinelfile mocks parsing a Sentinelfile
func (s *MockShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"name":        "MockAgent",
		"description": "This is a mock agent parsed from a Sentinelfile",
		"baseModel":   s.model,
		"capabilities": []string{
			"Mock capability 1",
			"Mock capability 2",
		},
	}, nil
}

// SupportsMultimodal returns whether this shim supports multimodal inputs
func (s *MockShim) SupportsMultimodal() bool {
	// Mock shim always supports multimodal for testing
	return true
}

// Close cleans up any resources used by the shim
func (s *MockShim) Close() error {
	// No resources to clean up for mock
	return nil
}
