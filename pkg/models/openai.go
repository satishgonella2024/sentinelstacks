package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIAdapter provides an interface to the OpenAI API
type OpenAIAdapter struct {
	APIKey  string
	Model   string
	Verbose bool
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// OpenAIMessage represents a message in the OpenAI chat format
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice in the OpenAI response
type Choice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewOpenAIAdapter creates a new adapter for OpenAI
func NewOpenAIAdapter(apiKey string, model string) *OpenAIAdapter {
	return &OpenAIAdapter{
		APIKey: apiKey,
		Model:  model,
	}
}

// Generate sends a prompt to OpenAI and returns the response
func (a *OpenAIAdapter) Generate(prompt string, systemPrompt string, options Options) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	
	// Create messages array
	messages := []OpenAIMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	
	// Create request
	request := OpenAIRequest{
		Model:       a.Model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		TopP:        options.TopP,
	}
	
	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	if a.Verbose {
		fmt.Printf("Sending request to OpenAI API with model %s...\n", a.Model)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.APIKey)
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse response
	var openaiResponse OpenAIResponse
	err = json.Unmarshal(body, &openaiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response: %w", err)
	}
	
	// Check if we have any choices
	if len(openaiResponse.Choices) == 0 {
		return "", fmt.Errorf("OpenAI returned no choices")
	}
	
	return openaiResponse.Choices[0].Message.Content, nil
}

// GetCapabilities returns the capabilities of the model
func (a *OpenAIAdapter) GetCapabilities() ModelCapabilities {
	// Define capabilities based on the model
	caps := ModelCapabilities{
		Streaming: true,
		MaxTokens: 4096, // Default
		Multimodal: false,
	}
	
	// Set model-specific capabilities
	switch a.Model {
	case "gpt-4", "gpt-4-turbo":
		caps.MaxTokens = 8192
		caps.FunctionCalling = true
	case "gpt-4-vision":
		caps.Multimodal = true
		caps.MaxTokens = 8192
		caps.FunctionCalling = true
	case "gpt-3.5-turbo":
		caps.MaxTokens = 4096
		caps.FunctionCalling = true
	}
	
	return caps
}
