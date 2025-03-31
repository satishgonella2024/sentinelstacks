// Package mock provides mock implementations for testing
package mock

import (
	"context"
	"fmt"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

// MockProvider is a mock implementation of the Provider interface
type MockProvider struct {
	supportMultimodal bool
}

// NewMockProvider creates a new mock provider
func NewMockProvider() interface{} {
	return &MockProvider{
		supportMultimodal: true,
	}
}

// Name returns the name of the provider
func (p *MockProvider) Name() string {
	return "mock"
}

// AvailableModels returns the available models
func (p *MockProvider) AvailableModels() []string {
	return []string{"mock-model"}
}

// GenerateResponse generates a response from the mock provider
func (p *MockProvider) GenerateResponse(ctx context.Context, prompt string, params map[string]interface{}) (string, error) {
	return fmt.Sprintf("Mock response to: %s", prompt), nil
}

// StreamResponse streams a response from the mock provider
func (p *MockProvider) StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error) {
	ch := make(chan string)

	go func() {
		defer close(ch)

		// Send a mock response in chunks
		words := []string{"Mock", "response", "to:", prompt}

		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case ch <- word + " ":
				// Continue
			}
		}
	}()

	return ch, nil
}

// GetEmbeddings gets embeddings for the given texts
func (p *MockProvider) GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// Return mock embeddings
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = []float32{0.1, 0.2, 0.3}
	}
	return embeddings, nil
}

// SupportsMultimodal checks if the provider supports multimodal inputs and outputs
func (p *MockProvider) SupportsMultimodal() bool {
	return p.supportMultimodal
}

// GenerateMultimodalResponse generates a multimodal response from the mock provider
func (p *MockProvider) GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error) {
	if !p.supportMultimodal {
		return nil, fmt.Errorf("multimodal not supported")
	}

	// Create a mock output
	output := multimodal.NewOutput()

	// Build a response based on the input
	hasImage := false
	textPrompt := ""

	for _, content := range input.Contents {
		if content.Type == multimodal.MediaTypeText {
			textPrompt += content.Text + " "
		} else if content.Type == multimodal.MediaTypeImage {
			hasImage = true
		}
	}

	response := ""
	if hasImage {
		response = fmt.Sprintf("Mock multimodal response to image and text: %s", textPrompt)
	} else {
		response = fmt.Sprintf("Mock text response to: %s", textPrompt)
	}

	output.AddText(response)
	return output, nil
}

// StreamMultimodalResponse streams a multimodal response from the mock provider
func (p *MockProvider) StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error) {
	if !p.supportMultimodal {
		return nil, fmt.Errorf("multimodal not supported")
	}

	ch := make(chan *multimodal.Chunk)

	go func() {
		defer close(ch)

		// Build a response based on the input
		hasImage := false
		textPrompt := ""

		for _, content := range input.Contents {
			if content.Type == multimodal.MediaTypeText {
				textPrompt += content.Text + " "
			} else if content.Type == multimodal.MediaTypeImage {
				hasImage = true
			}
		}

		// Use a different approach to avoid unused variable
		words := []string{"Mock"}
		if hasImage {
			words = append(words, "multimodal", "response", "to", "image", "and", "text:", textPrompt)
		} else {
			words = append(words, "text", "response", "to:", textPrompt)
		}

		// Split into chunks
		for i, word := range words {
			isLast := i == len(words)-1

			chunk := &multimodal.Chunk{
				Content: multimodal.NewTextContent(word + " "),
				IsFinal: isLast,
			}

			select {
			case <-ctx.Done():
				return
			case ch <- chunk:
				// Continue
			}
		}
	}()

	return ch, nil
}
