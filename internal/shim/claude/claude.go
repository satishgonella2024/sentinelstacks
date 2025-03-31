package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

// ClaudeShim implements the LLM provider interface for Claude
type ClaudeShim struct {
	APIKey   string
	Model    string
	Client   *http.Client
	Endpoint string
}

// ChatMessage represents a message in the chat history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the Claude chat API
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// ChatResponse represents a response from the Claude chat API
type ChatResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// AnthropicMessage represents a message in the Anthropic API format
type AnthropicMessage struct {
	Role    string                 `json:"role"`
	Content []AnthropicMessagePart `json:"content"`
}

// AnthropicMessagePart represents a part of a message in the Anthropic API
type AnthropicMessagePart struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Image *struct {
		Source *struct {
			Type      string `json:"type"`
			MediaType string `json:"media_type"`
			Data      string `json:"data,omitempty"`
		} `json:"source,omitempty"`
	} `json:"image,omitempty"`
}

// AnthropicRequest represents a request to the Anthropic API
type AnthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []AnthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature float64            `json:"temperature,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
	System      string             `json:"system,omitempty"`
}

// AnthropicResponse represents a response from the Anthropic API
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
	Message      string `json:"message,omitempty"`
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

// Provider is the Claude provider implementation
type Provider struct {
	client   *http.Client
	apiKey   string
	endpoint string
	model    string
}

// NewClaudeShim creates a new ClaudeShim
func NewClaudeShim(apiKey, model string) *ClaudeShim {
	if model == "" {
		model = "claude-3-5-sonnet-20240627"
	}

	return &ClaudeShim{
		APIKey:   apiKey,
		Model:    model,
		Client:   &http.Client{},
		Endpoint: "https://api.anthropic.com/v1/messages",
	}
}

// CompleteChatPrompt sends a chat completion request to Claude
func (s *ClaudeShim) CompleteChatPrompt(messages []ChatMessage) (string, error) {
	chatReq := ChatRequest{
		Model:       s.Model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	req, err := http.NewRequest("POST", s.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Claude: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Claude API returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract the response text from the content array
	if len(chatResp.Content) > 0 {
		return chatResp.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in Claude response")
}

// ParseSentinelfile uses Claude to parse a Sentinelfile into a structured format
func (s *ClaudeShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	// Create a system prompt that instructs the model how to parse the Sentinelfile
	systemPrompt := `You are a specialized AI that parses natural language Sentinelfiles into structured JSON. 
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

	// Create the message array
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: content},
	}

	// Get the completion from the model
	completion, err := s.CompleteChatPrompt(messages)
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

	return result, nil
}

// extractJSON tries to extract a JSON object from a string
func extractJSON(s string) string {
	// Look for opening and closing braces
	start := 0
	for i, c := range s {
		if c == '{' {
			start = i
			break
		}
	}

	end := len(s) - 1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '}' {
			end = i
			break
		}
	}

	if start < end {
		return s[start : end+1]
	}

	return s // Return original if no JSON structure found
}

// NewProvider creates a new Claude provider
func NewProvider() interface{} {
	return &Provider{
		client:   &http.Client{Timeout: 60 * time.Second},
		endpoint: "https://api.anthropic.com/v1/messages",
		model:    "claude-3-5-sonnet-20240627",
	}
}

// Initialize initializes the provider with the given configuration
func (p *Provider) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["APIKey"].(string); ok {
		p.apiKey = apiKey
	} else {
		return fmt.Errorf("API key is required for Claude provider")
	}

	if endpoint, ok := config["Endpoint"].(string); ok && endpoint != "" {
		p.endpoint = endpoint
	}

	if model, ok := config["Model"].(string); ok && model != "" {
		p.model = model
	}

	return nil
}

// Name returns the name of the provider
func (p *Provider) Name() string {
	return "claude"
}

// AvailableModels returns the available models for this provider
func (p *Provider) AvailableModels() []string {
	return []string{
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
		"claude-3-5-sonnet-20240627",
	}
}

// GenerateResponse generates a response from Claude
func (p *Provider) GenerateResponse(ctx context.Context, prompt string, params map[string]interface{}) (string, error) {
	messages := []AnthropicMessage{
		{
			Role: "user",
			Content: []AnthropicMessagePart{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	// Create request
	request := AnthropicRequest{
		Model:    p.model,
		Messages: messages,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	} else {
		request.MaxTokens = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		request.Temperature = temperature
	} else {
		request.Temperature = 0.7 // Default
	}

	if systemPrompt, ok := params["system"].(string); ok {
		request.System = systemPrompt
	}

	// Make API request
	response, err := p.makeAnthropicRequest(ctx, request)
	if err != nil {
		return "", err
	}

	// Extract text from response
	var result string
	for _, part := range response.Content {
		if part.Type == "text" {
			result += part.Text
		}
	}

	return result, nil
}

// StreamResponse streams a response from Claude
func (p *Provider) StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error) {
	messages := []AnthropicMessage{
		{
			Role: "user",
			Content: []AnthropicMessagePart{
				{
					Type: "text",
					Text: prompt,
				},
			},
		},
	}

	// Create request
	request := AnthropicRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   true,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	} else {
		request.MaxTokens = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		request.Temperature = temperature
	} else {
		request.Temperature = 0.7 // Default
	}

	if systemPrompt, ok := params["system"].(string); ok {
		request.System = systemPrompt
	}

	return p.streamAnthropicRequest(ctx, request)
}

// GetEmbeddings gets embeddings for the given texts
func (p *Provider) GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// Simple implementation for now
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = []float32{0.1, 0.2, 0.3} // Dummy embeddings
	}
	return embeddings, nil
}

// SupportsMultimodal checks if the provider supports multimodal inputs and outputs
func (p *Provider) SupportsMultimodal() bool {
	// Claude 3 supports multimodal inputs
	return true
}

// GenerateMultimodalResponse generates a multimodal response from Claude
func (p *Provider) GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error) {
	// Convert multimodal input to Anthropic format
	anthropicMessages, err := p.convertToAnthropicMessages(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to Anthropic format: %w", err)
	}

	// Create request
	request := AnthropicRequest{
		Model:    p.model,
		Messages: anthropicMessages,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	} else if input.MaxTokens > 0 {
		request.MaxTokens = input.MaxTokens
	} else {
		request.MaxTokens = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		request.Temperature = temperature
	} else if input.Temperature > 0 {
		request.Temperature = input.Temperature
	} else {
		request.Temperature = 0.7 // Default
	}

	if systemPrompt, ok := params["system"].(string); ok {
		request.System = systemPrompt
	} else if sysPrompt, ok := input.Metadata["system"].(string); ok {
		request.System = sysPrompt
	}

	// Make API request
	response, err := p.makeAnthropicRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// Convert Anthropic response to multimodal output
	output := multimodal.NewOutput()

	// Add content from response
	for _, part := range response.Content {
		if part.Type == "text" {
			output.AddText(part.Text)
		}
		// Handle other content types if needed
	}

	// Add usage information
	output.Metadata = map[string]interface{}{
		"input_tokens":  response.Usage.InputTokens,
		"output_tokens": response.Usage.OutputTokens,
		"used_tokens":   response.Usage.OutputTokens,
		"stop_reason":   response.StopReason,
	}

	return output, nil
}

// StreamMultimodalResponse streams a multimodal response from Claude
func (p *Provider) StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error) {
	// Convert multimodal input to Anthropic format
	anthropicMessages, err := p.convertToAnthropicMessages(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to Anthropic format: %w", err)
	}

	// Create request
	request := AnthropicRequest{
		Model:    p.model,
		Messages: anthropicMessages,
		Stream:   true,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	} else if input.MaxTokens > 0 {
		request.MaxTokens = input.MaxTokens
	} else {
		request.MaxTokens = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		request.Temperature = temperature
	} else if input.Temperature > 0 {
		request.Temperature = input.Temperature
	} else {
		request.Temperature = 0.7 // Default
	}

	if systemPrompt, ok := params["system"].(string); ok {
		request.System = systemPrompt
	} else if sysPrompt, ok := input.Metadata["system"].(string); ok {
		request.System = sysPrompt
	}

	// Set up streaming channel for Anthropic
	textChunks, err := p.streamAnthropicRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// Create output channel
	resultCh := make(chan *multimodal.Chunk)

	// Start goroutine to read from Anthropic and write to output
	go func() {
		defer close(resultCh)

		for text := range textChunks {
			// Create multimodal chunk with text content
			chunk := &multimodal.Chunk{
				Content: multimodal.NewTextContent(text),
				IsFinal: false,
			}
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				// Successfully sent chunk
			}
		}

		// Send final empty chunk to indicate completion
		finalChunk := &multimodal.Chunk{
			Content: multimodal.NewTextContent(""),
			IsFinal: true,
		}

		select {
		case <-ctx.Done():
			return
		case resultCh <- finalChunk:
			// Successfully sent final chunk
		}
	}()

	return resultCh, nil
}

// Helper methods

// convertToAnthropicMessages converts a multimodal input to Anthropic format
func (p *Provider) convertToAnthropicMessages(input *multimodal.Input) ([]AnthropicMessage, error) {
	if len(input.Contents) == 0 {
		return nil, fmt.Errorf("input contains no content")
	}

	// Claude API expects a single user message
	anthropicMessage := AnthropicMessage{
		Role:    "user",
		Content: []AnthropicMessagePart{},
	}

	// Convert each content item
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text part
			anthropicMessage.Content = append(anthropicMessage.Content, AnthropicMessagePart{
				Type: "text",
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Add image part
			part := AnthropicMessagePart{
				Type: "image",
				Image: &struct {
					Source *struct {
						Type      string `json:"type"`
						MediaType string `json:"media_type"`
						Data      string `json:"data,omitempty"`
					} `json:"source,omitempty"`
				}{
					Source: &struct {
						Type      string `json:"type"`
						MediaType string `json:"media_type"`
						Data      string `json:"data,omitempty"`
					}{
						Type:      "base64",
						MediaType: content.MimeType,
						Data:      content.ToBase64(),
					},
				},
			}
			anthropicMessage.Content = append(anthropicMessage.Content, part)

		default:
			return nil, fmt.Errorf("unsupported content type: %s", content.Type)
		}
	}

	// Create messages array with the single user message
	messages := []AnthropicMessage{anthropicMessage}

	return messages, nil
}

// makeAnthropicRequest makes a request to the Anthropic API
func (p *Provider) makeAnthropicRequest(ctx context.Context, request AnthropicRequest) (*AnthropicResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", p.apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Claude: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var response AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// streamAnthropicRequest streams a request to the Anthropic API
func (p *Provider) streamAnthropicRequest(ctx context.Context, request AnthropicRequest) (<-chan string, error) {
	// Ensure streaming is enabled
	request.Stream = true

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", p.apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")
	req.Header.Set("Accept", "application/json")

	// Create output channel
	ch := make(chan string)

	// Start goroutine to handle streaming
	go func() {
		defer close(ch)

		// Send request
		resp, err := p.client.Do(req)
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

// Helper functions needed for streaming
