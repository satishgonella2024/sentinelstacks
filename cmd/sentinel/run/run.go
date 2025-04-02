package run

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/registry"
	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
)

// NewRunCmd creates the run command
func NewRunCmd() *cobra.Command {
	var (
		env         []string
		interactive bool
		llmProvider string
		llmEndpoint string
		llmModel    string
		timeout     time.Duration
		imageFile   string
	)

	runCmd := &cobra.Command{
		Use:   "run [image_name]",
		Short: "Run a Sentinel Agent from an image",
		Long:  `Run a Sentinel Agent from a previously built image or from a registry`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse environment variables
			envMap := parseEnvironmentVariables(env)
			
			// Parse the image name
			imageName, imageTag := parseImageName(args[0])
			
			// Load the image
			image, err := loadImage(imageName, imageTag)
			if err != nil {
				return err
			}
			
			// Handle multimodal content if provided
			mmContent, err := loadMultimodalContent(imageFile)
			if err != nil {
				return err
			}
			
			// Configure LLM settings
			llmConfig, err := configureLLM(llmProvider, llmEndpoint, llmModel, image, mmContent != nil)
			if err != nil {
				return err
			}
			
			// Print configuration
			printRunConfiguration(imageName, imageTag, llmConfig, interactive, mmContent, timeout, envMap, &image.Definition)
			
			// Also print the actual endpoint being used for debugging
			fmt.Printf("Debug - Using LLM endpoint: %s\n", llmConfig.Endpoint)
			
			// Create agent runtime
			rt, err := runtime.GetRuntime()
			if err != nil {
				return fmt.Errorf("failed to get runtime: %w", err)
			}
			
			// Create the agent
			agent, err := rt.CreateAgent(image.Definition.Name, fmt.Sprintf("%s:%s", imageName, imageTag), llmConfig.Model)
			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}
			
			// Configure API key
			apiKey := getAPIKey(llmConfig.Provider)
			
			// Create multimodal agent
			mmAgent, err := rt.CreateMultimodalAgent(
				image.Definition.Name, 
				fmt.Sprintf("%s:%s", imageName, imageTag), 
				llmConfig.Model,
				llmConfig.Provider,
				apiKey,
				llmConfig.Endpoint,
			)
			if err != nil {
				return fmt.Errorf("failed to create multimodal agent: %w", err)
			}
			
			// Set up context with cancellation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			
			// Handle termination signals
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			
			go func() {
				<-signalCh
				fmt.Println("\nReceived termination signal, shutting down...")
				cancel()
			}()
			
			// Run the agent
			if interactive {
				return runInteractiveMode(ctx, mmAgent, mmContent)
			} else {
				// Start the agent in background mode
				if err := rt.StartAgent(agent.ID); err != nil {
					return fmt.Errorf("failed to start agent: %w", err)
				}
				
				fmt.Printf("Agent started with ID: %s\n", agent.ID)
				fmt.Println("Use 'sentinel logs " + agent.ID + "' to view logs")
				fmt.Println("Use 'sentinel stop " + agent.ID + "' to stop the agent")
				
				return nil
			}
		},
	}

	runCmd.Flags().StringSliceVarP(&env, "env", "e", []string{}, "Set environment variables")
	runCmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Run in interactive mode")
	runCmd.Flags().StringVar(&llmProvider, "llm", "", "Override the LLM provider (claude, openai, ollama)")
	runCmd.Flags().StringVar(&llmEndpoint, "llm-endpoint", "", "LLM provider endpoint URL")
	runCmd.Flags().StringVar(&llmModel, "llm-model", "", "Override the LLM model")
	runCmd.Flags().DurationVar(&timeout, "timeout", 60*time.Second, "Timeout for the agent run (e.g. 1h, 30m)")
	runCmd.Flags().StringVar(&imageFile, "image", "", "Path to an image file to include as multimodal input")

	return runCmd
}

// Helper functions

// parseEnvironmentVariables parses environment variables from a slice of strings
func parseEnvironmentVariables(env []string) map[string]interface{} {
	envMap := make(map[string]interface{})
	
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			// Try to convert value to appropriate type
			if value == "true" {
				envMap[key] = true
			} else if value == "false" {
				envMap[key] = false
			} else if intVal, err := parseInt(value); err == nil {
				envMap[key] = intVal
			} else if floatVal, err := parseFloat(value); err == nil {
				envMap[key] = floatVal
			} else {
				envMap[key] = value
			}
		}
	}
	
	return envMap
}

// parseInt tries to parse a string as an integer
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat tries to parse a string as a float
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseImageName parses an image name into name and tag parts
func parseImageName(imageName string) (string, string) {
	// Add :latest tag if no tag specified
	if !strings.Contains(imageName, ":") {
		imageName = imageName + ":latest"
	}

	// Split name and tag
	parts := strings.SplitN(imageName, ":", 2)
	name := parts[0]
	tag := "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}
	
	return name, tag
}

// loadImage loads an image from the registry
func loadImage(name, tag string) (*registry.Image, error) {
	// Get the registry
	reg, err := registry.GetLocalRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}

	// Load the image from the registry
	image, err := reg.Get(name, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	
	return image, nil
}

// loadMultimodalContent loads multimodal content from a file
func loadMultimodalContent(imageFile string) (*multimodal.Content, error) {
	if imageFile == "" {
		return nil, nil
	}
	
	// Validate that the image file exists
	if _, err := os.Stat(imageFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", imageFile)
	}

	// Load the image file
	imgData, err := os.ReadFile(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Create multimodal content
	mmContent := multimodal.NewImageContent(imgData, getContentType(imageFile))
	mmContent.Text = filepath.Base(imageFile) // Set alt text to filename

	fmt.Printf("Loaded input image: %s (%s, %d bytes)\n",
		filepath.Base(imageFile),
		mmContent.MimeType,
		len(mmContent.Data))
		
	return mmContent, nil
}

// LLMConfig represents the configuration for an LLM
type LLMConfig struct {
	Provider string
	Endpoint string
	Model    string
}

// configureLLM configures the LLM settings based on flags, config, and image
func configureLLM(provider, endpoint, model string, image *registry.Image, isMultimodal bool) (*LLMConfig, error) {
	config := &LLMConfig{}
	
	// Configure provider
	config.Provider = provider
	if config.Provider == "" {
		// Try to get from image definition
		if image.Definition.BaseModel != "" {
			// Extract provider from model name (simplistic approach)
			if strings.HasPrefix(image.Definition.BaseModel, "claude") {
				config.Provider = "claude"
			} else if strings.HasPrefix(image.Definition.BaseModel, "llama") || 
				   strings.HasPrefix(image.Definition.BaseModel, "mistral") {
				config.Provider = "ollama"
			} else if strings.HasPrefix(image.Definition.BaseModel, "gpt") {
				config.Provider = "openai"
			}
		}
		
		// If still not set, use from config
		if config.Provider == "" {
			config.Provider = viper.GetString("llm.provider")
			
			// Use default if not set in config
			if config.Provider == "" {
				config.Provider = "claude" // Default provider
			}
		}
	}
	
	// Configure endpoint
	config.Endpoint = endpoint
	if config.Endpoint == "" {
		config.Endpoint = viper.GetString("llm.endpoint")
		
		// If still not set, use provider-specific default
		if config.Endpoint == "" {
			if config.Provider == "ollama" {
				// Check if we have a custom Ollama endpoint in config
				customEndpoint := viper.GetString("ollama.endpoint")
				if customEndpoint != "" {
					config.Endpoint = customEndpoint
				} else {
					config.Endpoint = "http://localhost:11434/api/generate"
				}
			}
			// For Claude and OpenAI, use defaults in the shim
		}
	}
	
	// Configure model
	config.Model = model
	if config.Model == "" {
		// Try to get from image definition
		if image.Definition.BaseModel != "" {
			config.Model = image.Definition.BaseModel
		} else {
			// Try to get from config
			config.Model = viper.GetString("llm.model")
			
			// Set appropriate default based on provider
			if config.Model == "" {
				if config.Provider == "ollama" {
					config.Model = "llama3"
				} else if config.Provider == "claude" {
					config.Model = "claude-3-5-sonnet-20240627"
				} else if config.Provider == "openai" {
					config.Model = "gpt-4"
				}
			}
		}
	}
	
	// If multimodal content is provided, ensure the model supports it
	if isMultimodal {
		switch config.Provider {
		case "claude":
			// Ensure we're using a Claude 3 model that supports multimodal
			if !strings.HasPrefix(config.Model, "claude-3") {
				fmt.Println("Warning: Switching to claude-3-opus-20240229 for multimodal support")
				config.Model = "claude-3-opus-20240229"
			}
		case "openai":
			// Ensure we're using GPT-4 Vision that supports multimodal
			if !strings.Contains(config.Model, "vision") {
				fmt.Println("Warning: Switching to gpt-4-vision-preview for multimodal support")
				config.Model = "gpt-4-vision-preview"
			}
		case "ollama":
			// Ensure we're using a model that supports multimodal like llava
			if !strings.Contains(config.Model, "llava") &&
				!strings.Contains(config.Model, "bakllava") &&
				!strings.Contains(config.Model, "moondream") {
				fmt.Println("Warning: Switching to llava for multimodal support")
				config.Model = "llava"
			}
		}
	}
	
	return config, nil
}

// getAPIKey gets the API key for the specified provider
func getAPIKey(provider string) string {
	// Try generic key first
	apiKey := viper.GetString("llm.api_key")
	
	// If not set, try provider-specific key
	if apiKey == "" {
		apiKey = viper.GetString(fmt.Sprintf("%s.api_key", provider))
	}
	
	return apiKey
}

// printRunConfiguration prints the run configuration
func printRunConfiguration(name, tag string, llmConfig *LLMConfig, interactive bool, 
						  mmContent *multimodal.Content, timeout time.Duration, 
						  envMap map[string]interface{}, def *registry.ImageDefinition) {
	fmt.Printf("Running agent from image: %s:%s\n", name, tag)
	fmt.Printf("Using LLM provider: %s\n", llmConfig.Provider)
	if llmConfig.Endpoint != "" {
		fmt.Printf("LLM endpoint: %s\n", llmConfig.Endpoint)
	}
	fmt.Printf("LLM model: %s\n", llmConfig.Model)
	fmt.Printf("Interactive mode: %v\n", interactive)
	if mmContent != nil {
		fmt.Println("Mode: Multimodal (image input provided)")
	}
	fmt.Printf("Timeout: %s\n", timeout.String())

	if len(envMap) > 0 {
		fmt.Println("Environment variables:")
		for k, v := range envMap {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// Show agent details
	fmt.Println("\nAgent details:")
	fmt.Printf("  Name: %s\n", def.Name)
	fmt.Printf("  Description: %s\n", def.Description)
	if len(def.Capabilities) > 0 {
		fmt.Println("  Capabilities:")
		for _, cap := range def.Capabilities {
			fmt.Printf("    - %s\n", cap)
		}
	}
}

// runInteractiveMode runs the agent in interactive mode
func runInteractiveMode(ctx context.Context, agent *runtime.MultimodalAgent, mmContent *multimodal.Content) error {
	fmt.Println("\n=== Interactive Session Started ===")
	fmt.Println("Type 'exit' to end the session, 'help' for commands")
	
	// Create welcome message
	fmt.Printf("\nAgent> I'm %s. How can I help you today?\n", agent.Name)
	
	// If there's an image, process it first
	if mmContent != nil {
		fmt.Println("Agent> I see you've provided an image. Let me analyze it.")
		
		// Create multimodal input with the image
		input := multimodal.NewInput()
		input.AddContent(mmContent)
		input.AddText("Tell me about this image.")
		
		// Process the image
		output, err := agent.ProcessMultimodalInput(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to process image: %w", err)
		}
		
		// Print the response
		fmt.Printf("Agent> %s\n", extractTextFromOutput(output))
	}
	
	// Create scanner for user input
	scanner := bufio.NewScanner(os.Stdin)
	
	// Main interaction loop
	for {
		fmt.Print("\nUser> ")
		if !scanner.Scan() {
			break
		}
		
		// Get user input
		userInput := scanner.Text()
		if userInput == "" {
			continue
		}
		
		// Handle commands
		if userInput == "exit" {
			fmt.Println("\nEnding session.")
			break
		}
		
		if userInput == "help" {
			printHelp()
			continue
		}
		
		// Send request to the agent
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		
		// Process text input and print response
		fmt.Print("Agent> ")
		response, err := agent.ProcessTextInput(ctx, userInput)
		if err != nil {
			cancel()
			fmt.Printf("Error: %v\n", err)
			continue
		}
		
		// Print the response
		fmt.Println(response)
		
		// Cancel the context
		cancel()
	}
	
	fmt.Println("\n=== Interactive Session Ended ===")
	return nil
}

// printHelp prints available commands
func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  exit - End the session")
	fmt.Println("  help - Show this help message")
}

// getContentType returns the MIME type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// extractTextFromOutput extracts text content from a multimodal output
func extractTextFromOutput(output *multimodal.Output) string {
	var text string
	
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			text += content.Text
		}
	}
	
	return text
}
