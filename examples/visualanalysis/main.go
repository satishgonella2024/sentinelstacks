package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/shim"
)

func main() {
	// Parse command-line flags
	imageFile := flag.String("image", "", "Path to the image file to analyze")
	prompt := flag.String("prompt", "What can you tell me about this image?", "Prompt to send with the image")
	provider := flag.String("provider", "claude", "LLM provider to use (claude, openai)")
	model := flag.String("model", "claude-3-opus-20240229", "Model to use")
	apiKey := flag.String("api-key", "", "API key for the LLM provider")
	endpoint := flag.String("endpoint", "", "API endpoint for the LLM provider (optional)")
	flag.Parse()

	// Check required flags
	if *imageFile == "" {
		log.Fatal("Missing required --image flag")
	}

	if *apiKey == "" {
		// Try environment variable
		*apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if *provider == "openai" {
			*apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if *apiKey == "" {
			log.Fatal("Missing API key. Use --api-key flag or set ANTHROPIC_API_KEY/OPENAI_API_KEY environment variable")
		}
	}

	// Load the image file
	fmt.Printf("Loading image from %s...\n", *imageFile)
	imageContent, err := loadImage(*imageFile)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// Create and initialize the shim
	fmt.Printf("Initializing %s provider with model %s...\n", *provider, *model)
	shimInstance, err := createShim(*provider, *model, *apiKey, *endpoint)
	if err != nil {
		log.Fatalf("Failed to create shim: %v", err)
	}
	defer shimInstance.Close()

	// Check if multimodal is supported
	if !shimInstance.SupportsMultimodal() {
		log.Fatalf("Provider %s with model %s does not support multimodal inputs", *provider, *model)
	}

	// Prepare multimodal input
	input := buildInput(*prompt, imageContent)

	// Generate a response
	fmt.Println("Generating response...")
	output, err := shimInstance.GenerateMultimodal(context.Background(), input)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	// Print the response
	fmt.Printf("\n--- ANALYSIS RESULT ---\n\n")
	fmt.Println(output.GetText())
	fmt.Printf("\n--- END OF ANALYSIS ---\n\n")

	// Print token usage if available
	if output.UsedTokens > 0 {
		fmt.Printf("Used %d tokens\n", output.UsedTokens)
	}
}

// loadImage loads and processes an image file
func loadImage(filePath string) (*multimodal.Content, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Set image processing options (resize to reasonable dimensions)
	options := multimodal.ImageProcessingOptions{
		MaxWidth:  1024,
		MaxHeight: 1024,
		Compress:  true,
		Quality:   85,
	}

	// Load the image
	content, err := multimodal.LoadImageFromFile(filePath, &options)
	if err != nil {
		return nil, err
	}

	// Add the filename as alt text if not already set
	if content.Text == "" {
		content.Text = filepath.Base(filePath)
	}

	return content, nil
}

// createShim creates and initializes an LLM provider shim
func createShim(provider, model, apiKey, endpoint string) (shim.Shim, error) {
	config := shim.Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
	}

	if endpoint != "" {
		config.Endpoint = endpoint
	}

	return shim.CreateShim(provider, model, config)
}

// buildInput builds a multimodal input with text and image
func buildInput(text string, imageContent *multimodal.Content) *multimodal.Input {
	opts := map[string]interface{}{
		"max_tokens":    2000,
		"temperature":   0.7,
		"system_prompt": "You are a visual analysis AI assistant. Your task is to analyze images and provide detailed, accurate descriptions. Pay attention to details like objects, people, scenes, text, colors, and composition. Be thorough but concise in your analysis.",
	}

	return multimodal.BuildMultimodalInput(text, imageContent, opts)
}
