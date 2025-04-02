package multimodal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	// TODO: Add these dependencies to go.mod with:
	// go get github.com/briandowns/spinner
	// go get github.com/fatih/color
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/shim"
)

// NewMultimodalCmd creates a new multimodal command
func NewMultimodalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multimodal",
		Short: "Work with multimodal features",
		Long:  `Commands for working with multimodal features such as image and video analysis.`,
	}

	// Add the analyze-image subcommand
	cmd.AddCommand(newAnalyzeImageCmd())

	return cmd
}

// newAnalyzeImageCmd creates a new analyze-image command
func newAnalyzeImageCmd() *cobra.Command {
	var (
		imagePath   string
		prompt      string
		apiKey      string
		provider    string
		model       string
		temperature float64
		maxTokens   int
		noStreaming bool
		outputFile  string
	)

	cmd := &cobra.Command{
		Use:   "analyze-image",
		Short: "Analyze an image using a multimodal model",
		Long: `Analyze an image using a multimodal language model.
The command sends the image to the specified provider (Claude, OpenAI, Ollama, etc.) 
along with a prompt to analyze the image content.

Examples:
  sentinel multimodal analyze-image --image /path/to/image.jpg --prompt "Describe this image in detail"
  sentinel multimodal analyze-image --image /path/to/image.jpg --provider openai --model gpt-4-vision-preview
  sentinel multimodal analyze-image --image /path/to/image.jpg --provider ollama --model llava`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required parameters
			if imagePath == "" {
				return fmt.Errorf("image path is required")
			}

			// Default prompt if not provided
			if prompt == "" {
				prompt = "Describe this image in detail."
			}

			// Default provider if not provided
			if provider == "" {
				provider = "claude" // Default to Claude
			}

			// Create a map to hold additional parameters
			params := make(map[string]interface{})

			// Handle API key from environment if not provided
			if apiKey == "" {
				switch strings.ToLower(provider) {
				case "claude":
					apiKey = os.Getenv("ANTHROPIC_API_KEY")
					if apiKey == "" {
						return fmt.Errorf("API key not provided and ANTHROPIC_API_KEY environment variable not set")
					}
				case "openai":
					apiKey = os.Getenv("OPENAI_API_KEY")
					if apiKey == "" {
						return fmt.Errorf("API key not provided and OPENAI_API_KEY environment variable not set")
					}
				case "ollama":
					// Ollama typically runs locally without an API key, so we allow empty API key
					// Check if endpoint is provided in the environment
					endpoint := os.Getenv("OLLAMA_ENDPOINT")
					if endpoint != "" {
						// Set endpoint in the parameters
						params["endpoint"] = endpoint
					}
				default:
					return fmt.Errorf("unknown provider: %s", provider)
				}
			}

			// Set up model based on provider if not specified
			if model == "" {
				switch strings.ToLower(provider) {
				case "claude":
					model = "claude-3-opus-20240229"
				case "openai":
					model = "gpt-4-vision-preview"
				case "ollama":
					model = "llava"
				}
			}

			// Default temperature and max tokens if not provided
			if temperature == 0 {
				temperature = 0.7
			}
			if maxTokens == 0 {
				maxTokens = 4096
			}

			// Create shim configuration
			config := shim.Config{
				Provider: provider,
				Model:    model,
				APIKey:   apiKey,
				Timeout:  60 * time.Second,
			}

			// Set endpoint if provided
			if endpoint, ok := params["endpoint"].(string); ok && endpoint != "" {
				config.Endpoint = endpoint
			}

			// Create the shim
			shimProvider, err := shim.ShimFactory(provider, config.Endpoint, config.APIKey, config.Model)
			if err != nil {
				return fmt.Errorf("failed to create shim: %w", err)
			}
			defer shimProvider.Close()

			// Check if the shim supports multimodal
			if !shimProvider.SupportsMultimodal() {
				return fmt.Errorf("provider %s with model %s does not support multimodal inputs", provider, model)
			}

			// Read image file
			imgData, err := os.ReadFile(imagePath)
			if err != nil {
				return fmt.Errorf("failed to read image file: %w", err)
			}

			// Create multimodal input
			input := multimodal.NewInput()
			input.AddText(prompt)
			input.AddImage(imgData, filepath.Ext(imagePath)[1:]) // Remove the dot from extension
			input.Temperature = temperature
			input.MaxTokens = maxTokens

			// Create output file if specified
			var outputWriter io.Writer
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer f.Close()
				outputWriter = f
			} else {
				outputWriter = cmd.OutOrStdout()
			}

			// Show processing message
			fmt.Fprintln(cmd.OutOrStderr(), "Analyzing image...")

			// Handle streaming or non-streaming
			if noStreaming {
				// Non-streaming mode
				result, err := shimProvider.MultimodalCompletion(input, 120*time.Second)
				if err != nil {
					return fmt.Errorf("failed to generate response: %w", err)
				}

				// Get text from the result
				for _, content := range result.Contents {
					if content.Type == multimodal.MediaTypeText {
						fmt.Fprintln(outputWriter, content.Text)
					}
				}
			} else {
				// Streaming mode
				resultCh, err := shimProvider.StreamMultimodalCompletion(context.Background(), input)
				if err != nil {
					return fmt.Errorf("failed to start streaming: %w", err)
				}

				// Process the streaming response
				bold := color.New(color.Bold)
				bold.Fprintln(outputWriter, "Analysis:")
				for chunk := range resultCh {
					if chunk.Error != nil {
						return fmt.Errorf("error during streaming: %w", chunk.Error)
					}

					if chunk.Content.Type == multimodal.MediaTypeText {
						fmt.Fprint(outputWriter, chunk.Content.Text)
					}
				}
				fmt.Fprintln(outputWriter) // Add a newline at the end
			}

			// Save the analyzed image to a local directory for reference
			if err := saveAnalyzedImage(imagePath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to save analyzed image: %v\n", err)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&imagePath, "image", "i", "", "Path to the image file to analyze (required)")
	cmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Prompt to use for image analysis (default: \"Describe this image in detail.\")")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for the provider (defaults to environment variable based on provider)")
	cmd.Flags().StringVarP(&provider, "provider", "r", "", "Provider to use (claude, openai, ollama) (default: claude)")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model to use (defaults based on provider)")
	cmd.Flags().Float64VarP(&temperature, "temperature", "t", 0, "Temperature for generation (default: 0.7)")
	cmd.Flags().IntVarP(&maxTokens, "max-tokens", "x", 0, "Maximum tokens for generation (default: 4096)")
	cmd.Flags().BoolVarP(&noStreaming, "no-streaming", "n", false, "Disable streaming response (default: false)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save the analysis (default: stdout)")

	// Mark required flags
	cmd.MarkFlagRequired("image")

	return cmd
}

// saveAnalyzedImage saves a copy of the analyzed image to a local directory for reference
func saveAnalyzedImage(imagePath string) error {
	// Create the analyzed images directory in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	analyzedDir := filepath.Join(homeDir, ".sentinel", "analyzed_images")
	if err := os.MkdirAll(analyzedDir, 0755); err != nil {
		return err
	}

	// Read the source image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	// Create a filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Base(imagePath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	newFilename := fmt.Sprintf("%s_%s%s", name, timestamp, ext)
	destPath := filepath.Join(analyzedDir, newFilename)

	// Write the file
	return os.WriteFile(destPath, imageData, 0644)
}
