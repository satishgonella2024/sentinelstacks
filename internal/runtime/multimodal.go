// Package runtime provides the agent runtime implementation
package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentinelstacks/sentinel/internal/conversation"
	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim"
)

// MultimodalAgent extends the Agent struct with multimodal capabilities
type MultimodalAgent struct {
	*Agent
	History         *conversation.History
	LLM             shim.LLMShim
	MaxTokens       int
	Temperature     float64
	ConversationDir string
	metadata        map[string]interface{}
}

// NewMultimodalAgent creates a new multimodal agent
func NewMultimodalAgent(agent *Agent, config shim.Config) (*MultimodalAgent, error) {
	// Create conversation directory
	conversationDir := filepath.Join(agent.StateDir, "conversations")
	if err := os.MkdirAll(conversationDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create conversation directory: %w", err)
	}

	// Create a new conversation history
	history := conversation.NewHistory()
	history.SetID(fmt.Sprintf("session_%d", time.Now().UnixNano()))

	// Create the LLM shim
	llmShim, err := shim.ShimFactory(
		config.Provider,
		config.Endpoint,
		config.APIKey,
		config.Model,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create LLM shim: %w", err)
	}

	// Set the system prompt based on the agent definition
	systemPrompt := generateSystemPrompt(agent)
	llmShim.SetSystemPrompt(systemPrompt)

	// Check if the shim supports multimodal if needed
	if !llmShim.SupportsMultimodal() {
		// If the model doesn't support multimodal, it's not an error,
		// but we should log a warning
		fmt.Printf("Warning: Provider %s with model %s does not support multimodal inputs\n",
			config.Provider, config.Model)
	}

	// Create the multimodal agent
	return &MultimodalAgent{
		Agent:           agent,
		History:         history,
		LLM:             llmShim,
		MaxTokens:       4096, // Default
		Temperature:     0.7,  // Default
		ConversationDir: conversationDir,
		metadata:        make(map[string]interface{}),
	}, nil
}

// ProcessTextInput processes text input from the user
func (ma *MultimodalAgent) ProcessTextInput(ctx context.Context, text string) (string, error) {
	// Add the user message to the history
	ma.History.AddMessage("user", text)

	// Check if the LLM supports multimodal
	var responseText string
	var err error
	
	if ma.LLM.SupportsMultimodal() {
		// Use multimodal API for models that support it
		input := multimodal.NewInput()
		input.AddText(text)
		
		// Set generation parameters
		input.MaxTokens = ma.MaxTokens
		input.Temperature = ma.Temperature

		// Generate a response
		output, err := ma.LLM.MultimodalCompletionWithContext(ctx, input)
		if err != nil {
			return "", fmt.Errorf("failed to generate response: %w", err)
		}

		// Extract text from the output
		for _, content := range output.Contents {
			if content.Type == multimodal.MediaTypeText {
				responseText += content.Text
			}
		}
	} else {
		// Use text-only API for models that don't support multimodal
		responseText, err = ma.LLM.CompletionWithContext(ctx, text, ma.MaxTokens, ma.Temperature)
		if err != nil {
			return "", fmt.Errorf("failed to generate response: %w", err)
		}
	}

	// Add the assistant's response to the history
	ma.History.AddMessage("assistant", responseText)

	// Save the conversation
	if err := ma.saveConversation(); err != nil {
		// Just log the error, don't fail the response
		fmt.Printf("Warning: Failed to save conversation: %v\n", err)
	}

	return responseText, nil
}

// ProcessMultimodalInput processes multimodal input from the user (text + images)
func (ma *MultimodalAgent) ProcessMultimodalInput(ctx context.Context, userInput *multimodal.Input) (*multimodal.Output, error) {
	// Check if the LLM supports multimodal
	if !ma.LLM.SupportsMultimodal() {
		return nil, fmt.Errorf("the LLM does not support multimodal input")
	}

	// Add the multimodal message to history
	// For now, we just add the text part to history, but in a real implementation
	// we'd want to store the full multimodal content
	textContent := extractTextFromInput(userInput)
	ma.History.AddMessage("user", textContent)
	
	// Set generation parameters if not already set
	if userInput.MaxTokens <= 0 {
		userInput.MaxTokens = ma.MaxTokens
	}
	
	if userInput.Temperature <= 0 {
		userInput.Temperature = ma.Temperature
	}

	// Generate a response
	output, err := ma.LLM.MultimodalCompletionWithContext(ctx, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate multimodal response: %w", err)
	}

	// Extract text from the output for history
	responseText := extractTextFromOutput(output)
	
	// Add the assistant's response to the history
	ma.History.AddMessage("assistant", responseText)

	// Save the conversation
	if err := ma.saveConversation(); err != nil {
		// Just log the error, don't fail the response
		fmt.Printf("Warning: Failed to save conversation: %v\n", err)
	}

	return output, nil
}

// StreamResponse streams a response to a text input
func (ma *MultimodalAgent) StreamResponse(ctx context.Context, text string) (<-chan string, error) {
	// Add the user message to the history
	ma.History.AddMessage("user", text)

	// Stream the response
	// Note: we don't check SupportsMultimodal() here because StreamCompletion is text-only
	responseStream, err := ma.LLM.StreamCompletion(ctx, text, ma.MaxTokens, ma.Temperature)
	if err != nil {
		return nil, fmt.Errorf("failed to stream response: %w", err)
	}

	// Create a channel for the processed response
	processedStream := make(chan string)
	
	// Process the response in a goroutine
	go func() {
		defer close(processedStream)
		
		var fullResponse string
		
		// Read from the response stream and forward to the processed stream
		for chunk := range responseStream {
			// Forward the chunk
			select {
			case <-ctx.Done():
				return
			case processedStream <- chunk:
				// Accumulate the full response
				fullResponse += chunk
			}
		}
		
		// Add the full response to the history
		ma.History.AddMessage("assistant", fullResponse)
		
		// Save the conversation
		if err := ma.saveConversation(); err != nil {
			// Just log the error
			fmt.Printf("Warning: Failed to save conversation: %v\n", err)
		}
	}()
	
	return processedStream, nil
}

// StreamMultimodalResponse streams a response to a multimodal input
func (ma *MultimodalAgent) StreamMultimodalResponse(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// Check if the LLM supports multimodal
	if !ma.LLM.SupportsMultimodal() {
		return nil, fmt.Errorf("the LLM does not support multimodal input")
	}

	// Add the user message to the history (text part only for now)
	textContent := extractTextFromInput(input)
	ma.History.AddMessage("user", textContent)

	// Set parameters if not already set
	if input.MaxTokens <= 0 {
		input.MaxTokens = ma.MaxTokens
	}
	
	if input.Temperature <= 0 {
		input.Temperature = ma.Temperature
	}

	// Stream the response
	responseStream, err := ma.LLM.StreamMultimodalCompletion(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stream multimodal response: %w", err)
	}

	// Create a channel for the processed response
	processedStream := make(chan *multimodal.Chunk)
	
	// Process the response in a goroutine
	go func() {
		defer close(processedStream)
		
		var fullResponse string
		
		// Read from the response stream and forward to the processed stream
		for chunk := range responseStream {
			// Forward the chunk
			select {
			case <-ctx.Done():
				return
			case processedStream <- chunk:
				// Accumulate the full response if it's text
				if chunk.Content != nil && chunk.Content.Type == multimodal.MediaTypeText {
					fullResponse += chunk.Content.Text
				}
				
				// If this is the final chunk, add the full response to history
				if chunk.IsFinal {
					ma.History.AddMessage("assistant", fullResponse)
					
					// Save the conversation
					if err := ma.saveConversation(); err != nil {
						// Just log the error
						fmt.Printf("Warning: Failed to save conversation: %v\n", err)
					}
				}
			}
		}
	}()
	
	return processedStream, nil
}

// AddSystemPrompt adds a system message to the conversation history
func (ma *MultimodalAgent) AddSystemPrompt(text string) {
	ma.History.AddMessage("system", text)
	ma.LLM.SetSystemPrompt(text)
}

// SetTemperature sets the temperature for generation
func (ma *MultimodalAgent) SetTemperature(temperature float64) {
	ma.Temperature = temperature
}

// SetMaxTokens sets the maximum tokens for generation
func (ma *MultimodalAgent) SetMaxTokens(maxTokens int) {
	ma.MaxTokens = maxTokens
}

// GetConversationHistory returns the conversation history
func (ma *MultimodalAgent) GetConversationHistory() *conversation.History {
	return ma.History
}

// Close cleans up resources and saves the final state
func (ma *MultimodalAgent) Close() error {
	// Save the conversation history
	if err := ma.saveConversation(); err != nil {
		fmt.Printf("Warning: Failed to save conversation during close: %v\n", err)
	}
	
	// Close the LLM
	if err := ma.LLM.Close(); err != nil {
		return fmt.Errorf("failed to close LLM: %w", err)
	}
	
	return nil
}

// Helper methods

// saveConversation saves the conversation history to a file
func (ma *MultimodalAgent) saveConversation() error {
	// Create the conversation file path
	filePath := filepath.Join(ma.ConversationDir, ma.History.GetID()+".json")
	
	// Save the conversation to a file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create conversation file: %w", err)
	}
	defer file.Close()
	
	// Serialize and save the conversation
	if err := ma.History.Save(file); err != nil {
		return fmt.Errorf("failed to save conversation to file: %w", err)
	}
	
	return nil
}

// extractTextFromInput extracts text content from a multimodal input
func extractTextFromInput(input *multimodal.Input) string {
	var text string
	
	for _, content := range input.Contents {
		if content.Type == multimodal.MediaTypeText {
			text += content.Text + " "
		}
	}
	
	return text
}

// extractTextFromOutput extracts text content from a multimodal output
func extractTextFromOutput(output *multimodal.Output) string {
	var text string
	
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			text += content.Text
		}
	}
	
	return text
}

// generateSystemPrompt creates a system prompt for the agent
func generateSystemPrompt(agent *Agent) string {
	// In a real implementation, this would generate a system prompt based on
	// the agent's definition, capabilities, etc.
	return fmt.Sprintf("You are %s, an AI assistant. Your purpose is to help the user.", agent.Name)
}
