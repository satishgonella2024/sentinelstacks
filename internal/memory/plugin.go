// Package memory provides memory implementations
package memory

import (
	"fmt"
	"plugin"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// PluginMemoryStoreFactory is a factory that loads memory stores from plugins
type PluginMemoryStoreFactory struct {
	baseFactory types.MemoryStoreFactory
	plugins     map[string]*plugin.Plugin
}

// NewPluginMemoryStoreFactory creates a new plugin memory store factory
func NewPluginMemoryStoreFactory(baseFactory types.MemoryStoreFactory) *PluginMemoryStoreFactory {
	return &PluginMemoryStoreFactory{
		baseFactory: baseFactory,
		plugins:     make(map[string]*plugin.Plugin),
	}
}

// Create creates a new memory store from a plugin
func (f *PluginMemoryStoreFactory) Create(storeType types.MemoryStoreType, config types.MemoryConfig) (types.MemoryStore, error) {
	// First try to create from base factory
	store, err := f.baseFactory.Create(storeType, config)
	if err == nil {
		return store, nil
	}

	// Try to load from plugin
	pluginPath := fmt.Sprintf("memory_%s.so", storeType)
	p, ok := f.plugins[string(storeType)]
	if !ok {
		var err error
		p, err = plugin.Open(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load memory store plugin: %w", err)
		}
		f.plugins[string(storeType)] = p
	}

	// Get create function
	createSym, err := p.Lookup("CreateMemoryStore")
	if err != nil {
		return nil, fmt.Errorf("failed to find CreateMemoryStore in plugin: %w", err)
	}

	// Convert to function
	createFunc, ok := createSym.(func(types.MemoryConfig) (types.MemoryStore, error))
	if !ok {
		return nil, fmt.Errorf("invalid CreateMemoryStore function signature in plugin")
	}

	// Create store
	return createFunc(config)
}

// CreateVector creates a new vector store from a plugin
func (f *PluginMemoryStoreFactory) CreateVector(config types.MemoryConfig) (types.VectorStore, error) {
	// First try to create from base factory
	store, err := f.baseFactory.CreateVector(config)
	if err == nil {
		return store, nil
	}

	// Try to load from plugin
	pluginPath := fmt.Sprintf("vector_%s.so", config.CollectionName)
	p, ok := f.plugins[config.CollectionName]
	if !ok {
		var err error
		p, err = plugin.Open(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load vector store plugin: %w", err)
		}
		f.plugins[config.CollectionName] = p
	}

	// Get create function
	createSym, err := p.Lookup("CreateVectorStore")
	if err != nil {
		return nil, fmt.Errorf("failed to find CreateVectorStore in plugin: %w", err)
	}

	// Convert to function
	createFunc, ok := createSym.(func(types.MemoryConfig) (types.VectorStore, error))
	if !ok {
		return nil, fmt.Errorf("invalid CreateVectorStore function signature in plugin")
	}

	// Create store
	return createFunc(config)
}
