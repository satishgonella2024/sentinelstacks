// Package chat implements the chat command for the Sentinel CLI
package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/sentinelstacks/sentinel/internal/multimodal"
	"github.com/sentinelstacks/sentinel/internal/runtime"
	"github.com/spf13/cobra"
)

const defaultModel = "mock-model"

// NewChatCmd creates a new chat command
func NewChatCmd() *cobra.Command {
	var provider string
	var model string
	var apiKey string
	var endpoint string
	var imagePaths []string
	var temperature float64
	var maxTokens int
	var name string

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive chat session with an agent",
		Long: `Start an interactive chat session with an agent.
This command starts a real-time chat session with an agent, allowing you to
interact with the agent and receive responses in real-time.

You can optionally specify a provider and model to use for the session. If not
specified, the default provider and model will be used.

Examples:
  # Start a chat session with the default provider and model
  sentinel chat

  # Start a chat session with a specific provider and model
  sentinel chat --provider claude --model claude-3-opus-20240229

  # Start a chat session with images
  sentinel chat --images image1.jpg,image2.png`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get the API key from environment variable if not specified
			if apiKey == "" {
				switch provider {
				case "claude":
					apiKey = os.Getenv("ANTHROPIC_API_KEY")
				case "openai":
					apiKey = os.Getenv("OPENAI_API_KEY")
				}
			}

			// Create a runtime
			rt, err := runtime.NewRuntime("")
			if err != nil {
				fmt.Printf("Error creating runtime: %v\n", err)
				os.Exit(1)
			}

			// Create a temporary agent name if not specified
			if name == "" {
				name = fmt.Sprintf("chat-%d", os.Getpid())
			}

			// Prepare configuration for the agent
			agentImage := fmt.Sprintf("%s:%s", provider, model)

			// Create the multimodal agent
			fmt.Println("Creating multimodal agent...")
			agent, err := rt.CreateMultimodalAgent(name, agentImage, model, provider, apiKey, endpoint)
			if err != nil {
				fmt.Printf("Error creating agent: %v\n", err)
				os.Exit(1)
			}

			// Set agent parameters
			if temperature > 0 {
				agent.SetTemperature(temperature)
			}
			if maxTokens > 0 {
				agent.SetMaxTokens(maxTokens)
			}

			// Add system prompt
			agent.AddSystemPrompt("You are a helpful, harmless, and honest AI assistant. Respond concisely to user queries.")

			userColor := color.New(color.FgBlue).SprintFunc()
			assistantColor := color.New(color.FgGreen).SprintFunc()

			// Print welcome message
			fmt.Println("\n" + assistantColor("Assistant") + ": Hello! I'm an AI assistant powered by " + provider + ". How can I help you today?")

			// Check if initial images were provided
			if len(imagePaths) > 0 {
				// Process initial images
				contents, err := processImages(imagePaths)
				if err != nil {
					fmt.Printf("Error processing images: %v\n", err)
					os.Exit(1)
				}

				contents = append([]*multimodal.Content{multimodal.NewTextContent("What's in these images?")}, contents...)

				fmt.Println(userColor("You") + ": [Uploaded images with question: What's in these images?]")

				// Get response from agent
				output, err := agent.ProcessMultimodalInput(cmd.Context(), contents)
				if err != nil {
					fmt.Printf("Error from agent: %v\n", err)
					return
				}

				// Extract text response
				var responseText string
				for _, content := range output.Contents {
					if content.Type == multimodal.MediaTypeText {
						responseText = content.Text
						break
					}
				}

				// Print response
				fmt.Println(assistantColor("Assistant") + ": " + responseText)
			}

			// Main chat loop
			scanner := NewInputScanner()
			for {
				// Read user input
				fmt.Print(userColor("You") + ": ")
				text, err := scanner.ReadInput()
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					break
				}

				// Check for exit command
				if strings.ToLower(text) == "exit" || strings.ToLower(text) == "quit" {
					break
				}

				// Check for image upload command
				if strings.HasPrefix(text, "/image") || strings.HasPrefix(text, "/img") {
					// Extract image paths from command
					parts := strings.SplitN(text, " ", 2)
					if len(parts) < 2 {
						fmt.Println("Please specify image path(s). Example: /image path/to/image.jpg")
						continue
					}

					// Parse image paths
					imgPaths := strings.Split(parts[1], ",")
					for i, path := range imgPaths {
						imgPaths[i] = strings.TrimSpace(path)
					}

					// Ask for a description
					fmt.Print("Enter a question about the image(s): ")
					question, err := scanner.ReadInput()
					if err != nil {
						fmt.Printf("Error reading input: %v\n", err)
						continue
					}

					// Process images
					contents, err := processImages(imgPaths)
					if err != nil {
						fmt.Printf("Error processing images: %v\n", err)
						continue
					}

					// Add question to contents
					contents = append([]*multimodal.Content{multimodal.NewTextContent(question)}, contents...)

					// Get response from agent
					output, err := agent.ProcessMultimodalInput(cmd.Context(), contents)
					if err != nil {
						fmt.Printf("Error from agent: %v\n", err)
						continue
					}

					// Extract text response
					var responseText string
					for _, content := range output.Contents {
						if content.Type == multimodal.MediaTypeText {
							responseText = content.Text
							break
						}
					}

					// Print response
					fmt.Println(assistantColor("Assistant") + ": " + responseText)
					continue
				}

				// Process regular text input
				response, err := agent.ProcessTextInput(cmd.Context(), text)
				if err != nil {
					fmt.Printf("Error from agent: %v\n", err)
					continue
				}

				// Print response
				fmt.Println(assistantColor("Assistant") + ": " + response)
			}

			// Close the agent
			if err := agent.Close(); err != nil {
				fmt.Printf("Error closing agent: %v\n", err)
			}

			fmt.Println("Chat session ended.")
		},
	}

	// Add flags
	chatCmd.Flags().StringVar(&provider, "provider", "mock", "The LLM provider to use (claude, openai, ollama)")
	chatCmd.Flags().StringVar(&model, "model", defaultModel, "The model to use")
	chatCmd.Flags().StringVar(&apiKey, "api-key", "", "The API key to use (if not specified, will use environment variable)")
	chatCmd.Flags().StringVar(&endpoint, "endpoint", "", "The API endpoint to use (if not default)")
	chatCmd.Flags().StringSliceVar(&imagePaths, "images", []string{}, "Comma-separated list of image paths to include in the initial prompt")
	chatCmd.Flags().Float64Var(&temperature, "temperature", 0.7, "The temperature to use for generation")
	chatCmd.Flags().IntVar(&maxTokens, "max-tokens", 4096, "The maximum number of tokens to generate")
	chatCmd.Flags().StringVar(&name, "name", "", "The name to use for the agent")

	return chatCmd
}

// InputScanner handles user input with support for multiline input
type InputScanner struct {
	buffer []string
}

// NewInputScanner creates a new input scanner
func NewInputScanner() *InputScanner {
	return &InputScanner{
		buffer: make([]string, 0),
	}
}

// ReadInput reads user input, supporting multiline input with a '\'
func (s *InputScanner) ReadInput() (string, error) {
	var line string
	var err error

	// Read initial line
	line, err = readLine()
	if err != nil {
		return "", err
	}

	// Check if it ends with escape character
	for strings.HasSuffix(line, "\\") {
		// Remove the escape character
		line = line[:len(line)-1]

		// Add to buffer
		s.buffer = append(s.buffer, line)

		// Read next line
		fmt.Print("... ")
		line, err = readLine()
		if err != nil {
			return "", err
		}
	}

	// Add final line to buffer
	s.buffer = append(s.buffer, line)

	// Combine all lines
	result := strings.Join(s.buffer, "\n")

	// Clear buffer
	s.buffer = s.buffer[:0]

	return result, nil
}

// readLine reads a line from standard input
func readLine() (string, error) {
	var line string
	_, err := fmt.Scanln(&line)

	// Handle empty lines
	if err != nil && err.Error() == "unexpected newline" {
		return "", nil
	}

	return line, err
}

// processImages reads and processes image files
func processImages(paths []string) ([]*multimodal.Content, error) {
	contents := make([]*multimodal.Content, 0, len(paths))

	for _, path := range paths {
		// Expand the path
		expandedPath := expandPath(path)

		// Read the image file
		data, err := os.ReadFile(expandedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file %s: %w", expandedPath, err)
		}

		// Detect mime type based on extension
		mimeType := detectMimeType(expandedPath)

		// Add to contents
		contents = append(contents, multimodal.NewImageContent(data, mimeType))
	}

	return contents, nil
}

// expandPath expands a path with tilde to absolute path
func expandPath(path string) string {
	if path == "" {
		return path
	}

	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}

	return path
}

// detectMimeType detects the mime type of a file based on its extension
func detectMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	default:
		return "application/octet-stream"
	}
}
