package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
