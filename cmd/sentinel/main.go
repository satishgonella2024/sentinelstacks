package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/satishgonella2024/sentinelstacks/pkg/agentfile"
)

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "SentinelStacks - Docker for AI Agents",
	Long: `SentinelStacks is a platform for creating, running, sharing, and orchestrating AI agents.
It provides a natural language way to define agents that can leverage different model backends.`,
}

func init() {
	// Add commands
	rootCmd.AddCommand(agentfileCmd())
	rootCmd.AddCommand(agentCmd())
	rootCmd.AddCommand(registryCmd())
}

func agentfileCmd() *cobra.Command {
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

	createCmd.Flags().StringP("name", "n", "", "Name of the agent")
	createCmd.MarkFlagRequired("name")
	cmd.AddCommand(createCmd)

	return cmd
}

func agentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage and run agents",
		Long:  `Run, monitor, and manage your AI agents.`,
	}

	runCmd := &cobra.Command{
		Use:   "run [agent-name]",
		Short: "Run an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Running agent: %s\n", args[0])
			// Agent runtime code will go here
		},
	}

	cmd.AddCommand(runCmd)
	return cmd
}

func registryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Interact with the agent registry",
		Long:  `Push, pull, and search for agents in the registry.`,
	}

	pushCmd := &cobra.Command{
		Use:   "push [agent-name]",
		Short: "Push an agent to the registry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Pushing agent to registry: %s\n", args[0])
			// Registry code will go here
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull [agent-name]",
		Short: "Pull an agent from the registry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Pulling agent from registry: %s\n", args[0])
			// Registry code will go here
		},
	}

	cmd.AddCommand(pushCmd)
	cmd.AddCommand(pullCmd)
	return cmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
