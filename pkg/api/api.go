// Package api provides a unified API for the Sentinel Stacks system
package api

import (
	"fmt"

	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// APIConfig contains configuration for all services
type APIConfig struct {
	// StackConfig is the configuration for the stack service
	StackConfig StackServiceConfig

	// MemoryConfig is the configuration for the memory service
	MemoryConfig MemoryServiceConfig

	// RegistryConfig is the configuration for the registry service
	RegistryConfig RegistryServiceConfig
}

// API implements types.API
type API struct {
	stackService    *StackService
	memoryService   *MemoryService
	registryService *RegistryService
}

// NewAPI creates a new API with all services configured
func NewAPI(config APIConfig) (types.API, error) {
	// Create stack service
	stackService, err := NewStackService(config.StackConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack service: %w", err)
	}

	// Create memory service
	memoryService, err := NewMemoryService(config.MemoryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory service: %w", err)
	}

	// Create registry service
	registryService, err := NewRegistryService(config.RegistryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry service: %w", err)
	}

	return &API{
		stackService:    stackService,
		memoryService:   memoryService,
		registryService: registryService,
	}, nil
}

// Stack returns the stack service
func (api *API) Stack() types.StackService {
	return api.stackService
}

// Memory returns the memory service
func (api *API) Memory() types.MemoryService {
	return api.memoryService
}

// Registry returns the registry service
func (api *API) Registry() types.RegistryService {
	return api.registryService
}

// Close releases any resources held by the API services
func (api *API) Close() error {
	var lastErr error

	// Close memory service resources if it has a Close method
	if closer, ok := interface{}(api.memoryService).(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close memory service: %w", err)
		}
	}

	// Close stack service resources if it has a Close method
	if closer, ok := interface{}(api.stackService).(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close stack service: %w", err)
		}
	}

	// Close registry service resources if it has a Close method
	if closer, ok := interface{}(api.registryService).(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close registry service: %w", err)
		}
	}

	return lastErr
}
