package tools

import (
	"encoding/json"
	"fmt"
)

// Tool defines the interface that all agent tools must implement
type Tool interface {
	// ID returns the unique identifier for the tool
	ID() string

	// Name returns a user-friendly name for the tool
	Name() string

	// Description returns a detailed description of what the tool does
	Description() string

	// Version returns the semantic version of the tool
	Version() string

	// ParameterSchema returns the JSON schema for the tool's parameters
	ParameterSchema() map[string]interface{}

	// Execute runs the tool with the provided parameters and returns the result
	Execute(params map[string]interface{}) (interface{}, error)
}

// ToolManager manages the registration and execution of tools
type ToolManager struct {
	tools map[string]Tool
}

// NewToolManager creates a new tool manager
func NewToolManager() *ToolManager {
	return &ToolManager{
		tools: make(map[string]Tool),
	}
}

// RegisterTool adds a tool to the manager
func (tm *ToolManager) RegisterTool(tool Tool) error {
	id := tool.ID()
	if _, exists := tm.tools[id]; exists {
		return fmt.Errorf("tool with ID %s is already registered", id)
	}
	tm.tools[id] = tool
	return nil
}

// GetTool returns a tool by ID
func (tm *ToolManager) GetTool(id string) (Tool, error) {
	tool, exists := tm.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool with ID %s not found", id)
	}
	return tool, nil
}

// ListTools returns all registered tools
func (tm *ToolManager) ListTools() []Tool {
	tools := make([]Tool, 0, len(tm.tools))
	for _, tool := range tm.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ExecuteTool runs a tool by ID with the provided parameters
func (tm *ToolManager) ExecuteTool(id string, params map[string]interface{}) (interface{}, error) {
	tool, err := tm.GetTool(id)
	if err != nil {
		return nil, err
	}
	return tool.Execute(params)
}

// ToolManifest represents tool metadata in a format suitable for AI models
type ToolManifest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// GenerateManifests creates tool manifests for all registered tools
func (tm *ToolManager) GenerateManifests() []ToolManifest {
	manifests := make([]ToolManifest, 0, len(tm.tools))
	for _, tool := range tm.tools {
		manifest := ToolManifest{
			ID:          tool.ID(),
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  tool.ParameterSchema(),
		}
		manifests = append(manifests, manifest)
	}
	return manifests
}

// GenerateManifestsJSON returns tool manifests as a JSON string
func (tm *ToolManager) GenerateManifestsJSON() (string, error) {
	manifests := tm.GenerateManifests()
	data, err := json.MarshalIndent(manifests, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal tool manifests: %w", err)
	}
	return string(data), nil
}
