package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OllamaShim implements the LLM provider interface for Ollama
type OllamaShim struct {
	BaseURL   string
	Model     string
	APIKey    string // Optional, for authentication if required
	Client    *http.Client
	MaxTokens int
}

// ChatMessage represents a message in the chat history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the chat completion API
type ChatRequest struct {
	Model       string                 `json:"model"`
	Messages    []ChatMessage          `json:"messages"`
	Stream      bool                   `json:"stream,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ChatResponse represents a response from the chat completion API
type ChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// NewOllamaShim creates a new OllamaShim
func NewOllamaShim(baseURL, model, apiKey string, maxTokens int) *OllamaShim {
	// If no base URL is provided, use a default
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	// Ensure the base URL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Use a reasonable default for max tokens if not provided
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	return &OllamaShim{
		BaseURL:   baseURL,
		Model:     model,
		APIKey:    apiKey,
		Client:    &http.Client{},
		MaxTokens: maxTokens,
	}
}

// CompleteChatPrompt sends a chat completion request to Ollama
func (s *OllamaShim) CompleteChatPrompt(messages []ChatMessage) (string, error) {
	chatReq := ChatRequest{
		Model:       s.Model,
		Messages:    messages,
		Stream:      false,
		Temperature: 0.7,
		MaxTokens:   s.MaxTokens,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	url := fmt.Sprintf("%s/api/chat", s.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.APIKey)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return chatResp.Message.Content, nil
}

// ParseSentinelfile uses the LLM to parse a Sentinelfile into a structured format
func (s *OllamaShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
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

	// Extract the JSON part (in case the model outputs anything else)
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
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start >= 0 && end > start {
		return s[start : end+1]
	}

	return s // Return original if no JSON-like structure found
}
