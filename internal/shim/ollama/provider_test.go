package ollama

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

func TestNewProvider(t *testing.T) {
	provider := NewProvider().(*Provider)

	if provider.baseURL != "http://localhost:11434" {
		t.Errorf("Expected default baseURL to be http://localhost:11434, got %s", provider.baseURL)
	}

	if len(provider.multimodalModels) == 0 {
		t.Error("Expected multimodalModels to be populated")
	}
}

func TestSupportsMultimodal(t *testing.T) {
	provider := NewProvider().(*Provider)

	// Test with non-multimodal model
	provider.model = "llama3"
	if provider.SupportsMultimodal() {
		t.Errorf("Expected llama3 to not support multimodal")
	}

	// Test with multimodal model
	provider.model = "llava"
	if !provider.SupportsMultimodal() {
		t.Errorf("Expected llava to support multimodal")
	}

	// Test with versioned multimodal model
	provider.model = "llava:7b"
	if !provider.SupportsMultimodal() {
		t.Errorf("Expected llava:7b to support multimodal")
	}
}

func TestConvertToOllamaMultimodalMessage(t *testing.T) {
	provider := NewProvider().(*Provider)

	// Create a multimodal input with text
	input := multimodal.NewInput()
	input.AddText("Describe this image:")

	// Add dummy image data
	dummyImageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header
	input.AddImage(dummyImageData, "image/jpeg")

	// Convert to Ollama format
	message, err := provider.convertToOllamaMultimodalMessage(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check message role
	if message.Role != "user" {
		t.Errorf("Expected role to be user, got %s", message.Role)
	}

	// Check content parts
	if len(message.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(message.Content))
	}

	// Check text part
	if message.Content[0].Type != "text" {
		t.Errorf("Expected first part type to be text, got %s", message.Content[0].Type)
	}
	if message.Content[0].Text != "Describe this image:" {
		t.Errorf("Expected text content to be 'Describe this image:', got '%s'", message.Content[0].Text)
	}

	// Check image part
	if message.Content[1].Type != "image" {
		t.Errorf("Expected second part type to be image, got %s", message.Content[1].Type)
	}
	if message.Content[1].Image == nil {
		t.Fatalf("Expected image part to be populated")
	}
	if message.Content[1].Image.Type != "image/jpeg" {
		t.Errorf("Expected image type to be image/jpeg, got %s", message.Content[1].Image.Type)
	}
	if message.Content[1].Image.Data == "" {
		t.Errorf("Expected image data to be populated")
	}
}

func TestGenerateMultimodalResponse(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/api/chat" {
			t.Errorf("Expected path to be /api/chat, got %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"model": "llava",
			"message": {
				"role": "assistant",
				"content": "This is a test response for the image."
			},
			"done": true
		}`))
	}))
	defer server.Close()

	// Create provider with test server URL
	provider := NewProvider().(*Provider)
	provider.baseURL = server.URL
	provider.model = "llava"

	// Create multimodal input
	input := multimodal.NewInput()
	input.AddText("Describe this image:")

	// Add dummy image data
	dummyImageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header
	input.AddImage(dummyImageData, "image/jpeg")

	// Generate response
	output, err := provider.GenerateMultimodalResponse(context.Background(), input, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check output
	if output == nil {
		t.Fatalf("Expected output to be populated")
	}

	// Get text from output
	text := ""
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			text = content.Text
			break
		}
	}

	if text != "This is a test response for the image." {
		t.Errorf("Expected text to be 'This is a test response for the image.', got '%s'", text)
	}
}

func TestStreamMultimodalResponse(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/api/chat" {
			t.Errorf("Expected path to be /api/chat, got %s", r.URL.Path)
		}

		// Get query param for streaming
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			t.Errorf("Expected Content-Type to contain application/json")
		}

		// Set headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write streaming chunks
		chunks := []string{
			`{"model":"llava","response":"This ","done":false}`,
			`{"model":"llava","response":"is ","done":false}`,
			`{"model":"llava","response":"a ","done":false}`,
			`{"model":"llava","response":"test ","done":false}`,
			`{"model":"llava","response":"response.","done":true}`,
		}

		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n"))
			w.(http.Flusher).Flush()
		}
	}))
	defer server.Close()

	// Create provider with test server URL
	provider := NewProvider().(*Provider)
	provider.baseURL = server.URL
	provider.model = "llava"

	// Create multimodal input
	input := multimodal.NewInput()
	input.AddText("Describe this image:")

	// Add dummy image data
	dummyImageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header
	input.AddImage(dummyImageData, "image/jpeg")

	// Stream response
	chunkCh, err := provider.StreamMultimodalResponse(context.Background(), input, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Collect chunks
	var chunks []string
	var finalChunk bool
	for chunk := range chunkCh {
		if chunk.Error != nil {
			t.Fatalf("Expected no error in chunk, got %v", chunk.Error)
		}

		if chunk.Content.Type == multimodal.MediaTypeText {
			chunks = append(chunks, chunk.Content.Text)
		}

		if chunk.IsFinal {
			finalChunk = true
		}
	}

	// Check chunks
	if len(chunks) == 0 {
		t.Fatalf("Expected chunks to be populated")
	}

	// Check final chunk
	if !finalChunk {
		t.Errorf("Expected final chunk")
	}

	// Join chunks to check full response
	response := strings.Join(chunks, "")
	expected := "This is a test response."
	if response != expected {
		t.Errorf("Expected full response to be '%s', got '%s'", expected, response)
	}
}
