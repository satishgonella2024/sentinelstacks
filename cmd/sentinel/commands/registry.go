package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/pkg/registry"
)

// RegistryCmd returns the registry command
func RegistryCmd() *cobra.Command {
	var tags []string
	var visibility string

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
			agentName := args[0]
			
			// Create registry client
			client, err := registry.NewRegistryClient()
			if err != nil {
				fmt.Printf("Error creating registry client: %v\n", err)
				os.Exit(1)
			}
			
			// Get visibility flag
			visibilityFlag, _ := cmd.Flags().GetString("visibility")
			if visibilityFlag == "" {
				visibilityFlag = "public"
			}
			
			// Push agent to registry
			fmt.Printf("Pushing agent '%s' to registry with visibility '%s'...\n", agentName, visibilityFlag)
			err = client.PushAgent(agentName, visibilityFlag)
			if err != nil {
				fmt.Printf("Error pushing agent to registry: %v\n", err)
				os.Exit(1)
			}
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull [username/agent-name]",
		Short: "Pull an agent from the registry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			agentRef := args[0]
			
			// Create registry client
			client, err := registry.NewRegistryClient()
			if err != nil {
				fmt.Printf("Error creating registry client: %v\n", err)
				os.Exit(1)
			}
			
			// Pull agent from registry
			fmt.Printf("Pulling agent '%s' from registry...\n", agentRef)
			path, err := client.PullAgent(agentRef)
			if err != nil {
				fmt.Printf("Error pulling agent from registry: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Successfully pulled agent to '%s'\n", path)
		},
	}

	searchCmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for agents in the registry",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			
			tagsFlag, _ := cmd.Flags().GetStringSlice("tags")
			
			// Create registry client
			client, err := registry.NewRegistryClient()
			if err != nil {
				fmt.Printf("Error creating registry client: %v\n", err)
				os.Exit(1)
			}
			
			// Search registry
			fmt.Printf("Searching registry for '%s'", query)
			if len(tagsFlag) > 0 {
				fmt.Printf(" with tags: %s", strings.Join(tagsFlag, ", "))
			}
			fmt.Println("...")
			
			results, err := client.SearchAgents(query, tagsFlag)
			if err != nil {
				fmt.Printf("Error searching registry: %v\n", err)
				os.Exit(1)
			}
			
			if len(results) == 0 {
				fmt.Println("No agents found.")
				return
			}
			
			fmt.Printf("Found %d agents:\n\n", len(results))
			for i, agent := range results {
				fmt.Printf("%d. %s/%s@%s\n", i+1, agent.Author, agent.Name, agent.Version)
				fmt.Printf("   Description: %s\n", agent.Description)
				if len(agent.Models) > 0 {
					fmt.Printf("   Models: %s\n", strings.Join(agent.Models, ", "))
				}
				if len(agent.Capabilities) > 0 {
					fmt.Printf("   Capabilities: %s\n", strings.Join(agent.Capabilities, ", "))
				}
				if len(agent.Tags) > 0 {
					fmt.Printf("   Tags: %s\n", strings.Join(agent.Tags, ", "))
				}
				fmt.Printf("   Downloads: %d\n", agent.Downloads)
				fmt.Println()
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all agents in the registry",
		Run: func(cmd *cobra.Command, args []string) {
			// Create registry client
			client, err := registry.NewRegistryClient()
			if err != nil {
				fmt.Printf("Error creating registry client: %v\n", err)
				os.Exit(1)
			}
			
			// List agents
			fmt.Println("Listing agents in registry...")
			
			agents, err := client.ListAgents()
			if err != nil {
				fmt.Printf("Error listing agents: %v\n", err)
				os.Exit(1)
			}
			
			if len(agents) == 0 {
				fmt.Println("No agents found in the registry.")
				return
			}
			
			fmt.Printf("Found %d agents:\n\n", len(agents))
			for i, agent := range agents {
				fmt.Printf("%d. %s/%s@%s\n", i+1, agent.Author, agent.Name, agent.Version)
				fmt.Printf("   Description: %s\n", agent.Description)
				if len(agent.Models) > 0 {
					fmt.Printf("   Models: %s\n", strings.Join(agent.Models, ", "))
				}
				if len(agent.Capabilities) > 0 {
					fmt.Printf("   Capabilities: %s\n", strings.Join(agent.Capabilities, ", "))
				}
				fmt.Println()
			}
		},
	}

	pushCmd.Flags().StringVarP(&visibility, "visibility", "v", "public", "Visibility of the agent (public or private)")
	searchCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Filter by tags (comma-separated)")
	
	cmd.AddCommand(pushCmd)
	cmd.AddCommand(pullCmd)
	cmd.AddCommand(searchCmd)
	cmd.AddCommand(listCmd)

	return cmd
}
