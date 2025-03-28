package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
)

// AgentfileCmd returns the agentfile command
func AgentfileCmd() *cobra.Command {
	var modelEndpoint string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "agentfile",
		Short: "Manage Agentfiles",
		Long:  `Create, edit, and manage Agentfiles for your AI agents.`,
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Agentfile",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Creating new Agentfile '%s'\n", name)
			
			// Create directory for the agent
			agentDir := name
			if err := os.MkdirAll(agentDir, 0755); err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
				os.Exit(1)
			}
			
			// Create natural language file
			nlPath := filepath.Join(agentDir, "agentfile.natural.txt")
			nlContent := "This agent helps users with their tasks. It should be helpful, accurate, and concise.\n" +
				"It should use the Llama3 model and have basic conversation capabilities.\n"
			if err := os.WriteFile(nlPath, []byte(nlContent), 0644); err != nil {
				fmt.Printf("Error creating natural language file: %v\n", err)
				os.Exit(1)
			}
			
			// Create YAML file with default configuration
			defaultAgent := agentfile.DefaultAgentfile(name)
			yamlData, err := yaml.Marshal(defaultAgent)
			if err != nil {
				fmt.Printf("Error creating YAML: %v\n", err)
				os.Exit(1)
			}
			
			yamlPath := filepath.Join(agentDir, "agentfile.yaml")
			if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
				fmt.Printf("Error writing YAML file: %v\n", err)
				os.Exit(1)
			}
			
			// Create empty state file
			statePath := filepath.Join(agentDir, "agent.state.json")
			if err := os.WriteFile(statePath, []byte("{\n}\n"), 0644); err != nil {
				fmt.Printf("Error creating state file: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println("✓ Initialized agentfile.yaml")
			fmt.Println("✓ Created agentfile.natural.txt for natural language definition")
			fmt.Println("✓ Added default state schema")
			fmt.Println("\nDone! Edit agentfile.natural.txt to define your agent's purpose and behavior.")
		},
	}

	convertCmd := &cobra.Command{
		Use:   "convert [file-path]",
		Short: "Convert natural language to YAML",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			
			// Verify file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Printf("Error: File '%s' does not exist\n", filePath)
				os.Exit(1)
			}
			
			// Get the model endpoint from flag or use default
			endpoint, _ := cmd.Flags().GetString("endpoint")
			if endpoint == "" {
				endpoint = "http://model.gonella.co.uk"
			}
			
			// Debug info
			verboseFlag, _ := cmd.Flags().GetBool("verbose")
			if verboseFlag {
				fmt.Printf("DEBUG: Using model endpoint: %s\n", endpoint)
			}
			fmt.Printf("Converting '%s' to YAML using endpoint '%s'...\n", filePath, endpoint)
			
			// Create parser
			parser := agentfile.NewParser(endpoint)
			parser.SetVerbose(verboseFlag)
			
			// Parse file
			yamlPath, err := parser.ParseFile(filePath)
			if err != nil {
				fmt.Printf("Error converting file: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("✓ Successfully converted to '%s'\n", yamlPath)
		},
	}

	createCmd.Flags().StringP("name", "n", "", "Name of the agent")
	createCmd.MarkFlagRequired("name")
	convertCmd.Flags().StringVarP(&modelEndpoint, "endpoint", "e", "http://model.gonella.co.uk", "Ollama API endpoint URL")
	convertCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	
	cmd.AddCommand(createCmd)
	cmd.AddCommand(convertCmd)

	return cmd
}
