package shim

import (
	"context"
	"fmt"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim/claude"
	"github.com/sentinelstacks/sentinel/internal/shim/ollama"
)

// Config contains the configuration for a shim
type Config struct {
	Provider   string
	Model      string
	APIKey     string
	Endpoint   string
	Parameters map[string]interface{}
}

// GenerateInput contains the input for generation
type GenerateInput struct {
	Prompt      string
	MaxTokens   int
	Temperature float64
	Tools       []Tool
	Stream      bool
}

// GenerateOutput contains the output from generation
type GenerateOutput struct {
	Text       string
	FromCache  bool
	UsedTokens int
	ToolCalls  []ToolCall
}

// StreamChunk represents a chunk of a streaming response
type StreamChunk struct {
	Text    string
	IsFinal bool
	Error   error
}

// EmbeddingsInput contains the input for embeddings
type EmbeddingsInput struct {
	Texts []string
}

// EmbeddingsOutput contains the output from embeddings
type EmbeddingsOutput struct {
	Embeddings [][]float32
}

// Tool represents a tool that can be used by the LLM
type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

// ToolCall represents a call to a tool by the LLM
type ToolCall struct {
	Tool       string
	Parameters map[string]interface{}
	Result     interface{}
}

// Shim is the interface implemented by all LLM providers
type Shim interface {
	// Initialize initializes the shim with the given configuration
	Initialize(config Config) error

	// Generate generates a response from the LLM
	Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error)

	// Stream streams a response from the LLM
	Stream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, error)

	// GetEmbeddings gets embeddings for the given texts
	GetEmbeddings(ctx context.Context, input EmbeddingsInput) (*EmbeddingsOutput, error)

	// Close closes any resources used by the shim
	Close() error

	// NEW: Multimodal support methods

	// GenerateMultimodal generates a multimodal response from the LLM
	GenerateMultimodal(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error)

	// StreamMultimodal streams a multimodal response from the LLM
	StreamMultimodal(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error)

	// SupportsMultimodal checks if the shim supports multimodal inputs and outputs
	SupportsMultimodal() bool
}

// Provider represents the specific implementation of an LLM provider
type Provider interface {
	// Name returns the name of the provider
	Name() string

	// AvailableModels returns the available models for this provider
	AvailableModels() []string

	// GenerateResponse generates a response from the LLM
	GenerateResponse(ctx context.Context, prompt string, params map[string]interface{}) (string, error)

	// StreamResponse streams a response from the LLM
	StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error)

	// GetEmbeddings gets embeddings for the given texts
	GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error)

	// NEW: Multimodal support methods

	// GenerateMultimodalResponse generates a multimodal response from the LLM
	GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error)

	// StreamMultimodalResponse streams a multimodal response from the LLM
	StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error)

	// SupportsMultimodal checks if the provider supports multimodal inputs and outputs
	SupportsMultimodal() bool
}

// BaseShim is a base implementation of the Shim interface
type BaseShim struct {
	Config      Config
	Provider    Provider
	ActiveModel string
}

// Initialize initializes the shim with the given configuration
func (s *BaseShim) Initialize(config Config) error {
	s.Config = config
	s.ActiveModel = config.Model
	return nil
}

// Close closes any resources used by the shim
func (s *BaseShim) Close() error {
	return nil
}

// SupportsMultimodal checks if the shim supports multimodal inputs and outputs
func (s *BaseShim) SupportsMultimodal() bool {
	if s.Provider == nil {
		return false
	}
	return s.Provider.SupportsMultimodal()
}

// GenerateMultimodal generates a multimodal response from the LLM
func (s *BaseShim) GenerateMultimodal(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	if !s.SupportsMultimodal() {
		return nil, ErrMultimodalNotSupported
	}

	params := make(map[string]interface{})
	params["model"] = s.ActiveModel
	params["max_tokens"] = input.MaxTokens
	params["temperature"] = input.Temperature

	if input.Metadata != nil {
		for k, v := range input.Metadata {
			params[k] = v
		}
	}

	return s.Provider.GenerateMultimodalResponse(ctx, input, params)
}

// StreamMultimodal streams a multimodal response from the LLM
func (s *BaseShim) StreamMultimodal(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	if !s.SupportsMultimodal() {
		return nil, ErrMultimodalNotSupported
	}

	params := make(map[string]interface{})
	params["model"] = s.ActiveModel
	params["max_tokens"] = input.MaxTokens
	params["temperature"] = input.Temperature

	if input.Metadata != nil {
		for k, v := range input.Metadata {
			params[k] = v
		}
	}

	return s.Provider.StreamMultimodalResponse(ctx, input, params)
}

// Register registers the available providers
var providers = make(map[string]func() Provider)

// RegisterProvider registers a provider
func RegisterProvider(name string, factory func() Provider) {
	providers[name] = factory
}

// GetProvider returns a provider with the given name
func GetProvider(name string) (Provider, bool) {
	factory, exists := providers[name]
	if !exists {
		return nil, false
	}
	return factory(), true
}

// Error types
var (
	ErrProviderNotFound       = &shimError{"provider not found"}
	ErrModelNotFound          = &shimError{"model not found"}
	ErrInvalidInput           = &shimError{"invalid input"}
	ErrMultimodalNotSupported = &shimError{"multimodal not supported by this provider"}
)

// shimError is a general error type for shim operations
type shimError struct {
	msg string
}

func (e *shimError) Error() string {
	return e.msg
}

// Message represents a message in a conversation
type Message struct {
	Role    string // e.g., "system", "user", "assistant"
	Content string
}

// LLMShim defines the interface that all LLM provider implementations must satisfy
type LLMShim interface {
	// CompleteChatPrompt sends a conversation to the LLM and returns the response
	CompleteChatPrompt(messages []Message) (string, error)

	// ParseSentinelfile parses a Sentinelfile into a structured format
	ParseSentinelfile(content string) (map[string]interface{}, error)
}

// ShimFactory creates LLM shim implementations based on the provider name
func ShimFactory(provider, endpoint, apiKey string, model string) (LLMShim, error) {
	switch provider {
	case "ollama":
		ollamaShim, err := NewOllamaShim(endpoint, model, apiKey, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to create Ollama shim: %w", err)
		}
		return ollamaShim, nil
	case "claude":
		claudeShim, err := NewClaudeShim(apiKey, model)
		if err != nil {
			return nil, fmt.Errorf("failed to create Claude shim: %w", err)
		}
		return claudeShim, nil
	case "openai":
		// TODO: Implement OpenAI shim
		return nil, fmt.Errorf("OpenAI shim not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}

// NewOllamaShim creates a new Ollama shim
func NewOllamaShim(endpoint, model, apiKey string, maxTokens int) (LLMShim, error) {
	// Create a new Ollama shim
	shimImpl := ollama.NewOllamaShim(endpoint, model, apiKey, maxTokens)

	// Create an adapter that converts between the interface types and the implementation types
	return &ollamaShimAdapter{impl: shimImpl}, nil
}

// ollamaShimAdapter adapts the Ollama-specific implementation to the general LLMShim interface
type ollamaShimAdapter struct {
	impl *ollama.OllamaShim
}

func (a *ollamaShimAdapter) CompleteChatPrompt(messages []Message) (string, error) {
	// Convert messages to Ollama's format
	ollamaMessages := make([]ollama.ChatMessage, len(messages))
	for i, msg := range messages {
		ollamaMessages[i] = ollama.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return a.impl.CompleteChatPrompt(ollamaMessages)
}

func (a *ollamaShimAdapter) ParseSentinelfile(content string) (map[string]interface{}, error) {
	return a.impl.ParseSentinelfile(content)
}

// NewClaudeShim creates a new Claude shim
func NewClaudeShim(apiKey, model string) (LLMShim, error) {
	// Create a new Claude shim
	shimImpl := claude.NewClaudeShim(apiKey, model)

	// Create an adapter that converts between the interface types and the implementation types
	return &claudeShimAdapter{impl: shimImpl}, nil
}

// claudeShimAdapter adapts the Claude-specific implementation to the general LLMShim interface
type claudeShimAdapter struct {
	impl *claude.ClaudeShim
}

func (a *claudeShimAdapter) CompleteChatPrompt(messages []Message) (string, error) {
	// Convert messages to Claude's format
	claudeMessages := make([]claude.ChatMessage, len(messages))
	for i, msg := range messages {
		claudeMessages[i] = claude.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return a.impl.CompleteChatPrompt(claudeMessages)
}

func (a *claudeShimAdapter) ParseSentinelfile(content string) (map[string]interface{}, error) {
	return a.impl.ParseSentinelfile(content)
}
