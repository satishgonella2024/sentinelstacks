# Multimodal Support Implementation Plan

## Overview

This document outlines the implementation plan for adding multimodal support to SentinelStacks, enabling agents to process and generate not just text but also images, audio, and video content.

## Architecture Changes

### 1. Multimodal Types Package

We've created a new package `internal/multimodal` that defines the core types:

- `MediaType` - Enum for different media types (text, image, audio, video)
- `Content` - Struct to represent multimodal content with fields for type, data, MIME type, etc.
- `Input` - Container for multimodal inputs with methods to add text and image content
- `Output` - Container for multimodal outputs with methods to retrieve text and images
- `Chunk` - For streaming responses
- `ContentReader` - Interface for reading multimodal content

### 2. Updated Provider Interface

The LLM provider interface has been expanded to include multimodal capabilities:

```go
type Provider interface {
    // Existing methods...
    
    // Multimodal methods
    SupportsMultimodal() bool
    GenerateMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (*multimodal.Output, error)
    StreamMultimodalResponse(ctx context.Context, input *multimodal.Input, params map[string]interface{}) (<-chan *multimodal.Chunk, error)
}
```

### 3. Shim Interface Updates

The shim interface has been expanded to include:

```go
type Shim interface {
    // Existing methods...
    
    // Multimodal methods
    SupportsMultimodal() bool
    GenerateMultimodal(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error)
    StreamMultimodal(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error)
}
```

## Provider Implementations

### Claude Provider

Claude 3 models have strong multimodal capabilities. Implementation includes:

- Added `SupportsMultimodal()` method returning `true`
- Implemented `GenerateMultimodalResponse` and `StreamMultimodalResponse` methods
- Support for image inputs and text outputs

### OpenAI Provider (Planned)

Future implementation for OpenAI models (GPT-4 Vision):

- Will support image inputs and text outputs
- Will handle DALL-E 3 for image generation

### Ollama Provider (Planned)

- Multimodal support for models like LLaVA that support image inputs
- Implementation will depend on Ollama's API for multimodal capabilities

## CLI Enhancements

### New Commands

- `sentinel mm generate` - Generate a response using multimodal input
- `sentinel mm stream` - Stream a response using multimodal input

### Flags

- `--image` - Path to one or more image files to include
- `--audio` - Path to an audio file to include
- `--video` - Path to a video file to include
- `--output-dir` - Directory to save generated media outputs

## Sentinelfile Additions

New sections in the Sentinelfile for multimodal capabilities:

```yaml
multimodal:
  input_types:
    - text
    - image
  output_types:
    - text
    - image
  max_image_size: 4MB
  image_formats:
    - png
    - jpeg
```

## Testing Strategy

1. Unit tests for the multimodal types and functions
2. Integration tests with mock LLM providers
3. Full end-to-end tests with actual providers (Claude, OpenAI)
4. Performance testing with different media sizes and types

## Implementation Timeline

| Phase | Task | Status | ETA |
|-------|------|--------|-----|
| 1 | Core multimodal types | Completed | - |
| 1 | Shim interface updates | Completed | - |
| 2 | Claude provider implementation | In Progress | 2 weeks |
| 2 | CLI enhancements | Not Started | 3 weeks |
| 3 | OpenAI provider implementation | Not Started | 4 weeks |
| 3 | Ollama provider implementation | Not Started | 5 weeks |
| 4 | Documentation and examples | Not Started | 6 weeks |

## Future Enhancements

- Support for audio processing and generation
- Support for video processing and generation
- Vision-based agents that can analyze complex images
- Multi-turn multimodal conversations
- Agent-generated visual content (charts, diagrams, etc.)
- Multimodal state management for agents

## Challenges and Considerations

- **Performance**: Handling large media files efficiently
- **Security**: Ensuring safe processing of user-provided media
- **Compatibility**: Supporting different provider capabilities
- **Cost**: Managing token/API costs for multimodal interactions
- **Streaming**: Implementing efficient streaming for larger outputs

## Conclusion

Multimodal support represents a significant enhancement to SentinelStacks' capabilities, allowing developers to create more versatile and powerful agents. The implementation follows a phased approach, starting with the core types and interfaces, followed by provider-specific implementations and CLI enhancements. 