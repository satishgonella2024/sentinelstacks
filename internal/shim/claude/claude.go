package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/shim"
)

// Constants for Claude API
const (
	DefaultEndpoint    = "https://api.anthropic.com/v1/messages"
	DefaultModel       = "claude-3-5-sonnet-20240627"
	DefaultMaxTokens   = 4096
	DefaultTemperature = 0.7
	AnthropicVersion   = "2023-06-01"
)

// ClaudeShim implements the LLMShim interface for Claude
type ClaudeShim struct {
	APIKey       string
	Model        string
	Client       *http.Client
	Endpoint     string
	SystemPrompt string
	MaxRetries   int
}

// AnthropicMessage represents a message in the Claude API format
type AnthropicMessage struct {
	Role    string                 `json:"role"`
	Content []AnthropicMessagePart `json:"content"`
}

// AnthropicMessagePart represents a part of a message in the Claude API
type AnthropicMessagePart struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Image *struct {
		Source struct {
			Type      string `json:"type"`
			MediaType string `json:"media_type"`
			Data      string `json:"data"`
		} `json:"source"`
	} `json:"image,omitempty"`
}

// AnthropicRequest represents a request to the Claude API
type AnthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []AnthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature float64            `json:"temperature,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
	System      string             `json:"system,omitempty"`
}

// AnthropicResponse represents a response from the Claude API
type AnthropicResponse struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Role         string                 `json:"role"`
	Content      []AnthropicMessagePart `json:"content"`
	StopReason   string                 `json:"stop_reason,omitempty"`
	StopSequence string                 `json:"stop_sequence,omitempty"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// AnthropicStreamingChunk represents a chunk in a streaming response
type AnthropicStreamingChunk struct {
	Type         string `json:"type"`
	Index        int    `json:"index,omitempty"`
	ContentBlock struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"content_block,omitempty"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"delta,omitempty"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewClaudeShim creates a new ClaudeShim with config
func NewClaudeShim(config shim.Config) *ClaudeShim {
	// Set defaults if not provided
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	model := config.Model
	if model == "" {
		model = DefaultModel
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	return &ClaudeShim{
		APIKey:     config.APIKey,
		Model:      model,
		Client:     &http.Client{Timeout: timeout},
		Endpoint:   endpoint,
		MaxRetries: 2, // Default retry count
	}
}

// Completion sends a text prompt to Claude and returns the response
func (s *ClaudeShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.CompletionWithContext(ctx, prompt, maxTokens, temperature)
}

// CompletionWithContext sends a text prompt to Claude with context and returns the response
func (s *ClaudeShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	// Set defaults if not provided
	if maxTokens <= 0 {
		maxTokens = DefaultMaxTokens
	}

	if temperature <= 0 {
		temperature = DefaultTemperature
	}

	// Create message content for the request
	content := []AnthropicMessagePart{
		{
			Type: "text",
			Text: prompt,
		},
	}

	// Create a single user message
	messages := []AnthropicMessage{
		{
			Role:    "user",
			Content: content,
		},
	}

	// Create request
	request := AnthropicRequest{
		Model:       s.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Add system prompt if set
	if s.SystemPrompt != "" {
		request.System = s.SystemPrompt
	}

	// Make request with retries
	var response *AnthropicResponse
	var err error
	for attempt := 0; attempt <= s.MaxRetries; attempt++ {
		// Make the request
		response, err = s.makeRequest(ctx, request)
		
		// If successful or not a retryable error, break
		if err == nil || !isRetryableError(err) {
			break
		}
		
		// Calculate backoff delay (exponential backoff)
		delay := time.Duration(attempt+1) * 500 * time.Millisecond
		
		// Create a timer for the delay
		timer := time.NewTimer(delay)
		
		// Wait for the timer or context cancellation
		select {
		case <-ctx.Done():
			// Context canceled, abort retries
			timer.Stop()
			return "", ctx.Err()
		case <-timer.C:
			// Timer expired, continue to next attempt
		}
	}
	
	// If we still have an error after retries, return it
	if err != nil {
		return "", fmt.Errorf("Claude API request failed after retries: %w", err)
	}
	
	// Extract the text from the response
	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude API")
	}
	
	var responseText string
	for _, content := range response.Content {
		if content.Type == "text" {
			responseText += content.Text
		}
	}
	
	return responseText, nil
}

// MultimodalCompletion sends a multimodal input to Claude and returns the output
func (s *ClaudeShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.MultimodalCompletionWithContext(ctx, input)
}

// MultimodalCompletionWithContext sends a multimodal input to Claude with context and returns the output
func (s *ClaudeShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	// Check if model supports multimodal
	if !s.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", s.Model)
	}
	
	// Create Anthropic-format messages from multimodal input
	message, err := s.createAnthropicMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create message from input: %w", err)
	}
	
	// Create request
	request := AnthropicRequest{
		Model:    s.Model,
		Messages: []AnthropicMessage{message},
	}
	
	// Set parameters
	if input.MaxTokens > 0 {
		request.MaxTokens = input.MaxTokens
	} else {
		request.MaxTokens = DefaultMaxTokens
	}
	
	if input.Temperature > 0 {
		request.Temperature = input.Temperature
	} else {
		request.Temperature = DefaultTemperature
	}
	
	// Add system prompt if set
	if s.SystemPrompt != "" {
		request.System = s.SystemPrompt
	} else if systemPrompt, ok := input.Metadata["system"].(string); ok {
		request.System = systemPrompt
	}
	
	// Make request with retries
	var response *AnthropicResponse
	for attempt := 0; attempt <= s.MaxRetries; attempt++ {
		// Make the request
		response, err = s.makeRequest(ctx, request)
		
		// If successful or not a retryable error, break
		if err == nil || !isRetryableError(err) {
			break
		}
		
		// Calculate backoff delay (exponential backoff)
		delay := time.Duration(attempt+1) * 500 * time.Millisecond
		
		// Wait for the delay or context cancellation
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
			// Timer expired, continue to next attempt
		}
	}
	
	// If we still have an error after retries, return it
	if err != nil {
		return nil, fmt.Errorf("Claude API request failed after retries: %w", err)
	}
	
	// Create multimodal output from response
	output := multimodal.NewOutput()
	
	// Extract content from response
	for _, content := range response.Content {
		if content.Type == "text" {
			output.AddText(content.Text)
		}
		// Handle other content types if needed
	}
	
	// Add metadata
	output.Metadata = map[string]interface{}{
		"model":         s.Model,
		"input_tokens":  response.Usage.InputTokens,
		"output_tokens": response.Usage.OutputTokens,
		"stop_reason":   response.StopReason,
	}
	
	return output, nil
}

// StreamCompletion streams a text completion from Claude
func (s *ClaudeShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	// Set defaults if not provided
	if maxTokens <= 0 {
		maxTokens = DefaultMaxTokens
	}
	
	if temperature <= 0 {
		temperature = DefaultTemperature
	}
	
	// Create message content for the request
	content := []AnthropicMessagePart{
		{
			Type: "text",
			Text: prompt,
		},
	}
	
	// Create a single user message
	messages := []AnthropicMessage{
		{
			Role:    "user",
			Content: content,
		},
	}
	
	// Create request
	request := AnthropicRequest{
		Model:       s.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Stream:      true,
	}
	
	// Add system prompt if set
	if s.SystemPrompt != "" {
		request.System = s.SystemPrompt
	}
	
	// Make streaming request
	return s.streamRequest(ctx, request)
}

// StreamMultimodalCompletion streams a multimodal completion from Claude
func (s *ClaudeShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// Check if model supports multimodal
	if !s.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", s.Model)
	}
	
	// Create Anthropic-format messages from multimodal input
	message, err := s.createAnthropicMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create message from input: %w", err)
	}
	
	// Create request
	request := AnthropicRequest{
		Model:    s.Model,
		Messages: []AnthropicMessage{message},
		Stream:   true,
	}
	
	// Set parameters
	if input.MaxTokens > 0 {
		request.MaxTokens = input.MaxTokens
	} else {
		request.MaxTokens = DefaultMaxTokens
	}
	
	if input.Temperature > 0 {
		request.Temperature = input.Temperature
	} else {
		request.Temperature = DefaultTemperature
	}
	
	// Add system prompt if set
	if s.SystemPrompt != "" {
		request.System = s.SystemPrompt
	} else if systemPrompt, ok := input.Metadata["system"].(string); ok {
		request.System = systemPrompt
	}
	
	// Make streaming request
	textStream, err := s.streamRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	
	// Convert text stream to multimodal chunks
	return s.convertTextStreamToMultimodal(ctx, textStream), nil
}

// SetSystemPrompt sets the system prompt for the model
func (s *ClaudeShim) SetSystemPrompt(prompt string) {
	s.SystemPrompt = prompt
}

// ParseSentinelfile uses Claude to parse a Sentinelfile into a structured format
func (s *ClaudeShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	// Create a system prompt that instructs the model how to parse the Sentinelfile
	origPrompt := s.SystemPrompt
	s.SystemPrompt = `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
A Sentinelfile defines an AI agent's capabilities, behavior, and requirements. 
Your task is to extract information from the Sentinelfile and output a valid JSON structure with the following fields:
- name: The name of the agent (string, required)
- description: A short description of the agent (string, required)
- baseModel: The LLM model to use (string, required)
- capabilities: List of capabilities the agent should have (array of strings, optional)
- tools: List of tools the agent should have access to (array of strings, optional)
- stateSchema: Description of the state the agent should maintain (object, optional)
- parameters: Configuration parameters for the agent (object, optional)

Output ONLY the JSON object, no additional text or explanation.`

	defer func() {
		s.SystemPrompt = origPrompt
	}()

	// Create the prompt for the model
	prompt := fmt.Sprintf("Parse the following Sentinelfile into JSON:\n\n%s", content)

	// Send the completion request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	completion, err := s.CompletionWithContext(ctx, prompt, DefaultMaxTokens, 0.2)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion: %w", err)
	}

	// Extract the JSON part
	jsonStr := extractJSON(completion)

	// Parse the JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Validate required fields
	required := []string{"name", "description", "baseModel"}
	for _, field := range required {
		if _, ok := result[field]; !ok {
			return nil, fmt.Errorf("missing required field '%s' in Sentinelfile", field)
		}
	}

	return result, nil
}

// SupportsMultimodal returns whether this shim supports multimodal inputs
func (s *ClaudeShim) SupportsMultimodal() bool {
	// Check if model is in the list of multimodal-capable models
	multimodalModels := map[string]bool{
		"claude-3-opus-20240229":   true,
		"claude-3-sonnet-20240229": true,
		"claude-3-haiku-20240307":  true,
		"claude-3-5-sonnet-20240627": true,
		"claude-3-5-sonnet": true,
		"claude-3-opus": true,
		"claude-3-sonnet": true,
		"claude-3-haiku": true,
	}

	return multimodalModels[s.Model]
}

// Close cleans up any resources used by the shim
func (s *ClaudeShim) Close() error {
	// No special cleanup needed
	return nil
}

// Helper methods

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	// Rate limit errors, network errors, and 5xx server errors are retryable
	errStr := err.Error()
	return strings.Contains(errStr, "rate limit") || 
		strings.Contains(errStr, "timeout") || 
		strings.Contains(errStr, "connection") || 
		strings.Contains(errStr, "status code 429") || 
		strings.Contains(errStr, "status code 5")
}

// makeRequest makes a request to the Claude API
func (s *ClaudeShim) makeRequest(ctx context.Context, request AnthropicRequest) (*AnthropicResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("Anthropic-Version", AnthropicVersion)

	// Send the request
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Claude API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response AnthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// streamRequest streams a request to the Claude API
func (s *ClaudeShim) streamRequest(ctx context.Context, request AnthropicRequest) (<-chan string, error) {
	// Ensure streaming is enabled
	request.Stream = true

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("Anthropic-Version", AnthropicVersion)
	req.Header.Set("Accept", "application/json")

	// Create output channel
	ch := make(chan string)

	// Start goroutine to handle streaming
	go func() {
		defer close(ch)

		// Send request
		resp, err := s.Client.Do(req)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case ch <- fmt.Sprintf("Error: %v", err):
				return
			}
		}
		defer resp.Body.Close()

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			select {
			case <-ctx.Done():
				return
			case ch <- fmt.Sprintf("Error: API returned status code %d: %s", resp.StatusCode, string(body)):
				return
			}
		}

		// Process streaming response
		reader := bufio.NewReader(resp.Body)
		for {
			// Check if context is done
			select {
			case <-ctx.Done():
				return
			default:
				// Continue processing
			}

			// Read a line
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				select {
				case <-ctx.Done():
					return
				case ch <- fmt.Sprintf("Error reading stream: %v", err):
					return
				}
			}

			// Skip empty lines
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			// Skip "data: " prefix
			if bytes.HasPrefix(line, []byte("data: ")) {
				line = bytes.TrimPrefix(line, []byte("data: "))
			}

			// Skip "[DONE]" marker
			if string(line) == "[DONE]" {
				break
			}

			// Parse the JSON chunk
			var chunk AnthropicStreamingChunk
			if err := json.Unmarshal(line, &chunk); err != nil {
				select {
				case <-ctx.Done():
					return
				case ch <- fmt.Sprintf("Error parsing chunk: %v", err):
					return
				}
				continue
			}

			// Check for errors
			if chunk.Type == "error" {
				select {
				case <-ctx.Done():
					return
				case ch <- fmt.Sprintf("Error from API: %s", chunk.Error.Message):
					return
				}
			}

			// Extract content based on chunk type
			var content string
			if chunk.Type == "content_block_delta" {
				content = chunk.Delta.Text
			} else if chunk.Type == "content_block_start" {
				content = chunk.ContentBlock.Text
			}

			// Send content if not empty
			if content != "" {
				select {
				case <-ctx.Done():
					return
				case ch <- content:
					// Successfully sent content
				}
			}
		}
	}()

	return ch, nil
}

// createAnthropicMessage creates an AnthropicMessage from a multimodal input
func (s *ClaudeShim) createAnthropicMessage(input *multimodal.Input) (AnthropicMessage, error) {
	// Create the message
	message := AnthropicMessage{
		Role:    "user",
		Content: []AnthropicMessagePart{},
	}

	// Process each content item
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text content
			message.Content = append(message.Content, AnthropicMessagePart{
				Type: "text",
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Skip if no data and no URI
			if len(content.Data) == 0 && content.URI == "" {
				continue
			}

			// Create image part
			imagePart := AnthropicMessagePart{
				Type:  "image",
				Image: &struct{
					Source struct {
						Type      string `json:"type"`
						MediaType string `json:"media_type"`
						Data      string `json:"data"`
					} `json:"source"`
				}{},
			}

			// If we have data, use base64 encoding
			if len(content.Data) > 0 {
				base64Data := base64.StdEncoding.EncodeToString(content.Data)
				imagePart.Image.Source.Type = "base64"
				imagePart.Image.Source.MediaType = content.MimeType
				imagePart.Image.Source.Data = base64Data
			} else if content.URI != "" {
				// If we have a URI, use it - only for publicly accessible URLs
				// Note: This won't work for most cases as Claude API doesn't support URL references
				// Future: might need to download and convert to base64
				imagePart.Image.Source.Type = "url"
				imagePart.Image.Source.MediaType = content.MimeType
				imagePart.Image.Source.Data = content.URI
			}

			// Add the image part
			message.Content = append(message.Content, imagePart)

		default:
			return message, fmt.Errorf("unsupported content type: %s", content.Type)
		}
	}

	// If no content was added, return an error
	if len(message.Content) == 0 {
		return message, fmt.Errorf("no valid content to send to Claude")
	}

	return message, nil
}

// convertTextStreamToMultimodal converts a text stream to a multimodal chunk stream
func (s *ClaudeShim) convertTextStreamToMultimodal(ctx context.Context, textStream <-chan string) <-chan *multimodal.Chunk {
	chunkStream := make(chan *multimodal.Chunk)

	go func() {
		defer close(chunkStream)

		var isError bool
		for text := range textStream {
			// Check if this is an error message
			if strings.HasPrefix(text, "Error:") {
				chunk := multimodal.NewChunk(multimodal.NewTextContent(text), true)
				chunk.Error = fmt.Errorf(text)
				
				select {
				case <-ctx.Done():
					return
				case chunkStream <- chunk:
					isError = true
				}
				break
			}

			// Create a normal chunk
			chunk := multimodal.NewChunk(multimodal.NewTextContent(text), false)
			
			select {
			case <-ctx.Done():
				return
			case chunkStream <- chunk:
				// Successfully sent chunk
			}
		}

		// If we didn't encounter an error, send a final chunk
		if !isError {
			finalChunk := multimodal.NewChunk(multimodal.NewTextContent(""), true)
			
			select {
			case <-ctx.Done():
				return
			case chunkStream <- finalChunk:
				// Successfully sent final chunk
			}
		}
	}()

	return chunkStream
}

// extractJSON extracts a JSON object from a string
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	if start == -1 {
		return s
	}

	// Find matching closing brace
	depth := 0
	for i := start; i < len(s); i++ {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}

	// If no matching closing brace found, return from start to end
	return s[start:]
}
