// Package shim provides interfaces and implementations for LLM providers
package shim

import (
	"context"
	"fmt"
	"sync"

	"github.com/sentinelstacks/sentinel/internal/shim/claude"
	"github.com/sentinelstacks/sentinel/internal/shim/openai"
)

// ProviderRegistry manages the available LLM providers
type ProviderRegistry struct {
	providers map[string]func() Provider
	mu        sync.RWMutex
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]func() Provider),
	}
}

// Register registers a provider with the registry
func (r *ProviderRegistry) Register(name string, factory func() Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = factory
}

// Get returns a provider with the given name
func (r *ProviderRegistry) Get(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.providers[name]
	if !exists {
		return nil, false
	}
	return factory(), true
}

// List returns a list of all registered provider names
func (r *ProviderRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var names []string
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// ProviderInfo contains information about a provider
type ProviderInfo struct {
	Name            string   // Provider name
	AvailableModels []string // Available models
	SupportsChat    bool     // Whether the provider supports chat
	SupportsImages  bool     // Whether the provider supports images
	SupportsAudio   bool     // Whether the provider supports audio
	SupportsVideo   bool     // Whether the provider supports video
}

// GetInfo returns information about a provider
func (r *ProviderRegistry) GetInfo(name string) (*ProviderInfo, error) {
	provider, exists := r.Get(name)
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	info := &ProviderInfo{
		Name:            provider.Name(),
		AvailableModels: provider.AvailableModels(),
		SupportsChat:    true, // All providers support basic chat
		SupportsImages:  provider.SupportsMultimodal(),
		SupportsAudio:   false, // No providers support audio yet
		SupportsVideo:   false, // No providers support video yet
	}

	return info, nil
}

// Global provider registry instance
var globalRegistry = NewProviderRegistry()

// RegisterProviderFactory registers a provider with the global registry
func RegisterProviderFactory(name string, factory func() Provider) {
	globalRegistry.Register(name, factory)
}

// GetProviderFromRegistry returns a provider with the given name from the global registry
func GetProviderFromRegistry(name string) (Provider, bool) {
	return globalRegistry.Get(name)
}

// ListProviderNames returns a list of all registered provider names
func ListProviderNames() []string {
	return globalRegistry.List()
}

// GetProviderInfo returns information about a provider
func GetProviderInfo(name string) (*ProviderInfo, error) {
	return globalRegistry.GetInfo(name)
}

// CreateShim creates a new shim instance for the given provider and model
func CreateShim(providerName, model string, config Config) (Shim, error) {
	provider, exists := GetProviderFromRegistry(providerName)
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", providerName)
	}

	// Create a base shim with the provider
	baseShim := &BaseShim{
		Config:      config,
		Provider:    provider,
		ActiveModel: model,
	}

	// Initialize the shim
	if err := baseShim.Initialize(config); err != nil {
		return nil, fmt.Errorf("failed to initialize shim: %w", err)
	}

	return baseShimToFullShim(baseShim), nil
}

// baseShimToFullShim adds the missing methods to make BaseShim implement the full Shim interface
func baseShimToFullShim(base *BaseShim) Shim {
	return &fullShim{
		BaseShim: base,
	}
}

// fullShim wraps BaseShim and implements the full Shim interface
type fullShim struct {
	*BaseShim
}

// Generate implements the Generate method of the Shim interface
func (s *fullShim) Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error) {
	params := make(map[string]interface{})
	params["model"] = s.ActiveModel
	params["max_tokens"] = input.MaxTokens
	params["temperature"] = input.Temperature

	resp, err := s.Provider.GenerateResponse(ctx, input.Prompt, params)
	if err != nil {
		return nil, err
	}

	return &GenerateOutput{
		Text:       resp,
		FromCache:  false,
		UsedTokens: 0, // Unknown without provider info
		ToolCalls:  []ToolCall{},
	}, nil
}

// Stream implements the Stream method of the Shim interface
func (s *fullShim) Stream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, error) {
	params := make(map[string]interface{})
	params["model"] = s.ActiveModel
	params["max_tokens"] = input.MaxTokens
	params["temperature"] = input.Temperature

	respCh, err := s.Provider.StreamResponse(ctx, input.Prompt, params)
	if err != nil {
		return nil, err
	}

	outCh := make(chan StreamChunk)
	go func() {
		defer close(outCh)
		for text := range respCh {
			outCh <- StreamChunk{
				Text:    text,
				IsFinal: false,
				Error:   nil,
			}
		}
		// Send final chunk
		outCh <- StreamChunk{
			Text:    "",
			IsFinal: true,
			Error:   nil,
		}
	}()

	return outCh, nil
}

// GetEmbeddings implements the GetEmbeddings method of the Shim interface
func (s *fullShim) GetEmbeddings(ctx context.Context, input EmbeddingsInput) (*EmbeddingsOutput, error) {
	embeddings, err := s.Provider.GetEmbeddings(ctx, input.Texts)
	if err != nil {
		return nil, err
	}

	return &EmbeddingsOutput{
		Embeddings: embeddings,
	}, nil
}

// Check if the provider supports the given model
func SupportsModel(providerName, modelName string) (bool, error) {
	provider, exists := GetProviderFromRegistry(providerName)
	if !exists {
		return false, fmt.Errorf("provider not found: %s", providerName)
	}

	availableModels := provider.AvailableModels()
	for _, model := range availableModels {
		if model == modelName {
			return true, nil
		}
	}

	return false, nil
}

// Initialize registry with default providers
func init() {
	// Import and register providers when they're available
	RegisterProviderFactory("claude", func() Provider {
		provider := claude.NewProvider()
		return provider.(Provider)
	})
	RegisterProviderFactory("openai", func() Provider {
		provider := openai.NewProvider()
		return provider.(Provider)
	})
	// RegisterProviderFactory("ollama", ollama.NewProvider)
}
