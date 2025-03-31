# Multimodal Support

SentinelStacks provides robust support for multimodal interactions, allowing agents to process and analyze various types of content, including images. This document explains how to use the multimodal capabilities.

## Quick Start

### Analyzing Images with the CLI

You can use the `analyze-image` command to quickly analyze images without creating a long-running agent:

```bash
# Analyze an image with Claude
sentinel multimodal analyze-image --provider claude --image path/to/image.jpg "What's in this image?"

# Analyze an image with OpenAI
sentinel multimodal analyze-image --provider openai --image path/to/image.jpg "Describe this image in detail"

# Analyze multiple images
sentinel multimodal analyze-image --provider claude --image image1.jpg,image2.jpg "Compare these images"
```

### Interactive Chat with Images

For a more interactive experience, you can use the `chat` command which supports both text and image inputs:

```bash
# Start a chat session with default settings
sentinel chat

# Start a chat session with a specific provider and model
sentinel chat --provider claude --model claude-3-opus-20240229

# Start a chat session with initial images
sentinel chat --provider openai --images image1.jpg,image2.jpg
```

During the chat session, you can upload images using the `/image` command:

```
You: /image path/to/image.jpg
Enter a question about the image(s): What's in this image?
```

## Supported Providers

The following providers support multimodal capabilities:

| Provider | Multimodal Support | Supported Models |
|----------|-------------------|------------------|
| Claude   | Yes               | claude-3-opus, claude-3-sonnet, claude-3-haiku |
| OpenAI   | Yes               | gpt-4o, gpt-4-vision-preview |
| Ollama   | Yes (depends on model) | llava, bakllava, llava-llama3 |
| Mock     | Yes (for testing) | mock-model |

## Running Agents with Multimodal Input

You can run agents with multimodal input using the `run` command with the `--image` flag:

```bash
# Run an agent with an image
sentinel run my-agent --image path/to/image.jpg

# Run an agent with multiple images
sentinel run my-agent --image image1.jpg,image2.jpg
```

## Programmatic Usage

If you're building applications that use SentinelStacks for multimodal processing, you can use the runtime API:

```go
import (
    "github.com/sentinelstacks/sentinel/internal/multimodal"
    "github.com/sentinelstacks/sentinel/internal/runtime"
)

// Create a multimodal agent
agent, err := rt.CreateMultimodalAgent("my-agent", "claude:latest", "claude-3-opus", "claude", apiKey, "")

// Process text input
response, err := agent.ProcessTextInput(ctx, "Hello, world!")

// Process multimodal input
input := multimodal.NewInput()
input.AddText("What's in this image?")
input.AddImage(imageData, "image/jpeg")
output, err := agent.ProcessMultimodalInput(ctx, input.Contents)
```

## Best Practices

1. **Image Quality**: Provide clear, high-quality images for the best results.
2. **Specific Questions**: Ask specific questions about the images rather than general ones.
3. **Multiple Images**: When analyzing multiple images, be clear about which image you're referring to.
4. **Image Format**: JPEG, PNG, GIF, WEBP, and BMP formats are supported. JPEG and PNG are recommended.
5. **Image Size**: Large images (>10MB) may be resized or rejected by some providers.

## Troubleshooting

- **Error: Provider does not support multimodal**: Ensure you're using a model that supports multimodal input.
- **Error: Failed to read image file**: Check that the image path is correct and the file exists.
- **Poor analysis quality**: Try using a different provider or a more advanced model. 