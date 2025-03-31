# Multimodal Support Implementation Plan

This document outlines the plan for adding multimodal capabilities to SentinelStacks, allowing agents to process and generate various media types including images, audio, and potentially video in the future.

## Overview

Multimodal support will enable SentinelStacks agents to:

1. Process images alongside text
2. Generate or modify images
3. Understand visual content and respond based on it
4. Interact through rich media interfaces

## Implementation Phases

### Phase 1: Core Framework (Current)

- [x] Define core multimodal data structures
- [ ] Implement media content handling utilities
- [ ] Create provider interface extensions
- [ ] Add test fixtures and unit tests

### Phase 2: Provider Integration

- [ ] Implement Claude provider multimodal support
- [ ] Add OpenAI provider multimodal capabilities
- [ ] Create testing tools for multimodal agents

### Phase 3: CLI and User Experience

- [ ] Extend Sentinelfile format for multimodal definitions
- [ ] Implement CLI commands for multimodal interaction
- [ ] Add documentation and examples

### Phase 4: Advanced Features

- [ ] Image preprocessing capabilities
- [ ] Document analysis
- [ ] Performance optimization for large media files

## Technical Details

### Core Types (Implemented)

The multimodal package defines several key types:

```go
// MediaType represents the type of media content
type MediaType string

// Media type constants
const (
    MediaTypeText  MediaType = "text"
    MediaTypeImage MediaType = "image"
    MediaTypeAudio MediaType = "audio"
    MediaTypeVideo MediaType = "video"
)

// Content represents a piece of multimodal content
type Content struct {
    Type     MediaType              `json:"type"`
    Data     []byte                 `json:"-"` // Raw binary data
    MimeType string                 `json:"mime_type,omitempty"`
    Text     string                 `json:"text,omitempty"`
    URL      string                 `json:"url,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Input represents a multimodal input for generation
type Input struct {
    Contents    []*Content              `json:"contents"`
    MaxTokens   int                     `json:"max_tokens,omitempty"`
    Temperature float64                 `json:"temperature,omitempty"`
    Stream      bool                    `json:"stream,omitempty"`
    Metadata    map[string]interface{}  `json:"metadata,omitempty"`
}

// Output represents a multimodal output from generation
type Output struct {
    Contents   []*Content              `json:"contents"`
    UsedTokens int                     `json:"used_tokens,omitempty"`
    Metadata   map[string]interface{}  `json:"metadata,omitempty"`
}
```

### Provider Interface Extensions

The shim interface has been extended to support multimodal capabilities:

```go
// Multimodal support methods in the Shim interface
GenerateMultimodal(ctx context.Context, input *multimodal.Input) (*multimodal.Output, error)
StreamMultimodal(ctx context.Context, input *multimodal.Input) (<-chan *multimodal.Chunk, error)
SupportsMultimodal() bool
```

### Sentinelfile Extensions (Planned)

Sentinelfiles will be extended to support multimodal capabilities:

```yaml
name: VisualAnalysisAgent
description: An agent that analyzes images and provides insights
baseModel: claude-3-opus-20240229
multimodal:
  enabled: true
  supportedMediaTypes:
    - image/jpeg
    - image/png
    - image/gif
  maxImageSize: 5MB
  imageAnalysisCapabilities:
    - objectDetection
    - textRecognition
    - sceneClassification
```

### CLI Command Extensions (Planned)

The CLI will be extended with commands to handle multimodal interaction:

```
sentinel run --image path/to/image.jpg --prompt "Analyze this image"
sentinel run --input multimodal_input.json
```

## Provider-Specific Implementation

### Claude Provider

Claude 3 models support multimodal inputs natively. The implementation will:

1. Convert SentinelStacks multimodal content to Anthropic API format
2. Handle API responses and convert them back to SentinelStacks format
3. Support both synchronous and streaming responses

Example API request format:

```json
{
  "model": "claude-3-opus-20240229",
  "messages": [
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "What's in this image?"
        },
        {
          "type": "image",
          "source": {
            "type": "base64",
            "media_type": "image/jpeg",
            "data": "..."
          }
        }
      ]
    }
  ],
  "max_tokens": 1024
}
```

### OpenAI Provider

OpenAI's GPT-4V supports multimodal inputs. The implementation will:

1. Convert SentinelStacks multimodal content to OpenAI API format
2. Handle API responses and convert them back to SentinelStacks format
3. Support both synchronous and streaming responses

Example API request format:

```json
{
  "model": "gpt-4-vision-preview",
  "messages": [
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "What's in this image?"
        },
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/jpeg;base64,..."
          }
        }
      ]
    }
  ],
  "max_tokens": 1024
}
```

## Example Use Cases

### Visual Analysis Agent

```yaml
name: VisualAnalyzer
description: An agent that analyzes images and provides detailed descriptions
baseModel: claude-3-opus-20240229
multimodal:
  enabled: true
  supportedMediaTypes:
    - image/jpeg
    - image/png
  maxImageSize: 5MB
```

### Design Assistant

```yaml
name: DesignAssistant
description: An agent that helps with design tasks and can generate mockups
baseModel: gpt-4-vision-preview
multimodal:
  enabled: true
  supportedMediaTypes:
    - image/jpeg
    - image/png
    - image/svg+xml
  maxImageSize: 10MB
  generationCapabilities:
    - mockups
    - iconDesign
    - uiSuggestions
```

## Next Steps

1. Complete the core multimodal implementation
2. Implement Claude provider integration
3. Add OpenAI provider support
4. Extend CLI for multimodal interaction
5. Create documentation and examples
6. Test with real-world use cases

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Large media files impacting performance | Implement automatic resizing and compression |
| API costs for multimodal models | Add usage tracking and limits |
| Provider API changes | Abstract provider implementations behind stable interfaces |
| Security concerns with media content | Add content scanning and validation |

## Conclusion

Multimodal support will significantly enhance the capabilities of SentinelStacks agents, allowing them to interact with and understand visual content alongside text. The implementation plan outlined here provides a structured approach to adding these capabilities while maintaining compatibility with existing systems. 