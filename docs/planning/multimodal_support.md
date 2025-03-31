# Multimodal Support Implementation Plan

This document outlines the plan for adding multimodal support to SentinelStacks, enabling agents to process and generate visual content alongside text.

## Overview

Multimodal support will allow SentinelStacks agents to:

1. Process images as input alongside text
2. Generate or modify images as output
3. Understand and reason about visual content
4. Interact with users through rich media interfaces

This capability is particularly valuable for use cases such as:
- Visual analysis agents (e.g., medical imaging, satellite imagery)
- Design assistants (e.g., UI/UX design, graphic design)
- Document analysis (e.g., form processing, document understanding)
- Educational tools with visual explanations
- Visual content moderation
- Multimodal content creation

## Architecture Changes

### 1. Shim Interface Enhancements

Extend the current Shim interface to support multimodal inputs and outputs:

```go
type Shim interface {
    // Existing methods
    Initialize(config Config) error
    Generate(ctx context.Context, input GenerateInput) (*GenerateOutput, error)
    Stream(ctx context.Context, input GenerateInput) (<-chan StreamChunk, error)
    GetEmbeddings(ctx context.Context, input EmbeddingsInput) (*EmbeddingsOutput, error)
    Close() error
    
    // New methods for multimodal support
    GenerateMultimodal(ctx context.Context, input MultimodalInput) (*MultimodalOutput, error)
    StreamMultimodal(ctx context.Context, input MultimodalInput) (<-chan MultimodalChunk, error)
}

// New types for multimodal support
type MediaType string

const (
    MediaTypeText  MediaType = "text"
    MediaTypeImage MediaType = "image"
    MediaTypeAudio MediaType = "audio"
    MediaTypeVideo MediaType = "video"
)

type MediaContent struct {
    Type        MediaType
    Data        []byte      // Raw binary data for the media
    MimeType    string      // MIME type of the media
    Text        string      // For text content or alternative text
    URL         string      // Optional URL for external media
    Metadata    map[string]interface{} // Additional metadata
}

type MultimodalInput struct {
    Contents    []MediaContent  // Ordered list of content pieces (text, images, etc.)
    MaxTokens   int             // Maximum tokens to generate
    Temperature float64         // Temperature parameter
    Tools       []Tool          // Available tools
    Stream      bool            // Whether to stream the response
    Metadata    map[string]interface{} // Additional metadata
}

type MultimodalOutput struct {
    Contents    []MediaContent  // Generated content (text, images, etc.)
    UsedTokens  int             // Number of tokens used
    ToolCalls   []ToolCall      // Tool calls made during generation
    Metadata    map[string]interface{} // Additional metadata
}

type MultimodalChunk struct {
    Content     MediaContent    // Chunk of generated content
    IsFinal     bool            // Whether this is the final chunk
    Error       error           // Error if any
}
```

### 2. Provider Implementations

Update provider implementations to support multimodal capabilities where available:

#### Claude Provider (Anthropic)

```go
func (p *ClaudeProvider) GenerateMultimodal(ctx context.Context, input MultimodalInput) (*MultimodalOutput, error) {
    // Convert our MultimodalInput to Claude's message format
    messages := []anthropic.Message{}
    
    // Create a system message if provided
    if systemPrompt, ok := input.Metadata["system_prompt"].(string); ok && systemPrompt != "" {
        messages = append(messages, anthropic.Message{
            Role: "system",
            Content: systemPrompt,
        })
    }
    
    // Convert our media contents to Claude's content format
    userContent := []anthropic.Content{}
    for _, content := range input.Contents {
        switch content.Type {
        case MediaTypeText:
            userContent = append(userContent, anthropic.Content{
                Type: "text",
                Text: content.Text,
            })
        case MediaTypeImage:
            // Convert image to Claude's format
            mediaType := content.MimeType
            if mediaType == "" {
                mediaType = "image/jpeg" // Default
            }
            
            userContent = append(userContent, anthropic.Content{
                Type: "image",
                Source: &anthropic.Source{
                    Type: "base64",
                    MediaType: mediaType,
                    Data: base64.StdEncoding.EncodeToString(content.Data),
                },
            })
        }
    }
    
    // Add user message
    messages = append(messages, anthropic.Message{
        Role: "user",
        Content: userContent,
    })
    
    // Prepare request parameters
    params := anthropic.MessageRequest{
        Model: p.getModelName(input.Metadata),
        Messages: messages,
        MaxTokens: input.MaxTokens,
        Temperature: input.Temperature,
        System: "", // We handle system prompt in messages
    }
    
    // Call Claude API
    resp, err := p.client.Messages(ctx, params)
    if err != nil {
        return nil, err
    }
    
    // Convert Claude's response to our format
    output := &MultimodalOutput{
        Contents: []MediaContent{},
        UsedTokens: resp.Usage.OutputTokens,
        ToolCalls: []ToolCall{},
        Metadata: map[string]interface{}{
            "model": resp.Model,
            "stop_reason": resp.StopReason,
        },
    }
    
    // Convert response content
    for _, content := range resp.Content {
        switch content.Type {
        case "text":
            output.Contents = append(output.Contents, MediaContent{
                Type: MediaTypeText,
                Text: content.Text,
            })
        // Handle other content types if Claude returns them
        }
    }
    
    return output, nil
}
```

#### OpenAI Provider

```go
func (p *OpenAIProvider) GenerateMultimodal(ctx context.Context, input MultimodalInput) (*MultimodalOutput, error) {
    // Convert our MultimodalInput to OpenAI's message format
    messages := []openai.ChatCompletionMessage{}
    
    // Create a system message if provided
    if systemPrompt, ok := input.Metadata["system_prompt"].(string); ok && systemPrompt != "" {
        messages = append(messages, openai.ChatCompletionMessage{
            Role: "system",
            Content: systemPrompt,
        })
    }
    
    // Convert our media contents to OpenAI's content format
    userContent := []openai.ChatMessagePart{}
    for _, content := range input.Contents {
        switch content.Type {
        case MediaTypeText:
            userContent = append(userContent, openai.ChatMessagePart{
                Type: "text",
                Text: content.Text,
            })
        case MediaTypeImage:
            // Convert image to OpenAI's format
            userContent = append(userContent, openai.ChatMessagePart{
                Type: "image_url",
                ImageURL: &openai.ImageURL{
                    URL: fmt.Sprintf("data:%s;base64,%s", 
                          content.MimeType, 
                          base64.StdEncoding.EncodeToString(content.Data)),
                },
            })
        }
    }
    
    // Add user message
    messages = append(messages, openai.ChatCompletionMessage{
        Role: "user",
        Content: userContent,
    })
    
    // Prepare request parameters
    params := openai.ChatCompletionRequest{
        Model: p.getModelName(input.Metadata),
        Messages: messages,
        MaxTokens: input.MaxTokens,
        Temperature: input.Temperature,
    }
    
    // Call OpenAI API
    resp, err := p.client.CreateChatCompletion(ctx, params)
    if err != nil {
        return nil, err
    }
    
    // Convert OpenAI's response to our format
    output := &MultimodalOutput{
        Contents: []MediaContent{},
        UsedTokens: resp.Usage.CompletionTokens,
        ToolCalls: []ToolCall{},
        Metadata: map[string]interface{}{
            "model": resp.Model,
            "finish_reason": resp.Choices[0].FinishReason,
        },
    }
    
    // Convert response content
    output.Contents = append(output.Contents, MediaContent{
        Type: MediaTypeText,
        Text: resp.Choices[0].Message.Content,
    })
    
    return output, nil
}
```

### 3. Agent Definition Extensions

Update the Sentinelfile specification to include multimodal capabilities:

```yaml
name: visual-analysis-agent
description: An agent that can analyze and process images
capabilities:
  - Process images and provide detailed descriptions
  - Extract text from images
  - Identify objects and scenes in photos
  - Generate image modifications based on instructions
model:
  base: claude-3-opus
  vision: true  # Indicates this agent has vision capabilities
  parameters:
    temperature: 0.7
    top_p: 0.9
media_handling:
  supported_types:
    - image/jpeg
    - image/png
    - image/gif
  max_size: 5MB
  default_alt_text: "Image submitted for analysis"
  content_policy:
    adult_content: block
    violent_content: warn
tools:
  - image_analyzer:
      purpose: For detailed image analysis
  - ocr:
      purpose: For extracting text from images
  - image_editor:
      purpose: For simple image modifications
initialization:
  introduction: "Hello! I'm a visual analysis assistant. You can send me images, and I'll help analyze them."
termination:
  farewell: "Thank you for using the visual analysis assistant."
```

### 4. CLI Integration

Update the CLI to support image uploads and multimodal interactions:

```go
// New CLI command for multimodal file support
var cmdMultimodal = &cobra.Command{
    Use:   "run-multimodal [image] [agent]",
    Short: "Run an agent with multimodal input",
    Long:  `Run a Sentinel Agent with multimodal capabilities, providing both text and image inputs.`,
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // Parse arguments
        imagePath := args[0]
        agentName := ""
        if len(args) > 1 {
            agentName = args[1]
        }
        
        // Load image file
        imageData, err := os.ReadFile(imagePath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading image file: %v\n", err)
            os.Exit(1)
        }
        
        // Determine MIME type
        mimeType := http.DetectContentType(imageData)
        
        // Get initial prompt
        initialPrompt, _ := cmd.Flags().GetString("prompt")
        
        // Load agent
        agent, err := loadAgent(agentName)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error loading agent: %v\n", err)
            os.Exit(1)
        }
        
        // Create multimodal input
        input := MultimodalInput{
            Contents: []MediaContent{
                {
                    Type: MediaTypeImage,
                    Data: imageData,
                    MimeType: mimeType,
                },
                {
                    Type: MediaTypeText,
                    Text: initialPrompt,
                },
            },
            MaxTokens: viper.GetInt("generation.max_tokens"),
            Temperature: viper.GetFloat64("generation.temperature"),
        }
        
        // Process through agent
        output, err := agent.ProcessMultimodal(input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error processing multimodal input: %v\n", err)
            os.Exit(1)
        }
        
        // Display results
        for _, content := range output.Contents {
            if content.Type == MediaTypeText {
                fmt.Println(content.Text)
            } else if content.Type == MediaTypeImage {
                // Save generated image
                outputPath := "output." + strings.Split(content.MimeType, "/")[1]
                err := os.WriteFile(outputPath, content.Data, 0644)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Error saving output image: %v\n", err)
                } else {
                    fmt.Printf("Generated image saved to %s\n", outputPath)
                }
            }
        }
    },
}
```

### 5. Agent Runtime Extensions

Enhance the agent runtime to handle multimodal processing:

```go
// Agent extensions for multimodal support
func (a *Agent) ProcessMultimodal(input MultimodalInput) (*MultimodalOutput, error) {
    // Check if agent supports multimodal
    if !a.SupportsMultimodal() {
        return nil, fmt.Errorf("agent does not support multimodal inputs")
    }
    
    // Preprocess input (validate, resize images if needed, etc.)
    processedInput, err := a.preprocessMultimodalInput(input)
    if err != nil {
        return nil, err
    }
    
    // Apply any agent-specific configurations
    if a.config.MaxImageSize > 0 {
        processedInput = a.enforceImageSizeLimit(processedInput)
    }
    
    // Add system prompt based on agent definition
    if a.systemPrompt != "" {
        if processedInput.Metadata == nil {
            processedInput.Metadata = map[string]interface{}{}
        }
        processedInput.Metadata["system_prompt"] = a.systemPrompt
    }
    
    // Process through LLM shim
    output, err := a.shim.GenerateMultimodal(context.Background(), processedInput)
    if err != nil {
        return nil, err
    }
    
    // Postprocess output (validate, filter, etc.)
    return a.postprocessMultimodalOutput(output)
}

func (a *Agent) SupportsMultimodal() bool {
    // Check if agent definition includes multimodal capabilities
    return a.definition.Vision || a.definition.AudioProcessing
}
```

## Implementation Phases

### Phase 1: Core Framework (Weeks 1-2)

1. Design and implement core data structures for multimodal content
2. Extend shim interface to support multimodal inputs and outputs
3. Update agent definition format to include multimodal capabilities
4. Create basic media utilities (image loading/saving, format conversion)
5. Add unit tests for new components

### Phase 2: Provider Integration (Weeks 3-4)

1. Implement multimodal support for Claude (Anthropic) provider
2. Add multimodal support for OpenAI (GPT-4 Vision) provider
3. Create testing tools and mock responses for development
4. Develop provider-specific optimizations and error handling
5. Add integration tests with real API calls

### Phase 3: CLI and UX (Weeks 5-6)

1. Extend CLI to support multimodal inputs
2. Create interactive mode for multimodal conversations
3. Implement image display in terminal (where supported)
4. Add file export capabilities for generated media
5. Create documentation and examples

### Phase 4: Advanced Features (Weeks 7-8)

1. Implement image preprocessing tools (resize, crop, filter)
2. Add support for PDF and document analysis
3. Create specialized agents for visual tasks
4. Implement caching for multimodal content
5. Optimize for performance and resource usage

## Example Use Cases

### Visual Analysis Agent

```yaml
name: image-analyzer
description: An agent that can analyze images and provide detailed descriptions
model:
  base: claude-3-opus
  vision: true
capabilities:
  - Analyze image content and provide detailed descriptions
  - Identify objects, people, and scenes
  - Extract text from images
  - Detect potential safety concerns
```

### Design Assistant

```yaml
name: design-assistant
description: An agent that helps with UI/UX design tasks
model:
  base: gpt-4-vision
  vision: true
capabilities:
  - Critique UI designs
  - Suggest improvements to layouts
  - Generate wireframe descriptions
  - Provide feedback on color and typography
tools:
  - design_principles:
      purpose: For accessing design guidelines and principles
  - color_analyzer:
      purpose: For analyzing color schemes
```

### Medical Imaging Assistant

```yaml
name: medical-assistant
description: An agent that helps analyze medical images
model:
  base: claude-3-opus
  vision: true
compliance:
  regulatory_frameworks:
    - hipaa
  security_measures:
    - encryption_at_rest
    - secure_communications
capabilities:
  - Analyze X-ray, MRI, and CT images
  - Highlight potential areas of concern
  - Provide educational explanations
  - Assist with medical documentation
```

## Future Enhancements

- **Video Processing**: Support for analyzing video content
- **Audio Modality**: Add support for speech and audio processing
- **Real-time Camera Input**: Enable agents to process live camera feeds
- **Image Generation**: Integrate with image generation models like DALL-E
- **Multimodal RAG**: Implement retrieval-augmented generation for multimodal content
- **Collaborative Visual Work**: Enable multiple agents to work together on visual tasks

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| API costs for vision models | Higher operational costs | Implement efficient image preprocessing and caching |
| Privacy concerns with image data | Security and compliance issues | Add image anonymization and strict data handling policies |
| Provider API changes | Integration breakage | Create robust fallback mechanisms and monitoring |
| Performance issues with large media files | User experience degradation | Implement progressive loading and optimization strategies |
| Inappropriate content generation | Ethical and reputation risks | Apply content filtering and safety measures |

## Success Metrics

- **Provider Coverage**: Support for at least 3 major multimodal providers
- **Performance**: Average response time under 3 seconds for typical image analysis
- **Capability**: Successful handling of 95% of common image formats and scenes
- **User Adoption**: 30% of agents created using multimodal capabilities within 3 months
- **Error Rate**: Less than 5% failure rate for multimodal processing requests

---

**Next Steps:**
1. Create initial data structures for multimodal content
2. Implement provider interfaces for Claude and OpenAI
3. Develop prototype CLI commands for testing
4. Create example agents with multimodal capabilities 