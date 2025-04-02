package shim

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

func TestShimFactory(t *testing.T) {
	// Test creating each type of shim
	providers := []string{
		ProviderClaude,
		ProviderOpenAI,
		ProviderOllama,
		ProviderGoogle,
		ProviderMock,
	}

	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			// Use empty values for API test
			shim, err := ShimFactory(provider, "", "", "")
			if err != nil {
				t.Fatalf("Failed to create shim for provider %s: %v", provider, err)
			}
			
			if shim == nil {
				t.Fatalf("ShimFactory returned nil for provider %s", provider)
			}
			
			// Check that the shim implements the LLMShim interface
			switch provider {
			case ProviderClaude:
				_, ok := shim.(*ClaudeShim)
				if !ok {
					t.Errorf("Expected *ClaudeShim, got %T", shim)
				}
			case ProviderOpenAI:
				_, ok := shim.(*OpenAIShim)
				if !ok {
					t.Errorf("Expected *OpenAIShim, got %T", shim)
				}
			case ProviderOllama:
				_, ok := shim.(*OllamaShim)
				if !ok {
					t.Errorf("Expected *OllamaShim, got %T", shim)
				}
			case ProviderGoogle:
				_, ok := shim.(*GoogleShim)
				if !ok {
					t.Errorf("Expected *GoogleShim, got %T", shim)
				}
			case ProviderMock:
				_, ok := shim.(*MockShim)
				if !ok {
					t.Errorf("Expected *MockShim, got %T", shim)
				}
			}
		})
	}
}

func TestMockShim(t *testing.T) {
	// Create a mock shim for testing
	config := Config{
		Provider: ProviderMock,
		Model:    "test-model",
	}
	
	shim := NewMockShim(config)
	
	// Test system prompt
	systemPrompt := "You are a test assistant."
	shim.SetSystemPrompt(systemPrompt)
	
	// Test completion
	prompt := "Hello, world!"
	response, err := shim.Completion(prompt, 100, 0.7, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to generate completion: %v", err)
	}
	
	if response == "" {
		t.Error("Expected non-empty response, got empty string")
	}
	
	// Test context-aware completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	response, err = shim.CompletionWithContext(ctx, prompt, 100, 0.7)
	if err != nil {
		t.Fatalf("Failed to generate completion with context: %v", err)
	}
	
	if response == "" {
		t.Error("Expected non-empty response from CompletionWithContext, got empty string")
	}
	
	// Test multimodal completion
	input := multimodal.NewInput()
	input.AddText(prompt)
	
	output, err := shim.MultimodalCompletion(input, 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to generate multimodal completion: %v", err)
	}
	
	if output == nil {
		t.Fatal("Expected non-nil output, got nil")
	}
	
	if output.GetText() == "" {
		t.Error("Expected non-empty text in multimodal output, got empty string")
	}
	
	// Test streaming
	streamCh, err := shim.StreamCompletion(ctx, prompt, 100, 0.7)
	if err != nil {
		t.Fatalf("Failed to start streaming completion: %v", err)
	}
	
	var streamedText string
	for chunk := range streamCh {
		streamedText += chunk
	}
	
	if streamedText == "" {
		t.Error("Expected non-empty streamed text, got empty string")
	}
	
	// Test multimodal streaming
	multiStreamCh, err := shim.StreamMultimodalCompletion(ctx, input)
	if err != nil {
		t.Fatalf("Failed to start multimodal streaming: %v", err)
	}
	
	streamedText = ""
	for chunk := range multiStreamCh {
		if chunk.Content.Type == multimodal.MediaTypeText {
			streamedText += chunk.Content.Text
		}
	}
	
	if streamedText == "" {
		t.Error("Expected non-empty streamed multimodal text, got empty string")
	}
	
	// Test Sentinelfile parsing
	sentinelfileContent := `# Sentinelfile for TestAgent

This agent helps with testing.

It should be able to:
- Test things
- Verify results

The agent should use test-model as its base model.

Allow the agent to access the following tools:
- test/tool1
- test/tool2
`

	result, err := shim.ParseSentinelfile(sentinelfileContent)
	if err != nil {
		t.Fatalf("Failed to parse Sentinelfile: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil parsing result, got nil")
	}
	
	// Test multimodal support
	if !shim.SupportsMultimodal() {
		t.Error("Expected mock shim to support multimodal inputs")
	}
	
	// Test close
	err = shim.Close()
	if err != nil {
		t.Fatalf("Failed to close shim: %v", err)
	}
}

func TestOllamaShim(t *testing.T) {
	// Skip real API tests if not in integration test mode
	if os.Getenv("SENTINEL_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping Ollama API test; set SENTINEL_INTEGRATION_TEST=1 to run")
	}
	
	// Get endpoint from environment or use default
	endpoint := os.Getenv("SENTINEL_OLLAMA_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/generate"
	}
	
	// Get model from environment or use default
	model := os.Getenv("SENTINEL_OLLAMA_MODEL")
	if model == "" {
		model = "llama3"
	}
	
	// Create config
	config := Config{
		Provider: ProviderOllama,
		Model:    model,
		Endpoint: endpoint,
		Timeout:  30 * time.Second,
	}
	
	// Create shim
	shim := NewOllamaShim(config)
	
	// Test completion
	prompt := "Hello, world!"
	response, err := shim.Completion(prompt, 100, 0.7, 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to generate Ollama completion: %v", err)
	}
	
	if response == "" {
		t.Error("Expected non-empty response from Ollama, got empty string")
	}
	
	// Test Sentinelfile parsing
	sentinelfileContent := `# Sentinelfile for TestAgent

This agent helps with testing.

It should be able to:
- Test things
- Verify results

The agent should use test-model as its base model.

Allow the agent to access the following tools:
- test/tool1
- test/tool2
`

	result, err := shim.ParseSentinelfile(sentinelfileContent)
	if err != nil {
		t.Fatalf("Failed to parse Sentinelfile with Ollama: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil parsing result from Ollama, got nil")
	}
}

func TestGoogleShim(t *testing.T) {
	// Skip real API tests if not in integration test mode
	if os.Getenv("SENTINEL_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping Google API test; set SENTINEL_INTEGRATION_TEST=1 to run")
	}
	
	// Get API key from environment
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping Google API test; GOOGLE_API_KEY not set")
	}
	
	// Get model from environment or use default
	model := os.Getenv("SENTINEL_GOOGLE_MODEL")
	if model == "" {
		model = "gemini-1.5-pro"
	}
	
	// Create config
	config := Config{
		Provider: ProviderGoogle,
		Model:    model,
		APIKey:   apiKey,
		Endpoint: DefaultEndpoints[ProviderGoogle],
		Timeout:  30 * time.Second,
	}
	
	// Create shim
	shim := NewGoogleShim(config)
	
	// Test completion
	prompt := "Hello, world!"
	response, err := shim.Completion(prompt, 100, 0.7, 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to generate Google completion: %v", err)
	}
	
	if response == "" {
		t.Error("Expected non-empty response from Google, got empty string")
	}
	
	// Test Sentinelfile parsing
	sentinelfileContent := `# Sentinelfile for TestAgent

This agent helps with testing.

It should be able to:
- Test things
- Verify results

The agent should use test-model as its base model.

Allow the agent to access the following tools:
- test/tool1
- test/tool2
`

	result, err := shim.ParseSentinelfile(sentinelfileContent)
	if err != nil {
		t.Fatalf("Failed to parse Sentinelfile with Google: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil parsing result from Google, got nil")
	}
}

func TestProviderUtils(t *testing.T) {
	// Test IsMultimodalModel
	testCases := []struct {
		provider string
		model    string
		expected bool
	}{
		{ProviderClaude, "claude-3-opus", true},
		{ProviderClaude, "claude-3-sonnet", true},
		{ProviderClaude, "claude-2", false},
		{ProviderOpenAI, "gpt-4-vision-preview", true},
		{ProviderOpenAI, "gpt-3.5-turbo", false},
		{ProviderOllama, "llava", true},
		{ProviderOllama, "llava:latest", true},
		{ProviderOllama, "llama3", false},
		{ProviderGoogle, "gemini-pro-vision", true},
		{ProviderGoogle, "gemini-1.5-pro", true},
		{ProviderGoogle, "gemini-pro", false},
		{ProviderMock, "any-model", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.provider+"/"+tc.model, func(t *testing.T) {
			result := IsMultimodalModel(tc.provider, tc.model)
			if result != tc.expected {
				t.Errorf("IsMultimodalModel(%s, %s) = %v, expected %v", 
					tc.provider, tc.model, result, tc.expected)
			}
		})
	}
	
	// Test GetDefaultSystemPrompt
	systemPrompts := []string{
		GetDefaultSystemPrompt(ProviderClaude, "agent"),
		GetDefaultSystemPrompt(ProviderClaude, "parser"),
		GetDefaultSystemPrompt(ProviderOpenAI, "parser"),
		GetDefaultSystemPrompt(ProviderOllama, "parser"),
		GetDefaultSystemPrompt(ProviderGoogle, "parser"),
		GetDefaultSystemPrompt("unknown", "unknown"),
	}
	
	for _, prompt := range systemPrompts {
		if prompt == "" {
			t.Error("Expected non-empty system prompt, got empty string")
		}
	}
	
	// Test environment helpers
	os.Setenv("SENTINEL_LLM_PROVIDER", "google")
	os.Setenv("SENTINEL_LLM_MODEL", "gemini-1.5-pro")
	os.Setenv("GOOGLE_API_KEY", "test-key")
	
	provider := GetProviderFromEnv()
	if provider != "google" {
		t.Errorf("Expected provider 'google', got '%s'", provider)
	}
	
	model := GetModelFromEnv(provider)
	if model != "gemini-1.5-pro" {
		t.Errorf("Expected model 'gemini-1.5-pro', got '%s'", model)
	}
	
	apiKey := GetAPIKeyFromEnv()
	if apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", apiKey)
	}
	
	// Restore environment
	os.Unsetenv("SENTINEL_LLM_PROVIDER")
	os.Unsetenv("SENTINEL_LLM_MODEL")
	os.Unsetenv("GOOGLE_API_KEY")
}
