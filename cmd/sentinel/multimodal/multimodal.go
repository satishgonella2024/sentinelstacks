package multimodal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim"
	"github.com/spf13/cobra"
)

var (
	imagePath   string
	imagePrompt string
	provider    string
	model       string
	apiKey      string
	stream      bool
)

// NewMultimodalCmd creates a new multimodal command
func NewMultimodalCmd() *cobra.Command {
	// multimodalCmd represents the multimodal command
	multimodalCmd := &cobra.Command{
		Use:   "multimodal",
		Short: "Commands for multimodal interactions",
		Long:  `Commands for interacting with multimodal capabilities of LLMs.`,
	}

	// Add analyze-image command to multimodal
	multimodalCmd.AddCommand(newAnalyzeImageCmd())

	return multimodalCmd
}

// newAnalyzeImageCmd creates a new analyze-image command
func newAnalyzeImageCmd() *cobra.Command {
	// analyzeImageCmd represents the analyze-image command
	analyzeImageCmd := &cobra.Command{
		Use:   "analyze-image",
		Short: "Analyze an image with an LLM",
		Long:  `Send an image to an LLM for analysis with an optional text prompt.`,
		RunE:  runAnalyzeImage,
	}

	// Add flags to analyze-image command
	analyzeImageCmd.Flags().StringVarP(&imagePath, "image", "i", "", "Path to image file to analyze")
	analyzeImageCmd.Flags().StringVarP(&imagePrompt, "prompt", "p", "Analyze this image in detail.", "Text prompt to accompany the image")
	analyzeImageCmd.Flags().StringVarP(&provider, "provider", "r", "claude", "LLM provider to use (claude, openai)")
	analyzeImageCmd.Flags().StringVarP(&model, "model", "m", "", "Model to use (defaults to provider's default)")
	analyzeImageCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for the provider")
	analyzeImageCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Stream the response")

	// Mark image path as required
	analyzeImageCmd.MarkFlagRequired("image")

	return analyzeImageCmd
}

// runAnalyzeImage is the function for the analyze-image command
func runAnalyzeImage(cmd *cobra.Command, args []string) error {
	// Check if image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", imagePath)
	}

	// Create a base file name for logs
	baseName := filepath.Base(imagePath)
	ext := filepath.Ext(baseName)
	baseFileName := strings.TrimSuffix(baseName, ext)

	// Create a timestamp for log files
	timestamp := time.Now().Format("20060102_150405")

	// Get API key from environment if not provided
	if apiKey == "" {
		if provider == "claude" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		} else if provider == "openai" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		// Check if we have an API key
		if apiKey == "" {
			return fmt.Errorf("no API key provided and none found in environment variables")
		}
	}

	// Initialize the provider
	fmt.Printf("Initializing %s provider...\n", provider)

	// Initialize the shim
	shimConfig := shim.Config{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
	}

	llm, err := shim.CreateShim(provider, model, shimConfig)
	if err != nil {
		return fmt.Errorf("failed to create shim: %v", err)
	}

	// Load image
	fmt.Println("Loading image...")

	// Create options
	options := multimodal.DefaultImageOptions()

	imgContent, err := multimodal.LoadImageFromFile(imagePath, &options)
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}

	// Add alt text if available
	if imgContent.Text == "" {
		imgContent.Text = "Image uploaded by user for analysis"
	}

	// Build multimodal input
	input := multimodal.NewInput()
	input.AddText(imagePrompt)
	input.Contents = append(input.Contents, imgContent)
	input.SetMaxTokens(4096)
	input.SetTemperature(0.7)
	input.SetStream(stream)

	// Set a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Process the image
	if stream {
		// Stream response
		fmt.Println("Sending image for analysis (streaming mode)...")
		streamCh, err := llm.StreamMultimodal(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to stream multimodal response: %v", err)
		}

		// Create a logfile for the response
		logFilePath := fmt.Sprintf("%s_%s_response.txt", baseFileName, timestamp)
		logFile, err := os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to create log file: %v", err)
		}
		defer logFile.Close()

		// Process stream chunks
		for chunk := range streamCh {
			if chunk.Error != nil {
				return fmt.Errorf("error in stream: %v", chunk.Error)
			}

			if chunk.Content.Type == multimodal.MediaTypeText {
				fmt.Print(chunk.Content.Text)
				logFile.WriteString(chunk.Content.Text)
			}
		}

		fmt.Println("\n\nResponse saved to:", logFilePath)
	} else {
		// Get full response
		fmt.Println("Sending image for analysis...")
		output, err := llm.GenerateMultimodal(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to get multimodal response: %v", err)
		}

		// Extract text
		text := multimodal.ExtractTextFromOutput(output)

		// Create a logfile for the response
		logFilePath := fmt.Sprintf("%s_%s_response.txt", baseFileName, timestamp)
		logFile, err := os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to create log file: %v", err)
		}
		defer logFile.Close()

		// Write response
		fmt.Println(text)
		logFile.WriteString(text)

		fmt.Println("\nResponse saved to:", logFilePath)

		// Extract any images from response (for future use)
		images := multimodal.ExtractImagesFromOutput(output)
		for i, img := range images {
			imgFileName := fmt.Sprintf("%s_%s_response_image_%d%s", baseFileName, timestamp, i, ".jpg")
			err := multimodal.SaveImageToFile(img, imgFileName)
			if err != nil {
				fmt.Printf("Warning: Failed to save response image: %v\n", err)
				continue
			}
			fmt.Printf("Response image saved to: %s\n", imgFileName)
		}
	}

	return nil
}
