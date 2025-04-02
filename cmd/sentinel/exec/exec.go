package exec

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewExecCmd creates the exec command
func NewExecCmd() *cobra.Command {
	var (
		llmProvider string
		llmEndpoint string
		llmModel    string
		timeout     time.Duration
		maxTokens   int
		temperature float64
		pipe        bool
	)

	cmd := &cobra.Command{
		Use:   "exec [prompt]",
		Short: "Execute a one-time prompt",
		Long:  `Execute a one-time prompt without creating an agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var prompt string
			
			// Get prompt from arguments or stdin
			if pipe {
				// Read from stdin
				stdinBytes, err := os.ReadFile(os.Stdin.Name())
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
				prompt = string(stdinBytes)
			} else if len(args) > 0 {
				// Get prompt from args
				prompt = strings.Join(args, " ")
			} else {
				return fmt.Errorf("prompt required: provide as argument or pipe with --pipe")
			}
			
			// Get provider configuration
			if llmProvider == "" {
				llmProvider = viper.GetString("llm.provider")
				if llmProvider == "" {
					llmProvider = "claude" // Default
				}
			}
			
			if llmModel == "" {
				llmModel = viper.GetString("llm.model")
				// Set appropriate default based on provider
				if llmModel == "" {
					if llmProvider == "ollama" {
						llmModel = "llama3"
					} else if llmProvider == "claude" {
						llmModel = "claude-3-5-sonnet-20240627"
					} else if llmProvider == "openai" {
						llmModel = "gpt-4"
					} else if llmProvider == "google" {
						llmModel = "gemini-1.5-pro"
					}
				}
			}
			
			// Print configuration
			fmt.Printf("Executing prompt with %s (%s):\n", llmProvider, llmModel)
			fmt.Printf("Prompt: %s\n\n", truncateString(prompt, 100))
			
			// Set up context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			
			// Handle termination signals
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			
			go func() {
				<-signalCh
				fmt.Println("\nReceived termination signal, cancelling...")
				cancel()
			}()
			
			// TODO: Replace this with actual LLM call
			fmt.Println("Response:")
			fmt.Println("This is a simulated response to your prompt. In a real implementation,")
			fmt.Println("this would connect to the appropriate LLM provider and return the")
			fmt.Println("response from the model.")
			
			return nil
		},
	}

	cmd.Flags().StringVar(&llmProvider, "llm", "", "LLM provider (claude, openai, ollama, google)")
	cmd.Flags().StringVar(&llmEndpoint, "llm-endpoint", "", "LLM provider endpoint URL")
	cmd.Flags().StringVar(&llmModel, "llm-model", "", "LLM model to use")
	cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Timeout for the execution")
	cmd.Flags().IntVar(&maxTokens, "max-tokens", 2048, "Maximum tokens in the response")
	cmd.Flags().Float64Var(&temperature, "temperature", 0.7, "Temperature for the LLM call")
	cmd.Flags().BoolVar(&pipe, "pipe", false, "Read prompt from stdin")

	return cmd
}

// Helper functions

// truncateString truncates a string to the specified length and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
