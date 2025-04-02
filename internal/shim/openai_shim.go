package shim

import (
	"context"
	"fmt"
	"time"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim/openai"
)

// OpenAIShim is an implementation of LLMShim for OpenAI models
type OpenAIShim struct {
	config       Config
	systemPrompt string
	openai       *openai.OpenAIShim
}

// NewOpenAIShim creates a new shim for OpenAI
func NewOpenAIShim(config Config) *OpenAIShim {
	// Create the inner OpenAI shim
	openaiShim := openai.NewOpenAIShim(config.APIKey, config.Model)
	
	// Set endpoint if specified
	if config.Endpoint != "" {
		openaiShim.Endpoint = config.Endpoint
	}
	
	// Set timeout if specified
	if config.Timeout > 0 {
		openaiShim.Client.Timeout = config.Timeout
	}
	
	return &OpenAIShim{
		config: config,
		openai: openaiShim,
	}
}

// Completion generates a text completion using OpenAI
func (s *OpenAIShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.CompletionWithContext(ctx, prompt, maxTokens, temperature)
}

// CompletionWithContext generates a text completion using OpenAI with context
func (s *OpenAIShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	if s.openai != nil {
		// Set parameters for this request
		if maxTokens > 0 {
			s.openai.MaxTokens = maxTokens
		}
		
		if temperature > 0 {
			s.openai.Temperature = temperature
		}
		
		// Create messages with system prompt if set
		messages := []openai.ChatMessage{}
		
		if s.systemPrompt != "" {
			messages = append(messages, openai.ChatMessage{
				Role:    "system",
				Content: s.systemPrompt,
			})
		}
		
		// Add user message
		messages = append(messages, openai.ChatMessage{
			Role:    "user",
			Content: prompt,
		})
		
		// Make the request
		return s.openai.CompleteChatPrompt(messages)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return fmt.Sprintf("This is a response from OpenAI (%s) using the prompt: %s", s.config.Model, prompt), nil
}

// MultimodalCompletion generates a multimodal completion using OpenAI
func (s *OpenAIShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.MultimodalCompletionWithContext(ctx, input)
}

// MultimodalCompletionWithContext generates a multimodal completion using OpenAI with context
func (s *OpenAIShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	// For now, we don't have a direct connection to the openai.OpenAIShim multimodal methods
	// This would need to be implemented more comprehensively
	
	// Fallback to a placeholder if the inner shim is not available
	output := multimodal.NewOutput()
	output.AddText("This is a multimodal response from OpenAI")
	
	return output, nil
}

// StreamCompletion streams a text completion from OpenAI
func (s *OpenAIShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	// For now, we don't have a direct connection to the openai.OpenAIShim streaming methods
	// This would need to be implemented more comprehensively
	
	// Fallback to a placeholder stream
	resultCh := make(chan string)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming response ",
			"from OpenAI using ",
			"the prompt: " + prompt,
		}
		
		for _, chunk := range chunks {
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	
	return resultCh, nil
}

// StreamMultimodalCompletion streams a multimodal completion from OpenAI
func (s *OpenAIShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// For now, we don't have a direct connection to the openai.OpenAIShim streaming methods
	// This would need to be implemented more comprehensively
	
	// Fallback to a placeholder stream
	resultCh := make(chan *multimodal.Chunk)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming multimodal response ",
			"from OpenAI. ",
			"I'm analyzing the content you provided.",
		}
		
		for i, text := range chunks {
			chunk := multimodal.NewChunk(multimodal.NewTextContent(text), i == len(chunks)-1)
			
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	
	return resultCh, nil
}

// SetSystemPrompt sets the system prompt for OpenAI
func (s *OpenAIShim) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// ParseSentinelfile parses a Sentinelfile into a configuration map
func (s *OpenAIShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	if s.openai != nil {
		return s.openai.ParseSentinelfile(content)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return map[string]interface{}{
		"name":        "OpenAIAgent",
		"description": "An agent parsed by OpenAI",
		"baseModel":   s.config.Model,
	}, nil
}

// SupportsMultimodal returns whether OpenAI supports multimodal inputs
func (s *OpenAIShim) SupportsMultimodal() bool {
	// These models support multimodal inputs
	multimodalModels := map[string]bool{
		"gpt-4-vision-preview": true,
		"gpt-4-turbo-preview":  true,
		"gpt-4-turbo":          true,
		"gpt-4-1106-vision-preview": true,
		"gpt-4-1106-preview":   true,
		"gpt-4-vision":         true,
	}
	
	return multimodalModels[s.config.Model]
}

// Close cleans up any resources used by the OpenAI shim
func (s *OpenAIShim) Close() error {
	// No cleanup needed for OpenAI
	return nil
}
