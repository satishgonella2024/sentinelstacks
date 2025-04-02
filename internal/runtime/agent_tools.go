package runtime

import (
	"context"
	"fmt"
	
	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim"
	"github.com/sentinelstacks/sentinel/internal/tools"
	"github.com/sentinelstacks/sentinel/internal/tools/file"
	"github.com/sentinelstacks/sentinel/internal/tools/web"
)

// init registers all tools
func init() {
	// Register file tools
	if err := file.RegisterFileTools(); err != nil {
		fmt.Printf("Warning: Failed to register file tools: %v\n", err)
	}
	
	// Register web tools
	if err := web.RegisterWebTools(); err != nil {
		fmt.Printf("Warning: Failed to register web tools: %v\n", err)
	}
}

// ToolsCoordinator manages tool execution for an agent
type ToolsCoordinator struct {
	agentID      string
	executor     shim.ToolExecutor
	conversation map[string]interface{}
}

// NewToolsCoordinator creates a new tools coordinator
func NewToolsCoordinator(agentID string) (*ToolsCoordinator, error) {
	// Create tool executor
	executor, err := shim.NewDefaultToolExecutor()
	if err != nil {
		return nil, fmt.Errorf("failed to create tool executor: %w", err)
	}
	
	return &ToolsCoordinator{
		agentID:      agentID,
		executor:     executor,
		conversation: make(map[string]interface{}),
	}, nil
}

// ProcessMultimodalInput processes a multimodal input with tools support
func (c *ToolsCoordinator) ProcessMultimodalInput(ctx context.Context, mmAgent *MultimodalAgent, input *multimodal.Input, maxTurns int) (*multimodal.Output, error) {
	// Create tool-augmented input
	toolInput := shim.NewToolAugmentedInput(input, c.agentID, c.executor)
	
	// Prepare input with tools
	augmentedInput, err := toolInput.PrepareWithTools()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare input with tools: %w", err)
	}
	
	// Process with LLM
	output, err := mmAgent.ProcessMultimodalInput(ctx, augmentedInput)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %w", err)
	}
	
	// Extract text from output
	var responseText string
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			responseText += content.Text
		}
	}
	
	// Check for function calls in response
	functionCall, err := shim.ParseFunctionCallFromLLMResponse(responseText)
	if err != nil || functionCall == nil {
		// No function call, return original output
		return output, nil
	}
	
	// Execute function call
	if maxTurns <= 0 {
		// No more turns, return warning
		warningOutput := multimodal.NewOutput()
		warningOutput.AddText(fmt.Sprintf("Tool call detected (%s), but maximum tool call limit reached. Tool was not executed.", functionCall.Name))
		return warningOutput, nil
	}
	
	// Execute tool
	result, err := c.executor.ExecuteTool(ctx, c.agentID, functionCall.Name, functionCall.Parameters)
	if err != nil {
		// Tool execution failed
		errorOutput := multimodal.NewOutput()
		errorOutput.AddText(fmt.Sprintf("Tool call %s failed: %v", functionCall.Name, err))
		return errorOutput, nil
	}
	
	// Add result to conversation
	toolInput.AddToolResults(functionCall.Name, result)
	
	// Create follow-up input with tool results
	followupInput := multimodal.NewInput()
	followupInput.AddText(fmt.Sprintf("I executed the tool %s with the provided parameters. Here is the result:\n\n%s\n\nPlease analyze this result and respond accordingly.", 
		functionCall.Name,
		tools.FormatFunctionResult(&tools.FunctionResult{Name: functionCall.Name, Result: result})))
	
	// Set metadata from original input
	for key, value := range input.Metadata {
		followupInput.SetMetadata(key, value)
	}
	
	// Create new tool input with results
	newToolInput := shim.NewToolAugmentedInput(followupInput, c.agentID, c.executor)
	for name, res := range toolInput.LastFunctionResults {
		newToolInput.AddToolResults(name, res)
	}
	newToolInput.AddToolResults(functionCall.Name, result)
	
	// Recursively process with new input and reduced maxTurns
	return c.ProcessMultimodalInput(ctx, mmAgent, followupInput, maxTurns-1)
}