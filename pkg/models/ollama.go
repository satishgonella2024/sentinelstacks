package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// OllamaAdapter provides an interface to the Ollama API
type OllamaAdapter struct {
	Endpoint string
	Model    string
}

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model    string  `json:"model"`
	Prompt   string  `json:"prompt"`
	Stream   bool    `json:"stream,omitempty"`
	Options  Options `json:"options,omitempty"`
	System   string  `json:"system,omitempty"`
	Template string  `json:"template,omitempty"`
}

// Options represents optional parameters for the Ollama API
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Model     string  `json:"model"`
	Response  string  `json:"response"`
	Done      bool    `json:"done"`
	Context   []int   `json:"context,omitempty"`
	TotalDuration int64   `json:"total_duration,omitempty"`
	LoadDuration   int64   `json:"load_duration,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
}

// NewOllamaAdapter creates a new adapter for Ollama
func NewOllamaAdapter(endpoint string, model string) *OllamaAdapter {
	return &OllamaAdapter{
		Endpoint: endpoint,
		Model:    model,
	}
}

// Generate sends a prompt to Ollama and returns the response
func (a *OllamaAdapter) Generate(prompt string, systemPrompt string, options Options) (string, error) {
	url := fmt.Sprintf("%s/api/generate", a.Endpoint)
	
	request := OllamaRequest{
		Model:    a.Model,
		Prompt:   prompt,
		System:   systemPrompt,
		Options:  options,
	}
	
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return ollamaResponse.Response, nil
}
