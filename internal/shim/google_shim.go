package shim

import (
	"context"
	"fmt"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/shim/google"
)

// GoogleShim is an implementation of LLMShim for Google's Gemini models
type GoogleShim struct {
	config       Config
	systemPrompt string
	google       *google.GoogleShim
}

// NewGoogleShim creates a new shim for Google's Gemini
func NewGoogleShim(config Config) *GoogleShim {
	// Create the inner Google shim
	googleShim := google.NewGoogleShim(config.APIKey, config.Model)
	
	// Set endpoint if specified
	if config.Endpoint != "" {
		googleShim.Endpoint = config.Endpoint
	}
	
	// Set timeout if specified
	if config.Timeout > 0 {
		googleShim.Client.Timeout = config.Timeout
	}
	
	return &GoogleShim{
		config: config,
		google: googleShim,
	}
}

// Completion generates a text completion using Gemini
func (s *GoogleShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.CompletionWithContext(ctx, prompt, maxTokens, temperature)
}

// CompletionWithContext generates a text completion using Gemini with context
func (s *GoogleShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	if s.google != nil {
		// Set parameters for this request
		if maxTokens > 0 {
			s.google.MaxTokens = maxTokens
		}
		
		if temperature > 0 {
			s.google.Temperature = temperature
		}
		
		// Create messages with system prompt if set
		return s.google.GenerateContent(ctx, prompt, s.systemPrompt)
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return fmt.Sprintf("This is a response from Google Gemini (%s) using the prompt: %s", s.config.Model, prompt), nil
}

// MultimodalCompletion generates a multimodal completion using Gemini
func (s *GoogleShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.MultimodalCompletionWithContext(ctx, input)
}

// MultimodalCompletionWithContext generates a multimodal completion using Gemini with context
func (s *GoogleShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	if s.google != nil && s.SupportsMultimodal() {
		// Use the MultimodalGenerate function
		return s.google.MultimodalGenerate(ctx, input)
	}
	
	// Fallback to a placeholder if multimodal is not supported
	output := multimodal.NewOutput()
	output.AddText("This is a multimodal response from Google Gemini")
	
	return output, nil
}

// StreamCompletion streams a text completion from Gemini
func (s *GoogleShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	if s.google != nil {
		// Set parameters for this request
		if maxTokens > 0 {
			s.google.MaxTokens = maxTokens
		}
		
		if temperature > 0 {
			s.google.Temperature = temperature
		}
		
		// Start the streaming request
		return s.google.StreamContent(ctx, prompt, s.systemPrompt)
	}
	
	// Fallback to a placeholder stream
	resultCh := make(chan string)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming response ",
			"from Google Gemini using ",
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

// StreamMultimodalCompletion streams a multimodal completion from Gemini
func (s *GoogleShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	if s.google != nil && s.SupportsMultimodal() {
		// Use the StreamMultimodalContent function
		return s.google.StreamMultimodalContent(ctx, input)
	}
	
	// Fallback to a placeholder stream
	resultCh := make(chan *multimodal.Chunk)
	
	go func() {
		defer close(resultCh)
		
		// Simulate streaming with a few chunks
		chunks := []string{
			"This is a streaming multimodal response ",
			"from Google Gemini. ",
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

// SetSystemPrompt sets the system prompt for Gemini
func (s *GoogleShim) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// ParseSentinelfile parses a Sentinelfile into a configuration map
func (s *GoogleShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	if s.google != nil {
		// Save the current system prompt
		originalPrompt := s.systemPrompt
		
		// Set the system prompt for Sentinelfile parsing
		s.systemPrompt = `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
A Sentinelfile defines an AI agent's capabilities, behavior, and requirements. 
Your task is to extract information from the Sentinelfile and output a valid JSON structure with the following fields:
- name: The name of the agent (string, required)
- description: A short description of the agent (string, required)
- baseModel: The LLM model to use (string, required)
- capabilities: List of capabilities the agent should have (array of strings, optional)
- tools: List of tools the agent should have access to (array of strings, optional)
- stateSchema: Description of the state the agent should maintain (object, optional)
- parameters: Configuration parameters for the agent (object, optional)
- lifecycle: Object containing initialization and termination behaviors (object, optional)

Output ONLY the JSON object, no additional text or explanation.`
		
		// Restore the original system prompt when done
		defer func() {
			s.systemPrompt = originalPrompt
		}()
		
		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		// Generate the response
		response, err := s.CompletionWithContext(ctx, content, 4096, 0.2)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Sentinelfile: %v", err)
		}
		
		// Extract JSON from the response
		result, err := extractJSONFromResponse(response)
		if err != nil {
			return nil, fmt.Errorf("failed to extract JSON from response: %v", err)
		}
		
		return result, nil
	}
	
	// Fallback to a placeholder if the inner shim is not available
	return map[string]interface{}{
		"name":        "GoogleAgent",
		"description": "An agent parsed by Google Gemini",
		"baseModel":   s.config.Model,
	}, nil
}

// SupportsMultimodal returns whether this shim supports multimodal inputs
func (s *GoogleShim) SupportsMultimodal() bool {
	// Google multimodal-capable models
	multimodalModels := map[string]bool{
		"gemini-pro-vision":  true,
		"gemini-1.5-pro":     true,
		"gemini-1.5-flash":   true,
	}
	
	return multimodalModels[s.config.Model]
}

// Close cleans up any resources used by the shim
func (s *GoogleShim) Close() error {
	// No special cleanup needed
	return nil
}

// Helper function to extract JSON from a response
func extractJSONFromResponse(response string) (map[string]interface{}, error) {
	// Find the first { and last }
	start := strings.Index(response, "{")
	if start == -1 {
		return nil, fmt.Errorf("no JSON object found in response")
	}
	
	// Find the matching closing brace
	depth := 0
	end := -1
	
	for i := start; i < len(response); i++ {
		if response[i] == '{' {
			depth++
		} else if response[i] == '}' {
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}
	
	if end == -1 {
		return nil, fmt.Errorf("no valid JSON object found in response")
	}
	
	// Extract the JSON
	jsonStr := response[start : end+1]
	
	// Parse the JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	return result, nil
}
