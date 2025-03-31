package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// Provider is the Claude provider implementation
type Provider struct {
	client *Client
	// Add any provider-specific fields here
}

// Client is a simplified client for the Claude API
type Client struct {
	apiKey   string
	endpoint string
	// Add any client-specific fields here
}

// NewClaudeShim creates a new ClaudeShim
func NewClaudeShim(apiKey, model string) *ClaudeShim {
	if model == "" {
		model = "claude-3-7-sonnet-20240307"
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
		client: &Client{},
	}
}

// Initialize initializes the provider with the given configuration
func (p *Provider) Initialize(config map[string]interface{}) error {
	if apiKey, ok := config["APIKey"].(string); ok {
		p.client.apiKey = apiKey
	}
	if endpoint, ok := config["Endpoint"].(string); ok {
		p.client.endpoint = endpoint
	}
	// Additional initialization can be done here
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
	// Simple implementation for now
	return fmt.Sprintf("Claude response to: %s", prompt), nil
}

// StreamResponse streams a response from Claude
func (p *Provider) StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error) {
	// Simple implementation for now
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- fmt.Sprintf("Claude response (streaming) to: %s", prompt)
	}()
	return ch, nil
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
func (p *Provider) GenerateMultimodalResponse(ctx context.Context, input interface{}) (interface{}, error) {
	// This is a stub implementation that will be completed once we have
	// a proper interface for multimodal types without circular imports
	return map[string]interface{}{
		"text": "Claude multimodal response (stub implementation)",
	}, nil
}

// StreamMultimodalResponse streams a multimodal response from Claude
func (p *Provider) StreamMultimodalResponse(ctx context.Context, input interface{}) (<-chan interface{}, error) {
	// This is a stub implementation that will be completed once we have
	// a proper interface for multimodal types without circular imports
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		ch <- map[string]interface{}{
			"text":     "Claude multimodal response (streaming stub)",
			"is_final": true,
		}
	}()
	return ch, nil
}
