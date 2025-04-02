package google

import (
	"bufio"
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

// Constants for Google AI API
const (
	DefaultEndpoint    = "https://generativelanguage.googleapis.com/v1beta/models"
	DefaultModel       = "gemini-1.5-pro"
	DefaultMaxTokens   = 4096
	DefaultTemperature = 0.7
)

// GoogleShim implements the LLM interface for Google's Gemini models
type GoogleShim struct {
	APIKey      string
	Model       string
	Client      *http.Client
	Endpoint    string
	MaxTokens   int
	Temperature float64
}

// ContentPart represents a part of the content in the Google API
type ContentPart struct {
	Text  string             `json:"text,omitempty"`
	InlineData *InlineData   `json:"inlineData,omitempty"`
}

// InlineData represents image data in the Google API
type InlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// ContentRequest represents a request to the Google Gemini API
type ContentRequest struct {
	Contents     []Content           `json:"contents"`
	SafetySettings []SafetySetting   `json:"safetySettings,omitempty"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
}

// Content represents a content message in the Google API
type Content struct {
	Role    string         `json:"role,omitempty"`
	Parts   []ContentPart  `json:"parts"`
}

// SafetySetting represents a safety setting in the Google API
type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GenerationConfig represents generation parameters in the Google API
type GenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
}

// ContentResponse represents a response from the Google Gemini API
type ContentResponse struct {
	Candidates []struct {
		Content struct {
			Role  string         `json:"role"`
			Parts []ContentPart  `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

// NewGoogleShim creates a new GoogleShim
func NewGoogleShim(apiKey, model string) *GoogleShim {
	if model == "" {
		model = DefaultModel
	}

	return &GoogleShim{
		APIKey:      apiKey,
		Model:       model,
		Client:      &http.Client{Timeout: 60 * time.Second},
		Endpoint:    DefaultEndpoint,
		MaxTokens:   DefaultMaxTokens,
		Temperature: DefaultTemperature,
	}
}

// GenerateContent sends a text prompt to Gemini and returns the response
func (s *GoogleShim) GenerateContent(ctx context.Context, prompt, systemPrompt string) (string, error) {
	// Create content message
	contents := []Content{
		{
			Parts: []ContentPart{
				{
					Text: prompt,
				},
			},
		},
	}

	// Add system prompt if provided
	if systemPrompt != "" {
		// Google doesn't have native system prompts, so we prepend it to the user message
		contents[0].Parts = append([]ContentPart{
			{
				Text: systemPrompt + "\n\n",
			},
		}, contents[0].Parts...)
	}

	// Create request
	request := ContentRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature:     s.Temperature,
			MaxOutputTokens: s.MaxTokens,
		},
	}

	// Make the API request
	response, err := s.makeContentRequest(ctx, request)
	if err != nil {
		return "", err
	}

	// Extract text from the response
	return s.extractTextFromResponse(response)
}

// MultimodalGenerate handles multimodal content generation
func (s *GoogleShim) MultimodalGenerate(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error) {
	// Check if model supports multimodal
	if !s.supportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal input", s.Model)
	}

	// Prepare the content request with multimodal content
	contents := []Content{
		{
			Parts: []ContentPart{},
		},
	}

	// Extract text and images from input
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text content
			contents[0].Parts = append(contents[0].Parts, ContentPart{
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Skip if no data
			if len(content.Data) == 0 {
				continue
			}

			// Add image content
			base64Data := base64.StdEncoding.EncodeToString(content.Data)
			contents[0].Parts = append(contents[0].Parts, ContentPart{
				InlineData: &InlineData{
					MimeType: content.MimeType,
					Data:     base64Data,
				},
			})
		}
	}

	// Add system prompt if present in metadata
	if systemPrompt, ok := input.Metadata["system"].(string); ok && systemPrompt != "" {
		// For Google, prepend to the first text part or add a new text part
		if len(contents[0].Parts) > 0 && contents[0].Parts[0].Text != "" {
			contents[0].Parts[0].Text = systemPrompt + "\n\n" + contents[0].Parts[0].Text
		} else {
			// Insert system prompt at the beginning
			contents[0].Parts = append([]ContentPart{
				{
					Text: systemPrompt + "\n\n",
				},
			}, contents[0].Parts...)
		}
	}

	// Create request
	request := ContentRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature: input.Temperature,
		},
	}

	// Set max tokens if specified
	if input.MaxTokens > 0 {
		request.GenerationConfig.MaxOutputTokens = input.MaxTokens
	} else {
		request.GenerationConfig.MaxOutputTokens = s.MaxTokens
	}

	// Make the API request
	response, err := s.makeContentRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// Extract text from the response
	responseText, err := s.extractTextFromResponse(response)
	if err != nil {
		return nil, err
	}

	// Create multimodal output
	output := multimodal.NewOutput()
	output.AddText(responseText)

	// Add metadata
	output.Metadata = map[string]interface{}{
		"model": s.Model,
	}

	return output, nil
}

// StreamContent streams a text completion from Gemini
func (s *GoogleShim) StreamContent(ctx context.Context, prompt, systemPrompt string) (<-chan string, error) {
	// Create the return channel
	resultCh := make(chan string)

	// Create content message
	contents := []Content{
		{
			Parts: []ContentPart{
				{
					Text: prompt,
				},
			},
		},
	}

	// Add system prompt if provided
	if systemPrompt != "" {
		// Google doesn't have native system prompts, so we prepend it to the user message
		contents[0].Parts = append([]ContentPart{
			{
				Text: systemPrompt + "\n\n",
			},
		}, contents[0].Parts...)
	}

	// Create request
	request := ContentRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature:     s.Temperature,
			MaxOutputTokens: s.MaxTokens,
		},
	}

	// Make the streaming request
	go func() {
		defer close(resultCh)

		// Make streaming request
		err := s.makeStreamingRequest(ctx, request, func(chunk string) {
			select {
			case <-ctx.Done():
				return
			case resultCh <- chunk:
				// Successfully sent chunk
			}
		})

		if err != nil {
			select {
			case <-ctx.Done():
				return
			case resultCh <- fmt.Sprintf("Error: %v", err):
				// Sent error message
			}
		}
	}()

	return resultCh, nil
}

// StreamMultimodalContent streams a multimodal completion from Gemini
func (s *GoogleShim) StreamMultimodalContent(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error) {
	// Check if model supports multimodal
	if !s.supportsMultimodal() {
		return nil, fmt.Errorf("model %s does not support multimodal input", s.Model)
	}

	// Create the return channel
	resultCh := make(chan *multimodal.Chunk)

	// Prepare the content request with multimodal content
	contents := []Content{
		{
			Parts: []ContentPart{},
		},
	}

	// Extract text and images from input
	for _, content := range input.Contents {
		switch content.Type {
		case multimodal.MediaTypeText:
			// Add text content
			contents[0].Parts = append(contents[0].Parts, ContentPart{
				Text: content.Text,
			})

		case multimodal.MediaTypeImage:
			// Skip if no data
			if len(content.Data) == 0 {
				continue
			}

			// Add image content
			base64Data := base64.StdEncoding.EncodeToString(content.Data)
			contents[0].Parts = append(contents[0].Parts, ContentPart{
				InlineData: &InlineData{
					MimeType: content.MimeType,
					Data:     base64Data,
				},
			})
		}
	}

	// Add system prompt if present in metadata
	if systemPrompt, ok := input.Metadata["system"].(string); ok && systemPrompt != "" {
		// For Google, prepend to the first text part or add a new text part
		if len(contents[0].Parts) > 0 && contents[0].Parts[0].Text != "" {
			contents[0].Parts[0].Text = systemPrompt + "\n\n" + contents[0].Parts[0].Text
		} else {
			// Insert system prompt at the beginning
			contents[0].Parts = append([]ContentPart{
				{
					Text: systemPrompt + "\n\n",
				},
			}, contents[0].Parts...)
		}
	}

	// Create request
	request := ContentRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature: input.Temperature,
		},
	}

	// Set max tokens if specified
	if input.MaxTokens > 0 {
		request.GenerationConfig.MaxOutputTokens = input.MaxTokens
	} else {
		request.GenerationConfig.MaxOutputTokens = s.MaxTokens
	}

	// Make the streaming request
	go func() {
		defer close(resultCh)

		var isFinal bool
		var isError bool

		// Make streaming request
		err := s.makeStreamingRequest(ctx, request, func(chunk string) {
			// Convert to multimodal chunk
			multiChunk := multimodal.NewChunk(multimodal.NewTextContent(chunk), isFinal)
			
			// Send to channel
			select {
			case <-ctx.Done():
				return
			case resultCh <- multiChunk:
				// Successfully sent chunk
			}
		})

		if err != nil {
			// Create error chunk
			errorChunk := multimodal.NewChunk(
				multimodal.NewTextContent(fmt.Sprintf("Error: %v", err)),
				true,
			)
			errorChunk.Error = err
			isError = true
			
			select {
			case <-ctx.Done():
				return
			case resultCh <- errorChunk:
				// Sent error chunk
			}
		}

		// If no error occurred, send a final empty chunk
		if !isError {
			finalChunk := multimodal.NewChunk(multimodal.NewTextContent(""), true)
			
			select {
			case <-ctx.Done():
				return
			case resultCh <- finalChunk:
				// Sent final chunk
			}
		}
	}()

	return resultCh, nil
}

// makeContentRequest makes a request to the Google Gemini API for content generation
func (s *GoogleShim) makeContentRequest(ctx context.Context, request ContentRequest) (*ContentResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create URL for the model
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", s.Endpoint, s.Model, s.APIKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response ContentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// makeStreamingRequest makes a streaming request to the Google Gemini API
func (s *GoogleShim) makeStreamingRequest(ctx context.Context, request ContentRequest, chunkHandler func(string)) error {
	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create URL for the model
	url := fmt.Sprintf("%s/%s:streamGenerateContent?key=%s", s.Endpoint, s.Model, s.APIKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Google API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Read the streaming response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse the JSON chunk
		var response ContentResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue // Skip unparseable lines
		}

		// Extract text from the response
		if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
			for _, part := range response.Candidates[0].Content.Parts {
				if part.Text != "" {
					chunkHandler(part.Text)
				}
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

// extractTextFromResponse extracts text from a ContentResponse
func (s *GoogleShim) extractTextFromResponse(response *ContentResponse) (string, error) {
	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	var text string
	for _, part := range response.Candidates[0].Content.Parts {
		text += part.Text
	}

	return text, nil
}

// supportsMultimodal checks if the model supports multimodal input
func (s *GoogleShim) supportsMultimodal() bool {
	multimodalModels := map[string]bool{
		"gemini-pro-vision":  true,
		"gemini-1.5-pro":     true,
		"gemini-1.5-flash":   true,
	}

	return multimodalModels[s.Model]
}
