package openai

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

// OpenAIShim implements the LLM provider interface for OpenAI
type OpenAIShim struct {
	APIKey      string
	Model       string
	Client      *http.Client
	Endpoint    string
	Temperature float64
	MaxTokens   int
}

// ChatMessage represents a message in the chat history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the OpenAI chat API
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatResponse represents a response from the OpenAI chat API
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatStreamResponse represents a streaming response from the OpenAI chat API
type ChatStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// MultimediaMessage represents a message with multimedia content for OpenAI
type MultimediaMessage struct {
	Role    string              `json:"role"`
	Content []MultimediaContent `json:"content"`
}

// MultimediaContent represents multimedia content in a message
type MultimediaContent struct {
	Type  string                 `json:"type"`
	Text  string                 `json:"text,omitempty"`
	Image *MultimediaImageDetail `json:"image,omitempty"`
}

// MultimediaImageDetail represents image details for multimedia content
type MultimediaImageDetail struct {
	URL    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Provider is the OpenAI provider implementation
type Provider struct {
	client   *http.Client
	apiKey   string
	endpoint string
	model    string
}

// NewOpenAIShim creates a new OpenAIShim
func NewOpenAIShim(apiKey, model string) *OpenAIShim {
	if model == "" {
		model = "gpt-4-turbo"
	}

	return &OpenAIShim{
		APIKey:      apiKey,
		Model:       model,
		Client:      &http.Client{},
		Endpoint:    "https://api.openai.com/v1/chat/completions",
		Temperature: 0.7,
		MaxTokens:   4096,
	}
}

// CompleteChatPrompt sends a chat completion request to OpenAI
func (s *OpenAIShim) CompleteChatPrompt(messages []ChatMessage) (string, error) {
	chatReq := ChatRequest{
		Model:       s.Model,
		Messages:    messages,
		Temperature: s.Temperature,
		MaxTokens:   s.MaxTokens,
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
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract the response text from the choices array
	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in OpenAI response")
}

// ParseSentinelfile uses OpenAI to parse a Sentinelfile into a structured format
func (s *OpenAIShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
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

// NewProvider creates a new OpenAI provider
func NewProvider() interface{} {
	return &Provider{
		client:   &http.Client{Timeout: 60 * time.Second},
		endpoint: "https://api.openai.com/v1/chat/completions",
		model:    "gpt-4-vision-preview",
	}
}

// Initialize initializes the provider with the given configuration
func (p *Provider) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["APIKey"].(string); ok {
		p.apiKey = apiKey
	} else {
		return fmt.Errorf("API key is required for OpenAI provider")
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
	return "openai"
}

// AvailableModels returns the available models for this provider
func (p *Provider) AvailableModels() []string {
	return []string{
		"gpt-4-vision-preview",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}
}

// GenerateResponse generates a response from OpenAI
func (p *Provider) GenerateResponse(ctx context.Context, prompt string, params map[string]interface{}) (string, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// Create request
	request := ChatRequest{
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

	if stream, ok := params["stream"].(bool); ok {
		request.Stream = stream
	}

	// Make API request
	response, err := p.makeOpenAIRequest(ctx, request)
	if err != nil {
		return "", err
	}

	// Extract text from response
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in response")
}

// StreamResponse streams a response from OpenAI
func (p *Provider) StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// Create request
	request := ChatRequest{
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

	return p.streamOpenAIRequest(ctx, request)
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
	// Only vision models support multimodal
	return p.model == "gpt-4-vision-preview"
}

// GenerateMultimodalResponse generates a multimodal response from OpenAI
func (p *Provider) GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error) {
	// Check if the model supports multimodal
	if !p.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", p.model)
	}

	// Convert multimodal input to OpenAI format
	openAIMessage, err := p.convertToOpenAIMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to OpenAI format: %w", err)
	}

	// Create request with multimedia message
	messages := []json.RawMessage{
		mustMarshal(map[string]interface{}{
			"role":    openAIMessage.Role,
			"content": openAIMessage.Content,
		}),
	}

	// Check for system message in metadata
	if sysPrompt, ok := input.Metadata["system"].(string); ok {
		// Insert system message at the beginning
		sysMsg := mustMarshal(map[string]interface{}{
			"role":    "system",
			"content": sysPrompt,
		})
		messages = append([]json.RawMessage{sysMsg}, messages...)
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"model":    p.model,
		"messages": messages,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		requestBody["max_tokens"] = maxTokens
	} else if input.MaxTokens > 0 {
		requestBody["max_tokens"] = input.MaxTokens
	} else {
		requestBody["max_tokens"] = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		requestBody["temperature"] = temperature
	} else if input.Temperature > 0 {
		requestBody["temperature"] = input.Temperature
	} else {
		requestBody["temperature"] = 0.7 // Default
	}

	// Marshal the request body
	reqBody, err := json.Marshal(requestBody)
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert response to multimodal output
	output := multimodal.NewOutput()

	// Add content from response
	if len(response.Choices) > 0 {
		output.AddText(response.Choices[0].Message.Content)
	}

	// Add usage information
	output.Metadata = map[string]interface{}{
		"prompt_tokens": response.Usage.PromptTokens,
		"total_tokens":  response.Usage.TotalTokens,
		"used_tokens":   response.Usage.CompletionTokens,
	}

	return output, nil
}

// StreamMultimodalResponse streams a multimodal response from OpenAI
func (p *Provider) StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error) {
	// Check if the model supports multimodal
	if !p.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", p.model)
	}

	// Convert multimodal input to OpenAI format
	openAIMessage, err := p.convertToOpenAIMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to OpenAI format: %w", err)
	}

	// Create request with multimedia message
	messages := []json.RawMessage{
		mustMarshal(map[string]interface{}{
			"role":    openAIMessage.Role,
			"content": openAIMessage.Content,
		}),
	}

	// Check for system message in metadata
	if sysPrompt, ok := input.Metadata["system"].(string); ok {
		// Insert system message at the beginning
		sysMsg := mustMarshal(map[string]interface{}{
			"role":    "system",
			"content": sysPrompt,
		})
		messages = append([]json.RawMessage{sysMsg}, messages...)
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"model":    p.model,
		"messages": messages,
		"stream":   true,
	}

	// Add parameters
	if maxTokens, ok := params["max_tokens"].(int); ok {
		requestBody["max_tokens"] = maxTokens
	} else if input.MaxTokens > 0 {
		requestBody["max_tokens"] = input.MaxTokens
	} else {
		requestBody["max_tokens"] = 4096 // Default
	}

	if temperature, ok := params["temperature"].(float64); ok {
		requestBody["temperature"] = temperature
	} else if input.Temperature > 0 {
		requestBody["temperature"] = input.Temperature
	} else {
		requestBody["temperature"] = 0.7 // Default
	}

	// Marshal the request body
	reqBody, err := json.Marshal(requestBody)
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	// Create output channel
	resultCh := make(chan *multimodal.Chunk)

	// Start goroutine to handle request and response
	go func() {
		defer close(resultCh)

		// Send request
		resp, err := p.client.Do(req)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case resultCh <- &multimodal.Chunk{
				Content: multimodal.NewTextContent(fmt.Sprintf("Error: %v", err)),
				Error:   err,
				IsFinal: true,
			}:
				return
			}
		}
		defer resp.Body.Close()

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("OpenAI API returned status code %d: %s", resp.StatusCode, string(body))
			select {
			case <-ctx.Done():
				return
			case resultCh <- &multimodal.Chunk{
				Content: multimodal.NewTextContent(fmt.Sprintf("Error: %v", err)),
				Error:   err,
				IsFinal: true,
			}:
				return
			}
		}

		// Process the stream
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
				case resultCh <- &multimodal.Chunk{
					Content: multimodal.NewTextContent(fmt.Sprintf("Error reading stream: %v", err)),
					Error:   err,
					IsFinal: true,
				}:
					return
				}
			}

			// Skip empty lines
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			// Check for "data: " prefix
			if bytes.HasPrefix(line, []byte("data: ")) {
				line = bytes.TrimPrefix(line, []byte("data: "))

				// Check for "[DONE]" marker
				if string(line) == "[DONE]" {
					// Send final empty chunk to indicate completion
					select {
					case <-ctx.Done():
						return
					case resultCh <- &multimodal.Chunk{
						Content: multimodal.NewTextContent(""),
						IsFinal: true,
					}:
						// Successfully sent final chunk
					}
					break
				}

				// Parse the JSON chunk
				var streamResponse ChatStreamResponse
				if err := json.Unmarshal(line, &streamResponse); err != nil {
					// Skip unparseable chunks
					continue
				}

				// Check if there are choices
				if len(streamResponse.Choices) > 0 {
					// Extract the content
					content := streamResponse.Choices[0].Delta.Content
					if content == "" {
						continue
					}

					// Create and send chunk
					select {
					case <-ctx.Done():
						return
					case resultCh <- &multimodal.Chunk{
						Content: multimodal.NewTextContent(content),
						IsFinal: false,
					}:
						// Successfully sent chunk
					}
				}
			}
		}
	}()

	return resultCh, nil
}

// Helper methods

// convertToOpenAIMessage converts a multimodal input to OpenAI format
func (p *Provider) convertToOpenAIMessage(input *multimodal.Input) (*MultimediaMessage, error) {
	if len(input.Contents) == 0 {
		return nil, fmt.Errorf("input contains no content")
	}

	// Create a message with multimedia content
	message := &MultimediaMessage{
		Role:    "user",
		Content: []MultimediaContent{},
	}

	// Convert each content item
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text part
			message.Content = append(message.Content, MultimediaContent{
				Type: "text",
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Add image part
			imageContent := MultimediaContent{
				Type: "image",
				Image: &MultimediaImageDetail{
					Detail: "high", // Use high detail for best analysis
				},
			}

			// If we have a URL, use it
			if content.URI != "" {
				imageContent.Image.URL = content.URI
			} else if content.Data != nil {
				// Otherwise, use base64 data URL
				imageContent.Image.URL = content.ToDataURL()
			} else {
				return nil, fmt.Errorf("image has no URI or data")
			}

			message.Content = append(message.Content, imageContent)

		default:
			return nil, fmt.Errorf("unsupported content type: %s", content.Type)
		}
	}

	return message, nil
}

// makeOpenAIRequest makes a request to the OpenAI API
func (p *Provider) makeOpenAIRequest(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// streamOpenAIRequest streams a request to the OpenAI API
func (p *Provider) streamOpenAIRequest(ctx context.Context, request ChatRequest) (<-chan string, error) {
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Accept", "text/event-stream")

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
			var streamResponse ChatStreamResponse
			if err := json.Unmarshal(line, &streamResponse); err != nil {
				continue
			}

			// Extract content
			if len(streamResponse.Choices) > 0 {
				content := streamResponse.Choices[0].Delta.Content
				if content != "" {
					select {
					case <-ctx.Done():
						return
					case ch <- content:
						// Successfully sent content
					}
				}
			}
		}
	}()

	return ch, nil
}

// Helper functions

// mustMarshal marshals an object to JSON and panics on error
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(data)
}
