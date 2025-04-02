package shim

import (
	"context"
	"fmt"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/shim/claude"
)

// ClaudeShim is an implementation of LLMShim for Anthropic's Claude models
type ClaudeShim struct {
	config       Config
	systemPrompt string
	claude       *claude.ClaudeShim
}

// NewClaudeShim creates a new shim for Claude
func NewClaudeShim(config Config) *ClaudeShim {
	// Create the inner Claude shim
	claudeShim := claude.NewClaudeShim(Config{
		Provider: config.Provider,
		Model:    config.Model,
		APIKey:   config.APIKey,
		Endpoint: config.Endpoint,
		Timeout:  config.Timeout,
	})
	
	return &ClaudeShim{
		config: config,
		claude: claudeShim,
	}
}

// Completion generates a text completion using Claude
func (s *ClaudeShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.CompletionWithContext(ctx, prompt, maxTokens, temperature)
}

// CompletionWithContext generates a text completion using Claude with context
func (s *ClaudeShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	// Use the inner Claude shim to make the request
	if s.claude != nil {
		return s.claude.CompletionWithContext(ctx, prompt, maxTokens, temperature)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return fmt.Sprintf("This is a response from Claude (%s) using the prompt: %s", s.config.Model, prompt), nil
}

// MultimodalCompletion generates a multimodal completion using Claude
func (s *ClaudeShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.MultimodalCompletionWithContext(ctx, input)
}

// MultimodalCompletionWithContext generates a multimodal completion using Claude with context
func (s *ClaudeShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	// Check if the Claude client is available
	if s.claude != nil {
		return s.claude.MultimodalCompletionWithContext(ctx, input)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	output := multimodal.NewOutput()
	output.AddText("This is a multimodal response from Claude")
	
	return output, nil
}

// StreamCompletion streams a text completion from Claude
func (s *ClaudeShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	// Check if the Claude client is available
	if s.claude != nil {
		return s.claude.StreamCompletion(ctx, prompt, maxTokens, temperature)
	}
	
	// Fallback to a placeholder stream if the inner shim is not available
	resultCh := make(chan string)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming response ",
			"from Claude using ",
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

// StreamMultimodalCompletion streams a multimodal completion from Claude
func (s *ClaudeShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// Check if the Claude client is available
	if s.claude != nil {
		return s.claude.StreamMultimodalCompletion(ctx, input)
	}
	
	// Fallback to a placeholder stream if the inner shim is not available
	resultCh := make(chan *multimodal.Chunk)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming multimodal response ",
			"from Claude. ",
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

// SetSystemPrompt sets the system prompt for Claude
func (s *ClaudeShim) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
	
	// Also set it in the inner Claude shim if available
	if s.claude != nil {
		s.claude.SetSystemPrompt(prompt)
	}
}

// ParseSentinelfile parses a Sentinelfile into a configuration map
func (s *ClaudeShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	// Use the inner Claude shim if available
	if s.claude != nil {
		return s.claude.ParseSentinelfile(content)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return map[string]interface{}{
		"name":        "ClaudeAgent",
		"description": "An agent parsed by Claude",
		"baseModel":   s.config.Model,
	}, nil
}

// SupportsMultimodal returns whether Claude supports multimodal inputs
func (s *ClaudeShim) SupportsMultimodal() bool {
	// Use the inner Claude shim if available
	if s.claude != nil {
		return s.claude.SupportsMultimodal()
	}
	
	// Claude Opus and newer models support multimodal
	return true
}

// Close cleans up any resources used by the Claude shim
func (s *ClaudeShim) Close() error {
	// No cleanup needed for Claude
	return nil
}
