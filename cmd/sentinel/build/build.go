package build

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinelstacks/sentinel/internal/parser"
	"github.com/sentinelstacks/sentinel/internal/registry"
	"github.com/sentinelstacks/sentinel/pkg/agent"
)

// NewBuildCmd creates the build command
func NewBuildCmd() *cobra.Command {
	var (
		tag         string
		file        string
		noCache     bool
		llmProvider string
		llmEndpoint string
		llmModel    string
	)

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build a Sentinel Image from a Sentinelfile",
		Long:  `Build a Sentinel Image from a Sentinelfile, creating a packageable agent definition`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate tag
			if tag == "" {
				return fmt.Errorf("tag is required (use -t name:tag)")
			}

			if !strings.Contains(tag, ":") {
				tag = tag + ":latest"
			}

			// Check that Sentinelfile exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return fmt.Errorf("Sentinelfile not found at %s", file)
			}

			// If LLM provider was not specified, use the one from the config
			if llmProvider == "" {
				llmProvider = viper.GetString("llm.provider")
				if llmProvider == "" {
					llmProvider = "ollama" // Default
				}
			}

			// If LLM endpoint was not specified, use the one from the config
			if llmEndpoint == "" {
				llmEndpoint = viper.GetString("llm.endpoint")

				// Set appropriate default based on provider
				if llmEndpoint == "" {
					if llmProvider == "ollama" {
						llmEndpoint = "http://model.gonella.co.uk/api/generate"
					}
				}
			}

			// If LLM model was not specified, use the one from the config
			if llmModel == "" {
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

			// Set environment variables for the parser
			os.Setenv("SENTINEL_LLM_PROVIDER", llmProvider)
			os.Setenv("SENTINEL_LLM_ENDPOINT", llmEndpoint)
			os.Setenv("SENTINEL_LLM_MODEL", llmModel)

			// Get the API key from the config
			apiKey := viper.GetString("llm.api_key")
			if apiKey != "" {
				os.Setenv("SENTINEL_API_KEY", apiKey)
			}

			fmt.Printf("Building image %s from Sentinelfile\n", tag)
			fmt.Printf("Using LLM provider: %s\n", llmProvider)
			if llmEndpoint != "" {
				fmt.Printf("LLM endpoint: %s\n", llmEndpoint)
			}
			if llmModel != "" {
				fmt.Printf("LLM model: %s\n", llmModel)
			}
			fmt.Printf("Cache enabled: %v\n", !noCache)

			// Parse the Sentinelfile
			fmt.Println("Parsing Sentinelfile...")
			p := parser.NewSentinelfileParser(llmProvider)
			def, err := p.ParseFile(file)
			if err != nil {
				return fmt.Errorf("failed to parse Sentinelfile: %w", err)
			}

			// Validate the definition
			if err := parser.ValidateDefinition(def); err != nil {
				return fmt.Errorf("invalid agent definition: %w", err)
			}

			fmt.Println("Extracting agent capabilities...")
			if len(def.Capabilities) > 0 {
				fmt.Println("  Capabilities:")
				for _, cap := range def.Capabilities {
					fmt.Printf("   - %s\n", cap)
				}
			}

			fmt.Println("Configuring tools...")
			if len(def.Tools) > 0 {
				fmt.Println("  Tools:")
				for _, tool := range def.Tools {
					fmt.Printf("   - %s\n", tool)
				}
			}

			// Extract name and tagVersion from the tag
			tagParts := strings.SplitN(tag, ":", 2)
			imageName := tagParts[0]
			tagVersion := tagParts[1]

			// Create an image
			image := &agent.Image{
				ID:         generateImageID(),
				Name:       imageName,
				Tag:        tagVersion,
				CreatedAt:  time.Now().Unix(),
				Definition: *def,
			}

			// Get the registry
			reg, err := registry.GetLocalRegistry()
			if err != nil {
				return fmt.Errorf("failed to get registry: %w", err)
			}

			// Save the image to the registry
			regImage := registry.ConvertFromAgentImage(image)
			if err := reg.Save(regImage); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}

			fmt.Println("Creating Sentinel Image...")
			fmt.Printf("Successfully built image: %s\n", tag)
			fmt.Printf("Image ID: %s\n", image.ID)

			return nil
		},
	}

	buildCmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the 'name:tag' format")
	buildCmd.Flags().StringVarP(&file, "file", "f", "Sentinelfile", "Path to Sentinelfile")
	buildCmd.Flags().BoolVar(&noCache, "no-cache", false, "Do not use cache when building the image")
	buildCmd.Flags().StringVar(&llmProvider, "llm", "", "LLM provider to use for parsing (claude, ollama, etc.)")
	buildCmd.Flags().StringVar(&llmEndpoint, "llm-endpoint", "", "LLM provider endpoint URL")
	buildCmd.Flags().StringVar(&llmModel, "llm-model", "", "LLM model to use for parsing")

	buildCmd.MarkFlagRequired("tag")

	return buildCmd
}

// generateImageID generates a unique ID for an image
func generateImageID() string {
	return fmt.Sprintf("img_%x", time.Now().UnixNano())
}
