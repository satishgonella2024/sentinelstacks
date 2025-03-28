package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OllamaAdapter provides an interface to the Ollama API
type OllamaAdapter struct {
	Endpoint string
	Model    string
	Verbose  bool
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
	
	// Set stream to false to get a single response
	request := OllamaRequest{
		Model:    a.Model,
		Prompt:   prompt,
		System:   systemPrompt,
		Options:  options,
		Stream:   false, // Important: set to false to avoid streaming
	}
	
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	if a.Verbose {
		fmt.Printf("Sending request to %s with model %s...\n", a.Endpoint, a.Model)
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Print the raw response for debugging
	if a.Verbose {
		fmt.Printf("Raw response status: %s\n", resp.Status)
		fmt.Printf("Raw response body: %s\n", string(body))
	}
	
	// Try to parse the response
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		// If we can't parse as a single response, try to handle streaming format
		return a.handleStreamingResponse(body)
	}
	
	return ollamaResponse.Response, nil
}

// handleStreamingResponse handles streaming responses from Ollama
func (a *OllamaAdapter) handleStreamingResponse(responseBody []byte) (string, error) {
	if a.Verbose {
		fmt.Println("Handling streaming response format...")
	}
	
	// Split the response into lines
	lines := strings.Split(string(responseBody), "\n")
	
	// Accumulate the response parts
	var responseBuilder strings.Builder
	
	// Track if we've seen the done:true message
	seenDone := false
	
	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Parse each line as a separate JSON object
		var streamResponse OllamaResponse
		err := json.Unmarshal([]byte(line), &streamResponse)
		if err != nil {
			if a.Verbose {
				fmt.Printf("Error parsing line: %v\n", err)
			}
			continue // Skip this line if we can't parse it
		}
		
		// Accumulate the response parts
		responseBuilder.WriteString(streamResponse.Response)
		
		// Check if this is the final message
		if streamResponse.Done {
			seenDone = true
			// We could break here, but let's process all lines to be safe
		}
	}
	
	if !seenDone && a.Verbose {
		fmt.Println("Warning: Did not see 'done:true' in the streaming response")
	}
	
	finalResponse := responseBuilder.String()
	
	if a.Verbose {
		fmt.Printf("Reconstructed response:\n%s\n", finalResponse)
	}
	
	return finalResponse, nil
}

// GetCapabilities returns the capabilities of the model
func (a *OllamaAdapter) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		Streaming:       true,
		FunctionCalling: false, // Ollama doesn't support function calling yet
		MaxTokens:       4096,  // This varies by model, using a conservative estimate
		Multimodal:      false, // Some Ollama models support images, but we're not implementing that yet
	}
}
