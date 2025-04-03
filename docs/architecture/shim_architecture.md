# Shim Architecture: LLM Provider Integration

The Shim architecture is a core component of SentinelStacks that enables agent portability across different LLM providers. This document explains how the shim system works and how it integrates with various LLM providers.

## Overview

The shim system provides a unified abstraction layer over different LLM providers (Claude, OpenAI, Google, Ollama). This architecture allows agents to be provider-agnostic, meaning they can run on any supported LLM without code changes.

## Core Components

### LLMShim Interface

The LLMShim interface (`pkg/types/shim.go`) defines the standardized API for all LLM providers:

```go
// LLMShim is an interface for interacting with different LLM providers
type LLMShim interface {
    // Text completion methods
    Completion(prompt string, maxTokens int, temperature float64, timeout time.Duration) (string, error)
    CompletionWithContext(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error)
    
    // Multimodal methods
    MultimodalCompletion(input *Input, timeout time.Duration) (*Output, error)
    MultimodalCompletionWithContext(ctx context.Context, input *Input) (*Output, error)
    
    // Streaming methods
    StreamCompletion(ctx context.Context, prompt string, maxTokens int, temperature float64) (<-chan string, error)
    StreamMultimodalCompletion(ctx context.Context, input *Input) (<-chan *Chunk, error)
    
    // System prompts
    SetSystemPrompt(prompt string)
    
    // Utility methods
    ParseSentinelfile(content string) (map[string]interface{}, error)
    SupportsMultimodal() bool
    Close() error
}
```

### ShimFactory

The ShimFactory (`internal/shim/shim.go`) creates the appropriate shim implementation based on the specified provider:

```go
// ShimFactory creates a new LLM shim based on the provider
func ShimFactory(provider, endpoint, apiKey, model string) (LLMShim, error) {
    // Create config
    config := Config{
        Provider: provider,
        Model:    model,
        APIKey:   apiKey,
        Endpoint: endpoint,
        Timeout:  60 * time.Second, // Default timeout
    }
    
    // Create shim based on provider
    switch provider {
    case "claude":
        return NewClaudeShim(config), nil
    case "openai":
        return NewOpenAIShim(config), nil
    case "ollama":
        return NewOllamaShim(config), nil
    case "google":
        return NewGoogleShim(config), nil
    case "mock":
        return NewMockShim(config), nil
    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
}
```

## Provider-Specific Implementations

Each LLM provider has its own implementation that wraps the provider's API:

### Claude Shim

The Claude Shim (`internal/shim/claude_shim.go`) integrates with Anthropic's Claude API:

- Supports text and multimodal inputs
- Implements streaming for real-time responses
- Handles Claude-specific parameters and models

### OpenAI Shim

The OpenAI Shim (`internal/shim/openai_shim.go`) wraps OpenAI's API:

- Supports GPT models (3.5-turbo, 4, etc.)
- Handles multimodal inputs for vision-capable models
- Implements OpenAI's chat completion API with system prompts

### Google Shim

The Google Shim (`internal/shim/google_shim.go`) integrates with Google's Vertex AI and Gemini models:

- Supports Google's various AI models
- Handles Google-specific authentication and API calls
- Provides access to Gemini's multimodal capabilities

### Ollama Shim

The Ollama Shim (`internal/shim/ollama_shim.go`) enables using locally-hosted open-source models:

- Supports running models like Llama, Mistral, etc. via Ollama
- Implements a subset of features based on Ollama's capabilities
- Enables offline usage without external API dependencies

### Mock Shim

The Mock Shim is used for testing:

- Simulates LLM responses without calling actual APIs
- Useful for development and testing
- Supports all interface methods with predefined responses

## Key Features

### Unified Abstraction

The shim provides a common interface for all providers, allowing agent code to remain provider-agnostic:

```go
// Example usage in agent code
shim, err := ShimFactory("claude", "", apiKey, "claude-3-sonnet")
if err != nil {
    // Handle error
}

// The same code works regardless of provider
response, err := shim.Completion("What is SentinelStacks?", 1000, 0.7, 30*time.Second)
```

### Multimodal Support

All providers implement a standardized interface for handling multimodal content (text + images):

```go
// Create multimodal input
input := multimodal.NewInput()
input.AddText("What's in this image?")
input.AddImage(imageData, "image/jpeg")

// Process through any provider that supports multimodal
if shim.SupportsMultimodal() {
    output, err := shim.MultimodalCompletion(input, 60*time.Second)
    // Handle output
}
```

### Streaming Responses

The streaming API allows for real-time responses from all providers:

```go
// Stream responses from any provider
streamCh, err := shim.StreamCompletion(ctx, prompt, maxTokens, temperature)
if err != nil {
    // Handle error
}

// Process stream chunks
for chunk := range streamCh {
    // Process each chunk as it arrives
    fmt.Print(chunk)
}
```

## Implementation Pattern

Each shim implementation follows a common pattern:

1. **Configuration**: Accepts provider-specific configuration
2. **Delegation**: Delegates to internal provider-specific implementation
3. **Fallback**: Includes fallback mechanisms for graceful degradation
4. **Capability Checking**: Implements capability checks (e.g., multimodal support)
5. **Context Management**: Propagates context for cancellation and timeouts

## Using the Shim in Agents

In SentinelStacks, agents are created with a specific shim but can be migrated to different providers:

```go
// Create an agent with Claude
agent, err := agent.NewAgent("my-agent", "claude", "claude-3-sonnet", apiKey, "")

// The same agent definition can later run with OpenAI
agent, err := agent.NewAgent("my-agent", "openai", "gpt-4", apiKey, "")
```

## Extending with New Providers

To add support for a new LLM provider:

1. Create a new implementation file (e.g., `new_provider_shim.go`)
2. Implement the `LLMShim` interface for the new provider
3. Add the provider to the `ShimFactory` function
4. Update capability checks as needed

## Conclusion

The shim architecture is a key component that enables SentinelStacks' provider-agnostic approach to AI agents. By abstracting provider-specific details behind a unified interface, it allows agents to be portable across different LLMs while preserving their core functionality. 