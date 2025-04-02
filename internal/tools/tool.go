// Package tools provides interfaces and implementations for agent tools
package tools

import (
	"context"
	"fmt"
)

// Tool represents a capability that can be used by an agent
type Tool interface {
	// GetName returns the name of the tool
	GetName() string
	
	// GetDescription returns a description of the tool
	GetDescription() string
	
	// GetParameters returns the parameters required by the tool
	GetParameters() []Parameter
	
	// Execute runs the tool with the given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
	
	// RequiredPermission returns the permission required to use this tool
	RequiredPermission() Permission
}

// Parameter represents a parameter for a tool
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
}

// Schema represents the JSON schema for a tool to be exposed to LLM providers
type Schema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// GenerateSchema generates a JSON schema for the tool that can be used with LLM function calling
func GenerateSchema(tool Tool) Schema {
	// Create schema base
	schema := Schema{
		Name:        tool.GetName(),
		Description: tool.GetDescription(),
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
			"required": []string{},
		},
	}
	
	// Add parameters
	properties := schema.Parameters["properties"].(map[string]interface{})
	required := schema.Parameters["required"].([]string)
	
	for _, param := range tool.GetParameters() {
		// Add parameter definition
		paramSchema := map[string]interface{}{
			"type":        param.Type,
			"description": param.Description,
		}
		
		// Add default value if provided
		if param.Default != nil {
			paramSchema["default"] = param.Default
		}
		
		// Add to properties
		properties[param.Name] = paramSchema
		
		// Add to required list if needed
		if param.Required {
			required = append(required, param.Name)
		}
	}
	
	// Update required list
	schema.Parameters["required"] = required
	
	return schema
}

// ValidateParameters validates that the provided parameters match the tool's requirements
func ValidateParameters(tool Tool, params map[string]interface{}) error {
	// Get the tool's parameter definitions
	paramDefs := tool.GetParameters()
	
	// Check for required parameters
	for _, paramDef := range paramDefs {
		if paramDef.Required {
			if _, exists := params[paramDef.Name]; !exists {
				return fmt.Errorf("missing required parameter: %s", paramDef.Name)
			}
		}
	}
	
	// Check parameter types (basic type checking)
	for _, paramDef := range paramDefs {
		value, exists := params[paramDef.Name]
		if !exists {
			continue
		}
		
		// Type checking (simplified)
		switch paramDef.Type {
		case "string":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("parameter %s must be a string", paramDef.Name)
			}
		case "number", "integer":
			// Allow various number types (float64, int, etc.)
			switch value.(type) {
			case float64, float32, int, int64, int32, int16, int8:
				// Valid
			default:
				return fmt.Errorf("parameter %s must be a number", paramDef.Name)
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("parameter %s must be a boolean", paramDef.Name)
			}
		case "array":
			if _, ok := value.([]interface{}); !ok {
				return fmt.Errorf("parameter %s must be an array", paramDef.Name)
			}
		case "object":
			if _, ok := value.(map[string]interface{}); !ok {
				return fmt.Errorf("parameter %s must be an object", paramDef.Name)
			}
		}
	}
	
	return nil
}

// ExecuteTool validates parameters and executes the tool
func ExecuteTool(ctx context.Context, tool Tool, params map[string]interface{}) (interface{}, error) {
	// Validate parameters
	if err := ValidateParameters(tool, params); err != nil {
		return nil, err
	}
	
	// Execute the tool
	return tool.Execute(ctx, params)
}

// BaseTool provides common functionality for tools
type BaseTool struct {
	Name        string
	Description string
	Parameters  []Parameter
	Permission  Permission
}

// GetName returns the name of the tool
func (t *BaseTool) GetName() string {
	return t.Name
}

// GetDescription returns a description of the tool
func (t *BaseTool) GetDescription() string {
	return t.Description
}

// GetParameters returns the parameters required by the tool
func (t *BaseTool) GetParameters() []Parameter {
	return t.Parameters
}

// RequiredPermission returns the permission required to use this tool
func (t *BaseTool) RequiredPermission() Permission {
	return t.Permission
}