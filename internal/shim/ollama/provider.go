package ollama

import (
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
)

// Provider implements the LLM provider interface for Ollama
type Provider struct {
	baseURL          string
	model            string
	client           *http.Client
	multimodalModels []string
}

// MultimediaMessage represents a message with multimedia content
type MultimediaMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

// ContentPart represents a part of the message content
type ContentPart struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Image *struct {
		Data string `json:"data"`
		Type string `json:"type,omitempty"`
	} `json:"image,omitempty"`
}

// StreamingResponse represents a streaming response from Ollama
type StreamingResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// MultimodalRequest represents a request to Ollama with multimodal content
type MultimodalRequest struct {
	Model       string              `json:"model"`
	Messages    []MultimediaMessage `json:"messages"`
	Stream      bool                `json:"stream"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

// NewProvider creates a new Ollama provider
func NewProvider() interface{} {
	return &Provider{
		baseURL: "http://localhost:11434",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		// Models known to support multimodal in Ollama
		multimodalModels: []string{
			"llava",
			"llava:7b",
			"llava:13b",
			"bakllava",
			"moondream",
		},
	}
}

// Name returns the name of the provider
func (p *Provider) Name() string {
	return "ollama"
}

// AvailableModels returns the available models
func (p *Provider) AvailableModels() []string {
	// TODO: Call Ollama API to get available models
	return []string{
		"llama3",
		"mistral",
		"llava",
		"bakllava",
		"moondream",
	}
}

// GenerateResponse generates a response from the Ollama LLM
func (p *Provider) GenerateResponse(ctx context.Context, prompt string, params map[string]interface{}) (string, error) {
	// Get model from params
	model, ok := params["model"].(string)
	if !ok || model == "" {
		model = "llama3" // Default model
	}
	p.model = model

	// Extract max tokens and temperature from params
	maxTokens := 4096
	if mt, ok := params["max_tokens"].(int); ok && mt > 0 {
		maxTokens = mt
	}

	temperature := 0.7
	if temp, ok := params["temperature"].(float64); ok && temp >= 0 {
		temperature = temp
	}

	// Create chat messages
	messages := []ChatMessage{
		{Role: "user", Content: prompt},
	}

	// Create chat request
	chatReq := ChatRequest{
		Model:       model,
		Messages:    messages,
		Stream:      false,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	// Send request
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	// Create request URL
	url := fmt.Sprintf("%s/api/chat", p.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := p.client.Do(req)
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

// StreamResponse streams a response from the Ollama LLM
func (p *Provider) StreamResponse(ctx context.Context, prompt string, params map[string]interface{}) (<-chan string, error) {
	// Get model from params
	model, ok := params["model"].(string)
	if !ok || model == "" {
		model = "llama3" // Default model
	}
	p.model = model

	// Extract max tokens and temperature from params
	maxTokens := 4096
	if mt, ok := params["max_tokens"].(int); ok && mt > 0 {
		maxTokens = mt
	}

	temperature := 0.7
	if temp, ok := params["temperature"].(float64); ok && temp >= 0 {
		temperature = temp
	}

	// Create chat messages
	messages := []ChatMessage{
		{Role: "user", Content: prompt},
	}

	// Create chat request with streaming enabled
	chatReq := ChatRequest{
		Model:       model,
		Messages:    messages,
		Stream:      true,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	// Send request
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat request: %w", err)
	}

	// Create request URL
	url := fmt.Sprintf("%s/api/chat", p.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create output channel
	ch := make(chan string)

	// Start goroutine to handle streaming
	go func() {
		defer close(ch)

		// Send request
		resp, err := p.client.Do(req)
		if err != nil {
			ch <- fmt.Sprintf("Error: %v", err)
			return
		}
		defer resp.Body.Close()

		// Check for HTTP errors
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			ch <- fmt.Sprintf("Error: API returned status code %d: %s", resp.StatusCode, string(body))
			return
		}

		// Process the streaming response - Ollama returns each chunk as a complete JSON object
		scanner := json.NewDecoder(resp.Body)
		for {
			var streamResp StreamingResponse
			if err := scanner.Decode(&streamResp); err != nil {
				if err != io.EOF {
					select {
					case <-ctx.Done():
						return
					case ch <- fmt.Sprintf("Error: %v", err):
						return
					}
				}
				break
			}

			// Send the response chunk
			select {
			case <-ctx.Done():
				return
			case ch <- streamResp.Response:
				// Check if this is the last chunk
				if streamResp.Done {
					return
				}
			}
		}
	}()

	return ch, nil
}

// GetEmbeddings gets embeddings for the given texts
func (p *Provider) GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// Simple implementation for embeddings (most Ollama models don't have a separate embeddings endpoint)
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		// Return dummy embeddings for now - in a real implementation,
		// we would call the Ollama API to get actual embeddings
		embeddings[i] = []float32{0.1, 0.2, 0.3, 0.4, 0.5}
	}
	return embeddings, nil
}

// SupportsMultimodal checks if the provider supports multimodal inputs and outputs
func (p *Provider) SupportsMultimodal() bool {
	// Check if the current model supports multimodal
	for _, model := range p.multimodalModels {
		if strings.EqualFold(p.model, model) || strings.HasPrefix(strings.ToLower(p.model), model+":") {
			return true
		}
	}
	return false
}

// GenerateMultimodalResponse generates a multimodal response from the Ollama LLM
func (p *Provider) GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error) {
	// Check if the model supports multimodal
	if !p.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", p.model)
	}

	// Get model from params
	model, ok := params["model"].(string)
	if !ok || model == "" {
		model = p.model
	}

	// Extract max tokens and temperature from params
	maxTokens := 4096
	if mt, ok := params["max_tokens"].(int); ok && mt > 0 {
		maxTokens = mt
	} else if input.MaxTokens > 0 {
		maxTokens = input.MaxTokens
	}

	temperature := 0.7
	if temp, ok := params["temperature"].(float64); ok && temp >= 0 {
		temperature = temp
	} else if input.Temperature > 0 {
		temperature = input.Temperature
	}

	// Convert multimodal input to Ollama format
	multimodalMsg, err := p.convertToOllamaMultimodalMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert multimodal input: %w", err)
	}

	// Create multimodal request
	req := MultimodalRequest{
		Model:       model,
		Messages:    []MultimediaMessage{multimodalMsg},
		Stream:      false,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	// Send request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal multimodal request: %w", err)
	}

	// Create request URL
	url := fmt.Sprintf("%s/api/chat", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned non-200 status code %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Create multimodal output
	output := multimodal.NewOutput()
	output.AddText(chatResp.Message.Content)

	return output, nil
}

// StreamMultimodalResponse streams a multimodal response from the Ollama LLM
func (p *Provider) StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error) {
	// Check if the model supports multimodal
	if !p.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", p.model)
	}

	// Get model from params
	model, ok := params["model"].(string)
	if !ok || model == "" {
		model = p.model
	}

	// Extract max tokens and temperature from params
	maxTokens := 4096
	if mt, ok := params["max_tokens"].(int); ok && mt > 0 {
		maxTokens = mt
	} else if input.MaxTokens > 0 {
		maxTokens = input.MaxTokens
	}

	temperature := 0.7
	if temp, ok := params["temperature"].(float64); ok && temp >= 0 {
		temperature = temp
	} else if input.Temperature > 0 {
		temperature = input.Temperature
	}

	// Convert multimodal input to Ollama format
	multimodalMsg, err := p.convertToOllamaMultimodalMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert multimodal input: %w", err)
	}

	// Create multimodal request with streaming enabled
	req := MultimodalRequest{
		Model:       model,
		Messages:    []MultimediaMessage{multimodalMsg},
		Stream:      true,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	// Send request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal multimodal request: %w", err)
	}

	// Create request URL
	url := fmt.Sprintf("%s/api/chat", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Create output channel for multimodal chunks
	resultCh := make(chan *multimodal.Chunk)

	// Start goroutine to handle streaming
	go func() {
		defer close(resultCh)

		// Send request
		resp, err := p.client.Do(httpReq)
		if err != nil {
			chunk := &multimodal.Chunk{
				Content: multimodal.NewTextContent(fmt.Sprintf("Error: %v", err)),
				IsFinal: true,
				Error:   err,
			}
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				return
			}
		}
		defer resp.Body.Close()

		// Check for HTTP errors
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("Ollama API returned non-200 status code %d: %s", resp.StatusCode, string(body))
			chunk := &multimodal.Chunk{
				Content: multimodal.NewTextContent(fmt.Sprintf("Error: %s", err)),
				IsFinal: true,
				Error:   err,
			}
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				return
			}
		}

		// Process the streaming response
		scanner := json.NewDecoder(resp.Body)
		for {
			var streamResp StreamingResponse
			if err := scanner.Decode(&streamResp); err != nil {
				if err != io.EOF {
					errChunk := &multimodal.Chunk{
						Content: multimodal.NewTextContent(fmt.Sprintf("Error: %v", err)),
						IsFinal: true,
						Error:   err,
					}
					select {
					case <-ctx.Done():
						return
					case resultCh <- errChunk:
						return
					}
				}
				break
			}

			// Create chunk with text content
			chunk := &multimodal.Chunk{
				Content: multimodal.NewTextContent(streamResp.Response),
				IsFinal: streamResp.Done,
			}

			// Send chunk
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				// Check if this is the last chunk
				if streamResp.Done {
					return
				}
			}
		}
	}()

	return resultCh, nil
}

// convertToOllamaMultimodalMessage converts a multimodal input to Ollama format
func (p *Provider) convertToOllamaMultimodalMessage(input *multimodal.Input) (MultimediaMessage, error) {
	if len(input.Contents) == 0 {
		return MultimediaMessage{}, fmt.Errorf("input contains no content")
	}

	// Create message with user role
	message := MultimediaMessage{
		Role:    "user",
		Content: []ContentPart{},
	}

	// Process each content item
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text part
			message.Content = append(message.Content, ContentPart{
				Type: "text",
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Add image part
			// Ollama expects image data as base64
			imgData := base64.StdEncoding.EncodeToString(content.Data)

			// Create image part with base64 data
			message.Content = append(message.Content, ContentPart{
				Type: "image",
				Image: &struct {
					Data string `json:"data"`
					Type string `json:"type,omitempty"`
				}{
					Data: imgData,
					Type: content.MimeType,
				},
			})

		default:
			return MultimediaMessage{}, fmt.Errorf("unsupported content type: %s", content.Type)
		}
	}

	return message, nil
}
