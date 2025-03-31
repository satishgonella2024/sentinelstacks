package run

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinelstacks/sentinel/pkg/agent"
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
	)

	runCmd := &cobra.Command{
		Use:   "run [image_name]",
		Short: "Run a Sentinel Agent from an image",
		Long:  `Run a Sentinel Agent from a previously built image or from a registry`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName := args[0]

			// Add :latest tag if no tag specified
			if !strings.Contains(imageName, ":") {
				imageName = imageName + ":latest"
			}

			// Parse environment variables
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
				} else {
					return fmt.Errorf("invalid environment variable format: %s", e)
				}
			}

			// Load the image
			image, err := loadImage(imageName)
			if err != nil {
				return fmt.Errorf("failed to load image: %w", err)
			}

			// Set up the LLM configuration

			// If LLM provider was not specified, use the one from the image or config
			if llmProvider == "" {
				// First try to get from image definition
				if image.Definition.BaseModel != "" {
					// Extract provider from model name (simplistic approach)
					if strings.HasPrefix(image.Definition.BaseModel, "claude") {
						llmProvider = "claude"
					} else if strings.HasPrefix(image.Definition.BaseModel, "llama") {
						llmProvider = "ollama"
					} else if strings.HasPrefix(image.Definition.BaseModel, "gpt") {
						llmProvider = "openai"
					} else {
						// Default to provider from config
						llmProvider = viper.GetString("llm.provider")
					}
				}

				// If still not set, use default from config
				if llmProvider == "" {
					llmProvider = viper.GetString("llm.provider")
					if llmProvider == "" {
						llmProvider = "claude" // Ultimate default
					}
				}
			}

			// If LLM endpoint was not specified, use the one from the config
			if llmEndpoint == "" {
				llmEndpoint = viper.GetString("llm.endpoint")

				// Set appropriate default based on provider
				if llmEndpoint == "" {
					if llmProvider == "ollama" {
						// Check if we have a custom Ollama endpoint in the environment or config
						customEndpoint := viper.GetString("ollama.endpoint")
						if customEndpoint != "" {
							llmEndpoint = customEndpoint
						} else {
							llmEndpoint = "http://localhost:11434"
						}
					}
				}
			}

			// If LLM model was not specified, use the one from the image or config
			if llmModel == "" {
				// First try to get from image definition
				if image.Definition.BaseModel != "" {
					llmModel = image.Definition.BaseModel
				} else {
					// Try to get from config
					llmModel = viper.GetString("llm.model")

					// Set appropriate default based on provider
					if llmModel == "" {
						if llmProvider == "ollama" {
							llmModel = "llama3"
						} else if llmProvider == "claude" {
							llmModel = "claude-3.7-sonnet"
						}
					}
				}
			}

			// Get the API key from the config (will be used in future implementation)
			// apiKey := viper.GetString("llm.api_key")

			fmt.Printf("Running agent from image: %s\n", imageName)
			fmt.Printf("Using LLM provider: %s\n", llmProvider)
			if llmEndpoint != "" {
				fmt.Printf("LLM endpoint: %s\n", llmEndpoint)
			}
			fmt.Printf("LLM model: %s\n", llmModel)
			fmt.Printf("Interactive mode: %v\n", interactive)
			if timeout > 0 {
				fmt.Printf("Timeout: %s\n", timeout.String())
			} else {
				fmt.Println("Timeout: none")
			}

			if len(envMap) > 0 {
				fmt.Println("Environment variables:")
				for k, v := range envMap {
					fmt.Printf("  %s: %v\n", k, v)
				}
			}

			// Show agent details
			fmt.Println("\nAgent details:")
			fmt.Printf("  Name: %s\n", image.Definition.Name)
			fmt.Printf("  Description: %s\n", image.Definition.Description)
			if len(image.Definition.Capabilities) > 0 {
				fmt.Println("  Capabilities:")
				for _, cap := range image.Definition.Capabilities {
					fmt.Printf("    - %s\n", cap)
				}
			}

			// Simulate agent initialization and running
			fmt.Println("\nInitializing agent runtime...")
			fmt.Printf("Connecting to %s provider...\n", llmProvider)
			fmt.Println("Setting up agent state...")

			// Generate a fake agent ID
			agentID := "agent_" + fmt.Sprintf("%x", time.Now().UnixNano())

			// TODO: Implement actual agent runtime
			// This would include creating an LLM shim and initializing the agent

			if interactive {
				fmt.Println("\n=== Interactive Session Started ===")
				fmt.Println("Type 'exit' to end the session")
				fmt.Println("Agent> Hello! I'm ready to assist you.")
				fmt.Println("User> (Type your message here)")
				fmt.Println("\nSimulating interactive session...")
				fmt.Println("=== Interactive Session Ended ===")
			} else {
				fmt.Println("Agent running in background mode")
			}

			fmt.Printf("Agent started with ID: %s\n", agentID)
			fmt.Println("Use 'sentinel logs " + agentID + "' to view logs")
			fmt.Println("Use 'sentinel stop " + agentID + "' to stop the agent")

			return nil
		},
	}

	runCmd.Flags().StringSliceVarP(&env, "env", "e", []string{}, "Set environment variables")
	runCmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Run in interactive mode")
	runCmd.Flags().StringVar(&llmProvider, "llm", "", "Override the LLM provider specified in the image")
	runCmd.Flags().StringVar(&llmEndpoint, "llm-endpoint", "", "LLM provider endpoint URL")
	runCmd.Flags().StringVar(&llmModel, "llm-model", "", "Override the LLM model specified in the image")
	runCmd.Flags().DurationVar(&timeout, "timeout", 0, "Timeout for the agent run (e.g. 1h, 30m)")

	return runCmd
}

// loadImage loads an image by name
func loadImage(imageName string) (*agent.Image, error) {
	parts := strings.SplitN(imageName, ":", 2)
	name := parts[0]
	tag := parts[1]

	// Construct the file path for the image
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	imagesDir := filepath.Join(homeDir, ".sentinel/images")
	imagePath := filepath.Join(imagesDir, fmt.Sprintf("%s_%s.json", strings.ReplaceAll(name, "/", "_"), tag))

	// Check if the image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image %s not found", imageName)
	}

	// Read the image file
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Parse the image JSON
	var image agent.Image
	if err := json.Unmarshal(data, &image); err != nil {
		return nil, fmt.Errorf("failed to parse image: %w", err)
	}

	return &image, nil
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
