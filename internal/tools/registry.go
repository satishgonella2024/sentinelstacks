package tools

import (
	"fmt"
	"sync"
)

// Registry is a registry of tools
type Registry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// globalRegistry is the singleton registry instance
var globalRegistry *Registry
var once sync.Once

// GetRegistry returns the global registry instance
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			tools: make(map[string]Tool),
		}
		
		// Register default tools will happen here
		// This will be expanded as we implement specific tools
	})
	
	return globalRegistry
}

// RegisterTool registers a tool with the registry
func (r *Registry) RegisterTool(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	name := tool.GetName()
	
	// Check if tool already registered
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}
	
	// Register tool
	r.tools[name] = tool
	
	return nil
}

// GetTool returns a tool by name
func (r *Registry) GetTool(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}
	
	return tool, nil
}

// ListTools returns all registered tools
func (r *Registry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	
	return tools
}

// ListToolNames returns the names of all registered tools
func (r *Registry) ListToolNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	
	return names
}

// ListToolsWithPermission returns tools that require the given permission
func (r *Registry) ListToolsWithPermission(perm Permission) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tools := make([]Tool, 0)
	for _, tool := range r.tools {
		if tool.RequiredPermission() == perm || 
		   tool.RequiredPermission() == PermissionNone ||
		   perm == PermissionAll {
			tools = append(tools, tool)
		}
	}
	
	return tools
}

// RegisterTool registers a tool with the global registry
func RegisterTool(tool Tool) error {
	return GetRegistry().RegisterTool(tool)
}