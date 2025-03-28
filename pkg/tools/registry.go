package tools

import (
	"fmt"
	"os"
	"sync"
)

// ToolRegistry is a global registry of available tools
type ToolRegistry struct {
	manager   *ToolManager
	mutex     sync.RWMutex
	factories map[string]func() Tool
}

// Global singleton instance
var globalRegistry *ToolRegistry
var once sync.Once

// GetToolRegistry returns the global tool registry
func GetToolRegistry() *ToolRegistry {
	once.Do(func() {
		globalRegistry = &ToolRegistry{
			manager:   NewToolManager(),
			factories: make(map[string]func() Tool),
		}
		
		// Register built-in tools
		globalRegistry.RegisterFactory("calculator", func() Tool {
			return &CalculatorTool{}
		})
		
		globalRegistry.RegisterFactory("weather", func() Tool {
			apiKey := os.Getenv("OPENWEATHER_API_KEY")
			return NewWeatherTool(apiKey)
		})
		
		// Register Terraform tool
		globalRegistry.RegisterFactory("terraform", func() Tool {
			return &TerraformTool{}
		})

		// Register URLFetcher tool if it exists
		globalRegistry.RegisterFactory("urlfetcher", func() Tool {
			return &URLFetcherTool{}
		})
	})
	
	return globalRegistry
}

// RegisterFactory registers a tool factory function
func (r *ToolRegistry) RegisterFactory(id string, factory func() Tool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.factories[id] = factory
}

// GetTool creates and returns a tool instance by ID
func (r *ToolRegistry) GetTool(id string) (Tool, error) {
	r.mutex.RLock()
	factory, exists := r.factories[id]
	r.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no tool factory registered for ID: %s", id)
	}
	
	return factory(), nil
}

// CreateToolManager creates a tool manager with the specified tool IDs
func (r *ToolRegistry) CreateToolManager(toolIDs []string) (*ToolManager, error) {
	manager := NewToolManager()
	
	for _, id := range toolIDs {
		tool, err := r.GetTool(id)
		if err != nil {
			return nil, err
		}
		
		if err := manager.RegisterTool(tool); err != nil {
			return nil, err
		}
	}
	
	return manager, nil
}

// ListAvailableTools returns a list of all registered tool IDs
func (r *ToolRegistry) ListAvailableTools() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	tools := make([]string, 0, len(r.factories))
	for id := range r.factories {
		tools = append(tools, id)
	}
	
	return tools
}

// GetAllToolsInfo returns information about all available tools
func (r *ToolRegistry) GetAllToolsInfo() []map[string]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	result := make([]map[string]string, 0, len(r.factories))
	
	for id := range r.factories {
		tool := r.factories[id]()
		info := map[string]string{
			"id":          tool.ID(),
			"name":        tool.Name(),
			"description": tool.Description(),
			"version":     tool.Version(),
		}
		result = append(result, info)
	}
	
	return result
}
