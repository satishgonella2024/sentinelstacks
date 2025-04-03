// Package types contains common interfaces and types that are used across multiple packages
// to help prevent import cycles.
package types

import (
	"context"
	"time"
)

// LLMProvider defines the interface for LLM providers
type LLMProvider interface {
	// Complete generates a completion given a prompt
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
	
	// GetName returns the name of the provider
	GetName() string
	
	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string
}

// CompletionRequest represents a request to an LLM provider
type CompletionRequest struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
}

// CompletionResponse represents a response from an LLM provider
type CompletionResponse struct {
	Text         string
	FinishReason string
	Usage        struct {
		PromptTokens     int
		CompletionTokens int
		TotalTokens      int
	}
}

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeText  MediaType = "text"
	MediaTypeImage MediaType = "image"
)

// Content represents multimodal content
type Content struct {
	Type MediaType
	Text string
	URL  string
	Data []byte
}

// Input represents an input to a multimodal model
type Input struct {
	Contents   []Content
	MaxTokens  int
	Temperature float64
}

// Output represents the output from a multimodal model
type Output struct {
	Contents []Content
}

// Chunk represents a chunk of streaming output
type Chunk struct {
	Content *Content
	IsFinal bool
}

// NewContent creates a new text content
func NewTextContent(text string) *Content {
	return &Content{
		Type: MediaTypeText,
		Text: text,
	}
}

// NewInput creates a new input
func NewInput() *Input {
	return &Input{
		Contents: []Content{},
	}
}

// NewOutput creates a new output
func NewOutput() *Output {
	return &Output{
		Contents: []Content{},
	}
}

// NewChunk creates a new chunk
func NewChunk(content *Content, isFinal bool) *Chunk {
	return &Chunk{
		Content: content,
		IsFinal: isFinal,
	}
}

// LLMShimConfig represents configuration for an LLM shim
type LLMShimConfig struct {
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
	MultimodalCompletion(input *Input, timeout time.Duration) (*Output, error)
	MultimodalCompletionWithContext(ctx context.Context, input *Input) (*Output, error)
	
	// Streaming methods
	StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error)
	StreamMultimodalCompletion(ctx context.Context, input *Input) (<-chan *Chunk, error)
	
	// System prompts
	SetSystemPrompt(prompt string)
	
	// Utility methods
	ParseSentinelfile(content string) (map[string]interface{}, error)
	SupportsMultimodal() bool
	Close() error
}
