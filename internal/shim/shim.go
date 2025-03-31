package shim

import (
	"fmt"

	"github.com/sentinelstacks/sentinel/internal/shim/claude"
	"github.com/sentinelstacks/sentinel/internal/shim/ollama"
)

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
