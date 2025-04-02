package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FunctionCall represents a function call from an LLM
type FunctionCall struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// FunctionResult represents the result of a function call
type FunctionResult struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

// ToolHandler handles tool execution for LLM function calls
type ToolHandler struct {
	registry          *Registry
	permissionManager *PermissionManager
}

// NewToolHandler creates a new tool handler
func NewToolHandler() (*ToolHandler, error) {
	// Get registry
	registry := GetRegistry()
	
	// Create permission manager
	permManager, err := NewPermissionManager("")
	if err != nil {
		return nil, fmt.Errorf("failed to create permission manager: %w", err)
	}
	
	return &ToolHandler{
		registry:          registry,
		permissionManager: permManager,
	}, nil
}

// GetAvailableTools returns the tools available to an agent
func (h *ToolHandler) GetAvailableTools(agentID string) []Tool {
	// Get permissions for the agent
	permissions := h.permissionManager.GetPermissions(agentID)
	
	// If no permissions, return no tools
	if len(permissions) == 0 {
		return []Tool{}
	}
	
	// Check for "all" permission
	hasAllPermission := false
	for _, perm := range permissions {
		if perm == PermissionAll {
			hasAllPermission = true
			break
		}
	}
	
	// If agent has "all" permission, return all tools
	if hasAllPermission {
		return h.registry.ListTools()
	}
	
	// Otherwise, filter tools by permission
	var availableTools []Tool
	for _, tool := range h.registry.ListTools() {
		requiredPerm := tool.RequiredPermission()
		
		// Always include tools that require no permission
		if requiredPerm == PermissionNone {
			availableTools = append(availableTools, tool)
			continue
		}
		
		// Check if agent has the required permission
		for _, perm := range permissions {
			if perm == requiredPerm {
				availableTools = append(availableTools, tool)
				break
			}
		}
	}
	
	return availableTools
}

// GenerateToolSchemas generates JSON schemas for all available tools for an agent
func (h *ToolHandler) GenerateToolSchemas(agentID string) []Schema {
	// Get available tools
	tools := h.GetAvailableTools(agentID)
	
	// Generate schemas
	schemas := make([]Schema, 0, len(tools))
	for _, tool := range tools {
		schemas = append(schemas, GenerateSchema(tool))
	}
	
	return schemas
}

// HandleFunctionCall handles a function call from an LLM
func (h *ToolHandler) HandleFunctionCall(ctx context.Context, agentID string, functionCall *FunctionCall) *FunctionResult {
	// Create base result
	result := &FunctionResult{
		Name: functionCall.Name,
	}
	
	// Get tool from registry
	tool, err := h.registry.GetTool(functionCall.Name)
	if err != nil {
		result.Error = fmt.Sprintf("Unknown tool: %s", functionCall.Name)
		return result
	}
	
	// Check permission
	if !h.permissionManager.HasPermission(agentID, tool.RequiredPermission()) &&
	   !h.permissionManager.HasPermission(agentID, PermissionAll) &&
	   tool.RequiredPermission() != PermissionNone {
		result.Error = fmt.Sprintf("Permission denied for tool: %s", functionCall.Name)
		return result
	}
	
	// Create timeout context
	execCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	
	// Execute tool
	toolResult, err := ExecuteTool(execCtx, tool, functionCall.Parameters)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	
	// Set result
	result.Result = toolResult
	
	return result
}

// ParseFunctionCall parses a function call string from an LLM
func ParseFunctionCall(functionCallStr string) (*FunctionCall, error) {
	// Remove any potential prefixes like "function_call:" or similar
	functionCallStr = strings.TrimSpace(functionCallStr)
	
	// Simple parsing for function call format: name(param1=value1, param2=value2)
	if strings.Contains(functionCallStr, "(") && strings.HasSuffix(functionCallStr, ")") {
		parts := strings.SplitN(functionCallStr, "(", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid function call format")
		}
		
		name := strings.TrimSpace(parts[0])
		paramsStr := strings.TrimSuffix(parts[1], ")")
		
		// Parse parameters
		params := make(map[string]interface{})
		if paramsStr != "" {
			paramPairs := strings.Split(paramsStr, ",")
			for _, pair := range paramPairs {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) != 2 {
					continue
				}
				
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				
				// Handle different value types
				if value == "true" {
					params[key] = true
				} else if value == "false" {
					params[key] = false
				} else if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					// String value
					params[key] = value[1 : len(value)-1]
				} else if value == "null" || value == "nil" {
					params[key] = nil
				} else {
					// Try as number
					var num float64
					if _, err := fmt.Sscanf(value, "%f", &num); err == nil {
						// If it's a whole number, convert to int
						if num == float64(int(num)) {
							params[key] = int(num)
						} else {
							params[key] = num
						}
					} else {
						// Default to string
						params[key] = value
					}
				}
			}
		}
		
		return &FunctionCall{
			Name:       name,
			Parameters: params,
		}, nil
	}
	
	// Try JSON format
	var call FunctionCall
	if err := json.Unmarshal([]byte(functionCallStr), &call); err != nil {
		return nil, fmt.Errorf("invalid function call JSON: %w", err)
	}
	
	return &call, nil
}

// FormatFunctionResult formats a function result for the LLM
func FormatFunctionResult(result *FunctionResult) string {
	// If there's an error, return error message
	if result.Error != "" {
		return fmt.Sprintf("Function %s failed: %s", result.Name, result.Error)
	}
	
	// Marshal result to JSON
	resultJSON, err := json.MarshalIndent(result.Result, "", "  ")
	if err != nil {
		return fmt.Sprintf("Function %s succeeded but result could not be formatted", result.Name)
	}
	
	return fmt.Sprintf("Function %s returned:\n```json\n%s\n```", result.Name, string(resultJSON))
}

// GrantPermission grants a permission to an agent
func (h *ToolHandler) GrantPermission(agentID string, permission Permission) error {
	return h.permissionManager.Grant(agentID, permission)
}

// RevokePermission revokes a permission from an agent
func (h *ToolHandler) RevokePermission(agentID string, permission Permission) error {
	return h.permissionManager.Revoke(agentID, permission)
}

// GetPermissions returns the permissions for an agent
func (h *ToolHandler) GetPermissions(agentID string) []Permission {
	return h.permissionManager.GetPermissions(agentID)
}