package shim

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
)

// OllamaShim is an implementation of LLMShim for Ollama models
type OllamaShim struct {
	config       Config
	systemPrompt string
	httpClient   *http.Client
}

// OllamaGenerateRequest represents a request to the Ollama generate API
type OllamaGenerateRequest struct {
	Model       string                 `json:"model"`
	Prompt      string                 `json:"prompt"`
	System      string                 `json:"system,omitempty"`
	Template    string                 `json:"template,omitempty"`
	Context     []int                  `json:"context,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Raw         bool                   `json:"raw,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	KeepAlive   string                 `json:"keep_alive,omitempty"`
}

// OllamaGenerateResponse represents a response from the Ollama generate API
type OllamaGenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
	Context   []int     `json:"context,omitempty"`
	PromptEval    int   `json:"prompt_eval_count,omitempty"`
	Eval          int   `json:"eval_count,omitempty"`
	PromptTokens  int   `json:"prompt_tokens,omitempty"`
	Completion    int   `json:"completion_tokens,omitempty"`
	TotalDuration int64 `json:"total_duration,omitempty"`
	LoadDuration  int64 `json:"load_duration,omitempty"`
}

// NewOllamaShim creates a new shim for Ollama
func NewOllamaShim(config Config) *OllamaShim {
	// Set default timeout if not specified
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: timeout,
	}
	
	return &OllamaShim{
		config:     config,
		httpClient: httpClient,
	}
}

// Completion generates a text completion using Ollama
func (s *OllamaShim) Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.CompletionWithContext(ctx, prompt, maxTokens, temperature)
}

// CompletionWithContext generates a text completion using Ollama with context
func (s *OllamaShim) CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	// Configure the request payload
	request := OllamaGenerateRequest{
		Model:  s.config.Model,
		Prompt: prompt,
		System: s.systemPrompt,
		Stream: false, // Not streaming for regular completion
		Options: map[string]interface{}{
			"temperature": temperature,
		},
		KeepAlive: "5m", // Keep model loaded for 5 minutes
	}
	
	// Add max tokens if specified
	if maxTokens > 0 {
		request.Options["num_predict"] = maxTokens
	}
	
	// Marshal the request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create the HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.config.Endpoint,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Send the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	// For non-streaming mode, we need to handle the Ollama response format
	// which might return multiple JSON objects
	var fullResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		// Parse each line as a separate JSON response
		var response OllamaGenerateResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			return "", fmt.Errorf("failed to parse response: %w - response: %s", err, line)
		}
		
		// Accumulate the response text
		fullResponse += response.Response
		
		// If this is the final response, we're done
		if response.Done {
			break
		}
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}
	
	return fullResponse, nil
}

// MultimodalCompletion generates a multimodal completion using Ollama
func (s *OllamaShim) MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return s.MultimodalCompletionWithContext(ctx, input)
}

// MultimodalCompletionWithContext generates a multimodal completion using Ollama with context
func (s *OllamaShim) MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	// Check if model supports multimodal
	if !s.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", s.config.Model)
	}
	
	// For multimodal models like llava, we need to encode images and include them in the prompt
	// Extract text content from inputs
	var textContent string
	var imageContent []byte
	// Separate text and image content
	for _, content := range input.Contents {
		if content.Type == multimodal.MediaTypeText {
			textContent += content.Text + " "
		} else if content.Type == multimodal.MediaTypeImage {
			// For now, we only support a single image
			imageContent = content.Data
		}
	}
	
	// Prepare the request
	request := OllamaGenerateRequest{
		Model:  s.config.Model,
		Prompt: textContent,
		System: s.systemPrompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": input.Temperature,
		},
		KeepAlive: "5m",
	}
	
	// Add max tokens if specified
	if input.MaxTokens > 0 {
		request.Options["num_predict"] = input.MaxTokens
	}
	
	// If we have image content, we need to encode it as base64 and include it
	// However, the exact implementation might depend on the specific Ollama version and model
	// This is a simplified version and might need adjustments based on the Ollama API specs
	if len(imageContent) > 0 {
		// For now, this is a placeholder. In a complete implementation,
		// you would need to handle image encoding according to Ollama's multimodal API
		return nil, fmt.Errorf("multimodal support for Ollama is not fully implemented")
	}
	
	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.config.Endpoint,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	// For non-streaming mode, we need to handle the Ollama response format
	// which might return multiple JSON objects
	var fullResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		// Parse each line as a separate JSON response
		var response OllamaGenerateResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w - response: %s", err, line)
		}
		
		// Accumulate the response text
		fullResponse += response.Response
		
		// If this is the final response, we're done
		if response.Done {
			break
		}
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	
	// Create multimodal output with the response text
	output := multimodal.NewOutput()
	output.AddText(fullResponse)
	
	return output, nil
}

// StreamCompletion streams a text completion from Ollama
func (s *OllamaShim) StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error) {
	resultCh := make(chan string)
	
	// Configure the request payload
	request := OllamaGenerateRequest{
		Model:  s.config.Model,
		Prompt: prompt,
		System: s.systemPrompt,
		Stream: true, // Enable streaming
		Options: map[string]interface{}{
			"temperature": temperature,
		},
		KeepAlive: "5m", // Keep model loaded for 5 minutes
	}
	
	// Add max tokens if specified
	if maxTokens > 0 {
		request.Options["num_predict"] = maxTokens
	}
	
	// Marshal the request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.config.Endpoint,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		close(resultCh)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		close(resultCh)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	// Process streaming response in a goroutine
	go func() {
		defer resp.Body.Close()
		defer close(resultCh)
		
		// Create a new scanner to read the response line by line
		scanner := bufio.NewScanner(resp.Body)
		
		// Scan each line (JSON object) from the stream
		for scanner.Scan() {
			// Get the current line
			line := scanner.Text()
			if line == "" {
				continue
			}
			
			// Parse the line as JSON
			var response OllamaGenerateResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				// Skip malformed JSON
				continue
			}
			
			// Send the response text chunk to the channel
			select {
			case <-ctx.Done():
				// Context canceled, stop processing
				return
			case resultCh <- response.Response:
				// Successfully sent response chunk
			}
			
			// If this is the final response (done=true), exit
			if response.Done {
				break
			}
		}
		
		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			// Log scanner error (can't return it)
			fmt.Printf("Error reading stream: %v\n", err)
		}
	}()
	
	return resultCh, nil
}

// StreamMultimodalCompletion streams a multimodal completion from Ollama
func (s *OllamaShim) StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// Check if model supports multimodal
	if !s.SupportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal inputs", s.config.Model)
	}
	
	resultCh := make(chan *multimodal.Chunk)
	
	// Extract text content from inputs
	var textContent string
	var imageContent []byte
	
	// Separate text and image content
	for _, content := range input.Contents {
		if content.Type == multimodal.MediaTypeText {
			textContent += content.Text + " "
		} else if content.Type == multimodal.MediaTypeImage {
			// For now, we only support a single image
			imageContent = content.Data
		}
	}
	
	// Prepare the request
	request := OllamaGenerateRequest{
		Model:  s.config.Model,
		Prompt: textContent,
		System: s.systemPrompt,
		Stream: true, // Enable streaming
		Options: map[string]interface{}{
			"temperature": input.Temperature,
		},
		KeepAlive: "5m",
	}
	
	// Add max tokens if specified
	if input.MaxTokens > 0 {
		request.Options["num_predict"] = input.MaxTokens
	}
	
	// If we have image content, we need to encode it as base64 and include it
	// This is a placeholder for multimodal API implementation
	if len(imageContent) > 0 {
		// For now, mark that real multimodal handling will be needed
		// In a full implementation, you would encode the image for Ollama's API
		close(resultCh)
		return nil, fmt.Errorf("multimodal streaming support for Ollama is not fully implemented")
	}
	
	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		close(resultCh)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.config.Endpoint,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		close(resultCh)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		close(resultCh)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		close(resultCh)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	// Process the streaming response in a goroutine
	go func() {
		defer resp.Body.Close()
		defer close(resultCh)
		
		// Create a scanner to read the response line by line
		scanner := bufio.NewScanner(resp.Body)
		
		// Scan each line (JSON object) from the stream
		for scanner.Scan() {
			// Get the current line
			line := scanner.Text()
			if line == "" {
				continue
			}
			
			// Parse the JSON response
			var response OllamaGenerateResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				// Skip malformed JSON
				continue
			}
			
			// Create a multimodal chunk with the response text
			content := multimodal.NewTextContent(response.Response)
			chunk := multimodal.NewChunk(content, response.Done)
			
			// Send the chunk to the channel
			select {
			case <-ctx.Done():
				// Context canceled, stop processing
				return
			case resultCh <- chunk:
				// Successfully sent chunk
			}
			
			// If this is the final response, exit
			if response.Done {
				break
			}
		}
		
		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			// Log scanner error (can't return it from goroutine)
			fmt.Printf("Error reading stream: %v\n", err)
		}
	}()
	
	return resultCh, nil
}

// SetSystemPrompt sets the system prompt for Ollama
func (s *OllamaShim) SetSystemPrompt(prompt string) {
	s.systemPrompt = prompt
}

// ParseSentinelfile parses a Sentinelfile into a configuration map
func (s *OllamaShim) ParseSentinelfile(content string) (map[string]interface{}, error) {
	// Construct a request to parse the Sentinelfile using Ollama
	request := OllamaGenerateRequest{
		Model:  s.config.Model,
		Prompt: fmt.Sprintf("Parse this Sentinelfile and extract the key information into a JSON object:\n\n%s", content),
		System: "You are an AI assistant that specializes in parsing Sentinelfiles. Extract the name, description, baseModel, and capabilities from the Sentinelfile and return them as a valid JSON object.",
		Format: "json", // Request JSON output
		Options: map[string]interface{}{
			"temperature": 0.1, // Low temperature for more deterministic output
		},
	}
	
	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.config.Endpoint,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	// For non-streaming mode, we need to handle the Ollama response format
	// which might return multiple JSON objects
	var fullResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		// Parse each line as a separate JSON response
		var response OllamaGenerateResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w - response: %s", err, line)
		}
		
		// Accumulate the response text
		fullResponse += response.Response
		
		// If this is the final response, we're done
		if response.Done {
			break
		}
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	
	// Parse the JSON response text into a map
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(fullResponse), &result); err != nil {
		// If we can't parse the response as JSON, return a simple parsed structure
		return map[string]interface{}{
			"name":        "Unknown",
			"description": "Sentinelfile parsed by Ollama",
			"baseModel":   s.config.Model,
			"rawContent":  content,
		}, nil
	}
	
	// Ensure baseModel is set
	if _, ok := result["baseModel"]; !ok {
		result["baseModel"] = s.config.Model
	}
	
	return result, nil
}

// SupportsMultimodal returns whether Ollama supports multimodal inputs
func (s *OllamaShim) SupportsMultimodal() bool {
	// Ollama models that support multimodal inputs
	multimodalModels := map[string]bool{
		"llava": true,
		"llava:7b": true,
		"llava:13b": true,
		"llava:latest": true,
		"bakllava": true,
		"bakllava:latest": true,
		"moondream": true,
		"moondream:latest": true,
		"fuyu": true,
		"fuyu:latest": true,
		"yi-vl": true,
		"yi-vl:latest": true,
		"cogvlm": true,
		"cogvlm:latest": true,
	}
	
	// Check if the model is in the list of supported multimodal models
	// Also check for partial matches (e.g., if model starts with a supported model name)
	for modelName := range multimodalModels {
		if s.config.Model == modelName || strings.HasPrefix(s.config.Model, modelName+":") {
			return true
		}
	}
	
	return false
}

// Close cleans up any resources used by the Ollama shim
func (s *OllamaShim) Close() error {
	// Cancel any pending requests by closing the HTTP client's transport
	if s.httpClient != nil && s.httpClient.Transport != nil {
		if transport, ok := s.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
	
	return nil
}
