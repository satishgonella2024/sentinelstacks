package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/tools"
	"github.com/satishgonella2024/sentinelstacks/pkg/agent"
)

// AddToolsToAgent adds tool capabilities to an agent
func (a *MultimodalAgent) AddToolsToAgent() error {
	// Create tools coordinator
	coordinator, err := NewToolsCoordinator(a.Agent.ID)
	if err != nil {
		return fmt.Errorf("failed to create tools coordinator: %w", err)
	}
	
	// Store coordinator in agent metadata
	a.metadata["tools_coordinator"] = coordinator
	
	return nil
}

// ProcessInputWithTools processes input with tool support
func (a *MultimodalAgent) ProcessInputWithTools(ctx context.Context, input *multimodal.Input, maxToolCalls int) (*multimodal.Output, error) {
	// Get tools coordinator
	coordinator, ok := a.metadata["tools_coordinator"].(*ToolsCoordinator)
	if !ok {
		// No tools coordinator, fallback to normal processing
		return a.ProcessMultimodalInput(ctx, input)
	}
	
	// Process with tools
	return coordinator.ProcessMultimodalInput(ctx, a, input, maxToolCalls)
}

// ProcessTextInputWithTools processes text input with tool support
func (a *MultimodalAgent) ProcessTextInputWithTools(ctx context.Context, text string, maxToolCalls int) (string, error) {
	// Create multimodal input from text
	input := multimodal.NewInput()
	input.AddText(text)
	
	// Process with tools
	output, err := a.ProcessInputWithTools(ctx, input, maxToolCalls)
	if err != nil {
		return "", err
	}
	
	// Extract text from output
	var responseText string
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			responseText += content.Text
		}
	}
	
	return responseText, nil
}

// GrantToolPermission grants a tool permission to the agent
func (a *MultimodalAgent) GrantToolPermission(permission tools.Permission) error {
	// Create tool handler if not exists
	handler, err := tools.NewToolHandler()
	if err != nil {
		return fmt.Errorf("failed to create tool handler: %w", err)
	}
	
	// Grant permission
	return handler.GrantPermission(a.Agent.ID, permission)
}

// RevokeToolPermission revokes a tool permission from the agent
func (a *MultimodalAgent) RevokeToolPermission(permission tools.Permission) error {
	// Create tool handler if not exists
	handler, err := tools.NewToolHandler()
	if err != nil {
		return fmt.Errorf("failed to create tool handler: %w", err)
	}
	
	// Revoke permission
	return handler.RevokePermission(a.Agent.ID, permission)
}

// GetToolPermissions returns the tool permissions for the agent
func (a *MultimodalAgent) GetToolPermissions() ([]tools.Permission, error) {
	// Create tool handler if not exists
	handler, err := tools.NewToolHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to create tool handler: %w", err)
	}
	
	// Get permissions
	return handler.GetPermissions(a.Agent.ID), nil
}

// ConfigureToolsFromDefinition configures tools based on agent definition
func (a *MultimodalAgent) ConfigureToolsFromDefinition(def *agent.Definition) error {
	// Check if agent has tools
	if len(def.Tools) == 0 {
		return nil // No tools to configure
	}
	
	// Add tools to agent
	if err := a.AddToolsToAgent(); err != nil {
		return fmt.Errorf("failed to add tools to agent: %w", err)
	}
	
	// Grant permissions based on tools
	handler, err := tools.NewToolHandler()
	if err != nil {
		return fmt.Errorf("failed to create tool handler: %w", err)
	}
	
	// Map of tool prefixes to permissions
	permissionMap := map[string]tools.Permission{
		"file/": tools.PermissionFile,
		"web/":  tools.PermissionNetwork,
		"shell/": tools.PermissionShell,
		"api/":   tools.PermissionAPI,
	}
	
	// Keep track of permissions that need to be granted
	permissionsToGrant := make(map[tools.Permission]bool)
	
	// Check each tool
	for _, toolName := range def.Tools {
		// Look for matching prefix
		for prefix, permission := range permissionMap {
			if strings.HasPrefix(toolName, prefix) {
				permissionsToGrant[permission] = true
				break
			}
		}
	}
	
	// Grant all needed permissions
	for permission := range permissionsToGrant {
		if err := handler.GrantPermission(a.Agent.ID, permission); err != nil {
			return fmt.Errorf("failed to grant permission %s: %w", permission, err)
		}
	}
	
	return nil
}