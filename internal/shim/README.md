# SentinelStacks LLM Shims

This directory contains implementations of the LLM shim interface for different LLM providers. The shims provide a consistent interface for interacting with different LLM APIs.

## Overview

The shim system provides a unified interface for:

1. Text completions
2. Multimodal completions (text + images)
3. Streaming responses
4. Sentinelfile parsing for agent creation

## Supported Providers

The following LLM providers are supported:

- **Claude** (Anthropic): For Claude models (claude-3-opus, claude-3-sonnet, etc.)
- **OpenAI**: For GPT models (gpt-4-turbo, gpt-3.5-turbo, etc.)
- **Google**: For Gemini models (gemini-1.5-pro, gemini-pro-vision, etc.)
- **Ollama**: For local open-source models (llama3, mistral, etc.)
- **Mock**: For testing and development without API access

## Architecture

### Core Interface

The `LLMShim` interface in `shim.go` defines the contract that all provider implementations must fulfill:

```go
// LLMShim is an interface for interacting with different LLM providers
type LLMShim interface {
    // Text completion methods
    Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error)
    CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error)
    
    // Multimodal methods
    MultimodalCompletion(input *multimodal.Input, timeout time.Duration) (*multimodal.Output, error)
    MultimodalCompletionWithContext(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error)
    
    // Streaming methods
    StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error)
    StreamMultimodalCompletion(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error)
    
    // System prompts
    SetSystemPrompt(prompt string)
    
    // Utility methods
    ParseSentinelfile(content string) (map[string]interface{}, error)
    SupportsMultimodal() bool
    Close() error
}
```

### Provider-Specific Implementations

Each provider has its own implementation:

- `claude_shim.go` and `claude/claude.go`: Anthropic Claude implementation
- `openai_shim.go` and `openai/openai.go`: OpenAI GPT implementation
- `google_shim.go` and `google/google.go`: Google Gemini implementation
- `ollama_shim.go`: Ollama (local/self-hosted models) implementation
- `shim.go` (MockShim): Simple mock implementation for testing

### Factory Pattern

The `ShimFactory` function creates the appropriate shim based on the provider:

```go
func ShimFactory(provider, endpoint, apiKey, model string) (LLMShim, error)
```

## Configuration

Shims can be configured using environment variables:

- `SENTINEL_LLM_PROVIDER`: The LLM provider to use (claude, openai, google, ollama, mock)
- `SENTINEL_LLM_MODEL`: The model to use
- `SENTINEL_LLM_ENDPOINT`: The API endpoint URL
- `SENTINEL_API_KEY`: The API key for authentication

Provider-specific API keys can also be set:
- `ANTHROPIC_API_KEY`: For Claude models
- `OPENAI_API_KEY`: For OpenAI models
- `GOOGLE_API_KEY`: For Google models

## Multimodal Support

The following models support multimodal inputs (text + images):

- **Claude**: claude-3-opus, claude-3-sonnet, claude-3-haiku, claude-3.5-sonnet
- **OpenAI**: gpt-4-vision-preview, gpt-4-turbo
- **Google**: gemini-pro-vision, gemini-1.5-pro, gemini-1.5-flash
- **Ollama**: llava, bakllava, moondream, fuyu, yi-vl, cogvlm

## Sentinelfile Parsing

One of the key functions of the shims is to parse natural language Sentinelfiles into structured agent definitions. Each provider implementation includes a `ParseSentinelfile` method that:

1. Sets a specialized system prompt for Sentinelfile parsing
2. Sends the Sentinelfile content to the LLM
3. Extracts structured information as JSON
4. Validates and returns the parsed data

## Tool Integration

The `tool_integration.go` file contains utilities for integrating LLMs with tools:

- `ToolExecutor`: Interface for executing tools
- `ToolAugmentedInput`: Adds tool descriptions to inputs
- `ParseFunctionCallFromLLMResponse`: Extracts function calls from responses

## Testing

The `shim_test.go` file contains tests for the shim implementations. To run integration tests with real APIs:

```bash
SENTINEL_INTEGRATION_TEST=1 go test -v ./internal/shim
```

To test with specific providers:

```bash
# For Google
SENTINEL_INTEGRATION_TEST=1 GOOGLE_API_KEY=your_key go test -run TestGoogleShim -v ./internal/shim

# For Claude
SENTINEL_INTEGRATION_TEST=1 ANTHROPIC_API_KEY=your_key go test -run TestClaudeShim -v ./internal/shim

# For OpenAI
SENTINEL_INTEGRATION_TEST=1 OPENAI_API_KEY=your_key go test -run TestOpenAIShim -v ./internal/shim

# For Ollama
SENTINEL_INTEGRATION_TEST=1 go test -run TestOllamaShim -v ./internal/shim
```

## Usage Example

```go
// Create a shim based on environment variables
shim, err := shim.CreateShimFromEnv()
if err != nil {
    log.Fatalf("Failed to create shim: %v", err)
}

// Set a system prompt
shim.SetSystemPrompt("You are a helpful assistant.")

// Generate a completion
response, err := shim.Completion("Hello, world!", 100, 0.7, 30*time.Second)
if err != nil {
    log.Fatalf("Failed to generate completion: %v", err)
}

fmt.Println("Response:", response)
```

## Using Google Gemini

To use Google's Gemini models:

```go
// Set environment variables
os.Setenv("SENTINEL_LLM_PROVIDER", "google")
os.Setenv("SENTINEL_LLM_MODEL", "gemini-1.5-pro")
os.Setenv("GOOGLE_API_KEY", "your_google_api_key")

// Create a shim
shim, err := shim.CreateShimFromEnv()
if err != nil {
    log.Fatalf("Failed to create shim: %v", err)
}

// Or create it directly
config := shim.Config{
    Provider: shim.ProviderGoogle,
    Model:    "gemini-1.5-pro",
    APIKey:   "your_google_api_key",
}
googleShim := shim.NewGoogleShim(config)

// Use it like any other shim
response, err := googleShim.Completion("Hello, Gemini!", 100, 0.7, 30*time.Second)
```

## Multimodal Example

```go
// Create a multimodal input
input := multimodal.NewInput()
input.AddText("What's in this image?")

// Add an image
imageData, _ := ioutil.ReadFile("image.jpg")
input.AddImage(imageData, "image/jpeg")

// Generate a multimodal completion
output, err := shim.MultimodalCompletion(input, 30*time.Second)
if err != nil {
    log.Fatalf("Failed to generate multimodal completion: %v", err)
}

fmt.Println("Response:", output.GetText())
```

## Streaming Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Start streaming
streamCh, err := shim.StreamCompletion(ctx, "Generate a story about...", 1000, 0.7)
if err != nil {
    log.Fatalf("Failed to start streaming: %v", err)
}

// Consume the stream
for chunk := range streamCh {
    fmt.Print(chunk)
}
```
