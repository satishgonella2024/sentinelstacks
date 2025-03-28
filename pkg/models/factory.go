package models

import (
	"fmt"
	"os"
	"strings"
)

// ModelAdapterFactory creates model adapters based on configuration
type ModelAdapterFactory struct {
	DefaultOllamaEndpoint string
}

// NewModelAdapterFactory creates a new factory
func NewModelAdapterFactory() *ModelAdapterFactory {
	return &ModelAdapterFactory{
		DefaultOllamaEndpoint: "http://localhost:11434",
	}
}

// CreateAdapter creates a model adapter based on provider and model name
func (f *ModelAdapterFactory) CreateAdapter(provider string, model string, options map[string]interface{}) (ModelAdapter, error) {
	switch strings.ToLower(provider) {
	case "ollama":
		endpoint := f.DefaultOllamaEndpoint

		// First check environment variable
		if envEndpoint := os.Getenv("OLLAMA_ENDPOINT"); envEndpoint != "" {
			endpoint = envEndpoint
		}

		// Then check options (this takes precedence)
		if endpointVal, ok := options["endpoint"].(string); ok && endpointVal != "" {
			endpoint = endpointVal
		}

		return NewOllamaAdapter(endpoint, model), nil

	case "openai":
		// Get API key from environment
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
		}

		return NewOpenAIAdapter(apiKey, model), nil

	case "claude":
		// Get API key from environment
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
		}

		// Default to Claude 3 Sonnet if model not specified
		if model == "" {
			model = "claude-3-sonnet-20240229"
		}

		return NewClaudeAdapter(apiKey, model), nil

	default:
		return nil, fmt.Errorf("unsupported model provider: %s", provider)
	}
}
