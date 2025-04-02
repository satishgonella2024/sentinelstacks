package shim

import (
	"fmt"
	"sync"
)

// Registry is a registry of LLM shim providers
type Registry struct {
	providers map[string]ProviderFactory
	mu        sync.RWMutex
}

// ProviderFactory is a function that creates a new LLM shim
type ProviderFactory func(config Config) (LLMShim, error)

// globalRegistry is the singleton registry instance
var globalRegistry *Registry
var once sync.Once

// GetRegistry returns the global registry instance
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			providers: make(map[string]ProviderFactory),
		}
		
		// Register default providers
		globalRegistry.RegisterProvider("claude", func(config Config) (LLMShim, error) {
			return NewClaudeShim(config), nil
		})
		
		globalRegistry.RegisterProvider("openai", func(config Config) (LLMShim, error) {
			return NewOpenAIShim(config), nil
		})
		
		globalRegistry.RegisterProvider("ollama", func(config Config) (LLMShim, error) {
			return NewOllamaShim(config), nil
		})
		
		globalRegistry.RegisterProvider("google", func(config Config) (LLMShim, error) {
			return NewGoogleShim(config), nil
		})
		
		globalRegistry.RegisterProvider("mock", func(config Config) (LLMShim, error) {
			return NewMockShim(config), nil
		})
	})
	
	return globalRegistry
}

// RegisterProvider registers a new provider factory
func (r *Registry) RegisterProvider(name string, factory ProviderFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.providers[name] = factory
}

// GetProvider returns a provider factory by name
func (r *Registry) GetProvider(name string) (ProviderFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	factory, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	
	return factory, nil
}

// CreateShim creates a new LLM shim for the given provider and config
func CreateShim(provider string, config Config) (LLMShim, error) {
	registry := GetRegistry()
	
	factory, err := registry.GetProvider(provider)
	if err != nil {
		return nil, err
	}
	
	return factory(config)
}

// RegisterProvider registers a provider in the global registry
func RegisterProvider(name string, factory ProviderFactory) {
	GetRegistry().RegisterProvider(name, factory)
}
