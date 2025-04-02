package shim

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/tools"
)

// ToolExecutor interface for executing tools
type ToolExecutor interface {
	// ExecuteTool executes a tool with the given parameters
	ExecuteTool(ctx context.Context, agentID, toolName string, params map[string]interface{}) (interface{}, error)
	
	// GetAvailableTools returns the tools available to an agent
	GetAvailableTools(agentID string) []tools.Schema
}

// DefaultToolExecutor implements the ToolExecutor interface
type DefaultToolExecutor struct {
	handler *tools.ToolHandler
}

// NewDefaultToolExecutor creates a new default tool executor
func NewDefaultToolExecutor() (*DefaultToolExecutor, error) {
	// Create tool handler
	handler, err := tools.NewToolHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to create tool handler: %w", err)
	}
	
	return &DefaultToolExecutor{
		handler: handler,
	}, nil
}

// ExecuteTool executes a tool with the given parameters
func (e *DefaultToolExecutor) ExecuteTool(ctx context.Context, agentID, toolName string, params map[string]interface{}) (interface{}, error) {
	// Create function call
	functionCall := &tools.FunctionCall{
		Name:       toolName,
		Parameters: params,
	}
	
	// Execute function
	result := e.handler.HandleFunctionCall(ctx, agentID, functionCall)
	
	// Check for error
	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}
	
	return result.Result, nil
}

// GetAvailableTools returns the tools available to an agent
func (e *DefaultToolExecutor) GetAvailableTools(agentID string) []tools.Schema {
	return e.handler.GenerateToolSchemas(agentID)
}

// ToolAugmentedInput adds tool descriptions to a multimodal input
type ToolAugmentedInput struct {
	Input               *multimodal.Input
	AgentID             string
	ToolExecutor        ToolExecutor
	LastFunctionResults map[string]interface{}
}

// NewToolAugmentedInput creates a new tool-augmented input
func NewToolAugmentedInput(input *multimodal.Input, agentID string, executor ToolExecutor) *ToolAugmentedInput {
	return &ToolAugmentedInput{
		Input:               input,
		AgentID:             agentID,
		ToolExecutor:        executor,
		LastFunctionResults: make(map[string]interface{}),
	}
}

// AddToolResults adds function results to the input
func (t *ToolAugmentedInput) AddToolResults(name string, result interface{}) {
	t.LastFunctionResults[name] = result
}

// PrepareWithTools augments the input with tool descriptions
func (t *ToolAugmentedInput) PrepareWithTools() (*multimodal.Input, error) {
	// Get available tools
	availableTools := t.ToolExecutor.GetAvailableTools(t.AgentID)
	if len(availableTools) == 0 {
		// No tools available, return original input
		return t.Input, nil
	}
	
	// Create a new input with tools information
	augmentedInput := multimodal.NewInput()
	
	// Copy original contents
	for _, content := range t.Input.Contents {
		augmentedInput.AddContent(content)
	}
	
	// Get system message
	systemPrompt, hasSystem := t.Input.GetMetadata("system")
	if !hasSystem {
		systemPrompt = ""
	}
	
	// Add tools information to system prompt
	toolsInfo := generateToolsDescription(availableTools)
	
	// Create new system prompt with tools
	newSystemPrompt := fmt.Sprintf(`%s

%s

%s`,
		systemPrompt,
		toolsInfo,
		generateFunctionResultsInfo(t.LastFunctionResults),
	)
	
	// Set new system prompt
	augmentedInput.SetMetadata("system", newSystemPrompt)
	
	// Copy other metadata
	for key, value := range t.Input.Metadata {
		if key != "system" {
			augmentedInput.SetMetadata(key, value)
		}
	}
	
	// Set max tokens and temperature
	augmentedInput.MaxTokens = t.Input.MaxTokens
	augmentedInput.Temperature = t.Input.Temperature
	augmentedInput.Stream = t.Input.Stream
	
	return augmentedInput, nil
}

// generateToolsDescription generates a description of available tools
func generateToolsDescription(schemas []tools.Schema) string {
	if len(schemas) == 0 {
		return ""
	}
	
	var builder strings.Builder
	
	builder.WriteString("You have access to the following tools:\n\n")
	
	for _, schema := range schemas {
		builder.WriteString(fmt.Sprintf("- %s: %s\n", schema.Name, schema.Description))
		
		// Get parameters
		if properties, ok := schema.Parameters["properties"].(map[string]interface{}); ok {
			builder.WriteString("  Parameters:\n")
			
			for paramName, paramSchema := range properties {
				paramInfo, ok := paramSchema.(map[string]interface{})
				if !ok {
					continue
				}
				
				// Get parameter type and description
				paramType, _ := paramInfo["type"].(string)
				paramDesc, _ := paramInfo["description"].(string)
				
				builder.WriteString(fmt.Sprintf("  - %s (%s): %s\n", paramName, paramType, paramDesc))
			}
		}
		
		builder.WriteString("\n")
	}
	
	builder.WriteString("To use a tool, respond with a function call in this format:\n")
	builder.WriteString("```\ntool_name(param1=value1, param2=value2, ...)\n```\n\n")
	builder.WriteString("For example:\n")
	builder.WriteString("```\nfile/read(path=\"/tmp/example.txt\")\n```\n\n")
	builder.WriteString("You can also format it as a JSON object if you prefer:\n")
	builder.WriteString("```json\n{\"name\": \"tool_name\", \"parameters\": {\"param1\": \"value1\", \"param2\": \"value2\"}}\n```\n\n")
	
	return builder.String()
}

// generateFunctionResultsInfo generates a description of function results
func generateFunctionResultsInfo(results map[string]interface{}) string {
	if len(results) == 0 {
		return ""
	}
	
	var builder strings.Builder
	
	builder.WriteString("Recent function call results:\n\n")
	
	for name, result := range results {
		// Marshal result to JSON
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			continue
		}
		
		builder.WriteString(fmt.Sprintf("Function %s returned:\n```json\n%s\n```\n\n", name, string(resultJSON)))
	}
	
	return builder.String()
}

// ParseFunctionCallFromLLMResponse parses a function call from an LLM response
func ParseFunctionCallFromLLMResponse(response string) (*tools.FunctionCall, error) {
	// Look for function call format in code blocks
	codeBlockStart := strings.Index(response, "```")
	if codeBlockStart >= 0 {
		// Find the end of the code block
		codeBlockEnd := strings.Index(response[codeBlockStart+3:], "```")
		if codeBlockEnd >= 0 {
			// Extract the code block content
			codeBlockContent := response[codeBlockStart+3 : codeBlockStart+3+codeBlockEnd]
			
			// Remove language specifier if present
			firstNewline := strings.Index(codeBlockContent, "\n")
			if firstNewline >= 0 {
				if !strings.Contains(codeBlockContent[:firstNewline], "(") {
					// Likely a language specifier
					codeBlockContent = codeBlockContent[firstNewline+1:]
				}
			}
			
			// Trim spaces
			codeBlockContent = strings.TrimSpace(codeBlockContent)
			
			// Try to parse function call
			functionCall, err := tools.ParseFunctionCall(codeBlockContent)
			if err == nil {
				return functionCall, nil
			}
		}
	}
	
	// Try to find function call format directly in the text (tool_name(params))
	// This regex would look for patterns like: tool_name(param1=value1, param2=value2)
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Check if line contains a function call
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			// Try to parse function call
			functionCall, err := tools.ParseFunctionCall(line)
			if err == nil {
				return functionCall, nil
			}
		}
		
		// Try to parse JSON object
		if strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}") {
			var functionCall tools.FunctionCall
			if err := json.Unmarshal([]byte(line), &functionCall); err == nil {
				if functionCall.Name != "" {
					return &functionCall, nil
				}
			}
		}
	}
	
	return nil, fmt.Errorf("no function call found in response")
}