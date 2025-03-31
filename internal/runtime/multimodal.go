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
	Shim            shim.Shim
	ShimConfig      shim.Config
	MaxTokens       int
	Temperature     float64
	ConversationDir string
}

// NewMultimodalAgent creates a new multimodal agent
func NewMultimodalAgent(agent *Agent, shimConfig shim.Config) (*MultimodalAgent, error) {
	// Create conversation directory
	conversationDir := filepath.Join(agent.StateDir, "conversations")
	if err := os.MkdirAll(conversationDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create conversation directory: %w", err)
	}

	// Create a new conversation history
	history := conversation.NewHistory(agent.ID, fmt.Sprintf("session_%d", time.Now().UnixNano()))

	// Create the shim
	shimInstance, err := shim.CreateShim(shimConfig.Provider, shimConfig.Model, shimConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create shim: %w", err)
	}

	// Check if the shim supports multimodal if needed
	if !shimInstance.SupportsMultimodal() {
		return nil, fmt.Errorf("provider %s with model %s does not support multimodal", shimConfig.Provider, shimConfig.Model)
	}

	// Create the multimodal agent
	return &MultimodalAgent{
		Agent:           agent,
		History:         history,
		Shim:            shimInstance,
		ShimConfig:      shimConfig,
		MaxTokens:       4096, // Default
		Temperature:     0.7,  // Default
		ConversationDir: conversationDir,
	}, nil
}

// ProcessTextInput processes text input from the user
func (ma *MultimodalAgent) ProcessTextInput(ctx context.Context, text string) (string, error) {
	// Add the user message to the history
	ma.History.AddUserMessage(text)

	// Convert the conversation to a multimodal input
	input, err := ma.History.ToMultimodalInput(10) // Use last 10 messages
	if err != nil {
		return "", fmt.Errorf("failed to convert history to input: %w", err)
	}

	// Set generation parameters
	input.MaxTokens = ma.MaxTokens
	input.Temperature = ma.Temperature

	// Generate a response
	output, err := ma.Shim.GenerateMultimodal(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract text from the output
	var responseText string
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			responseText = content.Text
			break
		}
	}

	// Add the assistant's response to the history
	ma.History.AddAssistantMessage(responseText)

	// Save the conversation
	conversationFile := filepath.Join(ma.ConversationDir, ma.History.ID+".json")
	if err := ma.History.SaveToFile(conversationFile); err != nil {
		return responseText, fmt.Errorf("failed to save conversation: %w", err)
	}

	return responseText, nil
}

// ProcessMultimodalInput processes multimodal input from the user (text + images)
func (ma *MultimodalAgent) ProcessMultimodalInput(ctx context.Context, contents []*multimodal.Content) (*multimodal.Output, error) {
	// Add the user message to the history
	ma.History.AddUserMultimodalMessage(contents)

	// Convert the conversation to a multimodal input
	input, err := ma.History.ToMultimodalInput(10) // Use last 10 messages
	if err != nil {
		return nil, fmt.Errorf("failed to convert history to input: %w", err)
	}

	// Set generation parameters
	input.MaxTokens = ma.MaxTokens
	input.Temperature = ma.Temperature

	// Generate a response
	output, err := ma.Shim.GenerateMultimodal(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract text from the output
	var responseText string
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			responseText = content.Text
			break
		}
	}

	// Add the assistant's response to the history
	ma.History.AddAssistantMessage(responseText)

	// Save the conversation
	conversationFile := filepath.Join(ma.ConversationDir, ma.History.ID+".json")
	if err := ma.History.SaveToFile(conversationFile); err != nil {
		return output, fmt.Errorf("failed to save conversation: %w", err)
	}

	return output, nil
}

// StreamMultimodalInput processes multimodal input and streams the response
func (ma *MultimodalAgent) StreamMultimodalInput(ctx context.Context, contents []*multimodal.Content) (<-chan *multimodal.Chunk, error) {
	// Add the user message to the history
	ma.History.AddUserMultimodalMessage(contents)

	// Convert the conversation to a multimodal input
	input, err := ma.History.ToMultimodalInput(10) // Use last 10 messages
	if err != nil {
		return nil, fmt.Errorf("failed to convert history to input: %w", err)
	}

	// Set generation parameters
	input.MaxTokens = ma.MaxTokens
	input.Temperature = ma.Temperature
	input.Stream = true

	// Stream the response
	chunks, err := ma.Shim.StreamMultimodal(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stream response: %w", err)
	}

	// Create a channel for processed chunks
	processedChunks := make(chan *multimodal.Chunk)

	// Process chunks in a goroutine
	go func() {
		defer close(processedChunks)

		var fullResponse string

		// Process each chunk
		for chunk := range chunks {
			// Forward the chunk to the caller
			select {
			case <-ctx.Done():
				return
			case processedChunks <- chunk:
				// Accumulate the text response
				if chunk.Content.Type == multimodal.MediaTypeText {
					fullResponse += chunk.Content.Text
				}

				// If this is the final chunk, add the response to history
				if chunk.IsFinal {
					// Add the assistant's response to the history
					ma.History.AddAssistantMessage(fullResponse)

					// Save the conversation
					conversationFile := filepath.Join(ma.ConversationDir, ma.History.ID+".json")
					_ = ma.History.SaveToFile(conversationFile) // Ignore error
				}
			}
		}
	}()

	return processedChunks, nil
}

// AddSystemPrompt adds a system prompt to the conversation
func (ma *MultimodalAgent) AddSystemPrompt(text string) {
	ma.History.AddSystemMessage(text)
}

// GetConversationHistory returns the full conversation history
func (ma *MultimodalAgent) GetConversationHistory() *conversation.History {
	return ma.History
}

// SetTemperature sets the temperature for generation
func (ma *MultimodalAgent) SetTemperature(temperature float64) {
	ma.Temperature = temperature
}

// SetMaxTokens sets the maximum tokens for generation
func (ma *MultimodalAgent) SetMaxTokens(maxTokens int) {
	ma.MaxTokens = maxTokens
}

// Close closes the agent resources
func (ma *MultimodalAgent) Close() error {
	// Save the conversation
	conversationFile := filepath.Join(ma.ConversationDir, ma.History.ID+".json")
	if err := ma.History.SaveToFile(conversationFile); err != nil {
		return fmt.Errorf("failed to save conversation: %w", err)
	}

	// Close the shim
	return ma.Shim.Close()
}
