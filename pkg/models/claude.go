package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ClaudeAdapter provides an interface to the Anthropic Claude API
type ClaudeAdapter struct {
	APIKey  string
	Model   string
	Verbose bool
}

// ClaudeRequest represents a request to the Claude API
type ClaudeRequest struct {
	Model       string         `json:"model"`
	Messages    []ClaudeMessage `json:"messages"`
	MaxTokens   int            `json:"max_tokens,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	TopP        float64        `json:"top_p,omitempty"`
	System      string         `json:"system,omitempty"`
}

// ClaudeMessage represents a message in the Claude chat format
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents a response from the Claude API
type ClaudeResponse struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Role    string        `json:"role"`
	Content []ClaudeContent `json:"content"`
	Model   string        `json:"model"`
	Usage   ClaudeUsage     `json:"usage"`
}

// ClaudeContent represents a content block in the Claude response
type ClaudeContent struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
}

// ClaudeUsage represents token usage information
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// NewClaudeAdapter creates a new adapter for Claude
func NewClaudeAdapter(apiKey string, model string) *ClaudeAdapter {
	return &ClaudeAdapter{
		APIKey: apiKey,
		Model:  model,
	}
}

// Generate sends a prompt to Claude and returns the response
func (a *ClaudeAdapter) Generate(prompt string, systemPrompt string, options Options) (string, error) {
	url := "https://api.anthropic.com/v1/messages"
	
	// Create messages array
	messages := []ClaudeMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}
	
	// Create request
	request := ClaudeRequest{
		Model:       a.Model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		TopP:        options.TopP,
		System:      systemPrompt,
	}
	
	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	if a.Verbose {
		fmt.Printf("Sending request to Claude API with model %s...\n", a.Model)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Claude: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse response
	var claudeResponse ClaudeResponse
	err = json.Unmarshal(body, &claudeResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Claude response: %w", err)
	}
	
	// Extract text from content blocks
	var responseText string
	for _, content := range claudeResponse.Content {
		if content.Type == "text" {
			responseText += content.Text
		}
	}
	
	return responseText, nil
}

// GetCapabilities returns the capabilities of the model
func (a *ClaudeAdapter) GetCapabilities() ModelCapabilities {
	// Define capabilities based on the model
	caps := ModelCapabilities{
		Streaming:       true,
		FunctionCalling: false,
		MaxTokens:       100000, // Claude has very large context windows
		Multimodal:      false,
	}
	
	// Set model-specific capabilities
	switch a.Model {
	case "claude-3-opus-20240229":
		caps.MaxTokens = 200000
		caps.Multimodal = true
	case "claude-3-sonnet-20240229":
		caps.MaxTokens = 180000
		caps.Multimodal = true
	case "claude-3-haiku-20240307":
		caps.MaxTokens = 150000
		caps.Multimodal = true
	}
	
	return caps
}
