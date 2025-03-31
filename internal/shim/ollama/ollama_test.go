package ollama

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOllamaShim(t *testing.T) {
	// Test with defaults
	shim := NewOllamaShim("", "llama3", "", 0)
	if shim.BaseURL != "http://localhost:11434" {
		t.Errorf("Expected default BaseURL to be http://localhost:11434, got %s", shim.BaseURL)
	}
	if shim.Model != "llama3" {
		t.Errorf("Expected Model to be llama3, got %s", shim.Model)
	}
	if shim.MaxTokens != 4096 {
		t.Errorf("Expected default MaxTokens to be 4096, got %d", shim.MaxTokens)
	}

	// Test with custom values
	shim = NewOllamaShim("https://custom.endpoint", "custom-model", "api-key", 1000)
	if shim.BaseURL != "https://custom.endpoint" {
		t.Errorf("Expected BaseURL to be https://custom.endpoint, got %s", shim.BaseURL)
	}
	if shim.Model != "custom-model" {
		t.Errorf("Expected Model to be custom-model, got %s", shim.Model)
	}
	if shim.APIKey != "api-key" {
		t.Errorf("Expected APIKey to be api-key, got %s", shim.APIKey)
	}
	if shim.MaxTokens != 1000 {
		t.Errorf("Expected MaxTokens to be 1000, got %d", shim.MaxTokens)
	}

	// Test URL trimming
	shim = NewOllamaShim("https://endpoint.com/", "model", "", 0)
	if shim.BaseURL != "https://endpoint.com" {
		t.Errorf("Expected BaseURL to have trailing slash removed, got %s", shim.BaseURL)
	}
}

func TestCompleteChatPrompt(t *testing.T) {
	// Create a test server that returns a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request URL
		if r.URL.Path != "/api/chat" {
			t.Errorf("Expected request to /api/chat, got %s", r.URL.Path)
		}

		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", contentType)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"model": "llama3",
			"message": {
				"role": "assistant",
				"content": "This is a test response"
			},
			"done": true
		}`))
	}))
	defer server.Close()

	// Create a shim that uses the test server
	shim := NewOllamaShim(server.URL, "llama3", "", 0)

	// Test the completion function
	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "Hello, world!"},
	}
	result, err := shim.CompleteChatPrompt(messages)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the result
	expected := "This is a test response"
	if result != expected {
		t.Errorf("Expected result to be '%s', got '%s'", expected, result)
	}
}

// Skip the ParseSentinelfile test for now as it's more complex
// and would require mocking the CompleteChatPrompt function
