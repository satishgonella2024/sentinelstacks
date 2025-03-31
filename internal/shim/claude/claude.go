package claude

import (
	"bytes"
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
