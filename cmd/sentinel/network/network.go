package network

import (
	"fmt"

	"github.com/sentinelstacks/sentinel/cmd/sentinel"
	"github.com/spf13/cobra"
)

// NewNetworkCmd creates the network command group
func NewNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage agent networks",
		Long:  `Create and manage networks for agent-to-agent communication`,
	}

	// Add subcommands
	cmd.AddCommand(newNetworkCreateCmd())
	cmd.AddCommand(newNetworkListCmd())
	cmd.AddCommand(newNetworkConnectCmd())
	cmd.AddCommand(newNetworkDisconnectCmd())
	cmd.AddCommand(newNetworkRemoveCmd())
	cmd.AddCommand(newNetworkInspectCmd())

	return cmd
}

// newNetworkCreateCmd creates the network create command
func newNetworkCreateCmd() *cobra.Command {
	var driver string

	cmd := &cobra.Command{
		Use:   "create [network_name]",
		Short: "Create a new agent network",
		Long:  `Create a new network for agents to communicate with each other`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			networkName := args[0]
			
			fmt.Printf("Creating network '%s' with driver '%s'\n", networkName, driver)
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// Create the network
			network, err := sp.CreateNetwork(ctx, networkName, driver)
			if err != nil {
				return fmt.Errorf("failed to create network: %w", err)
			}
			
			fmt.Printf("Network '%s' created successfully\n", networkName)
			return nil
		},
	}

	cmd.Flags().StringVar(&driver, "driver", "default", "Network driver to use")
	return cmd
}

// newNetworkListCmd creates the network list command
func newNetworkListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List networks",
		Long:    `List all networks available for agent communication`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			
			fmt.Println("Listing all networks:")
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// List networks
			networks, err := sp.ListNetworks(ctx)
			if err != nil {
				return fmt.Errorf("failed to list networks: %w", err)
			}
			
			if len(networks) == 0 {
				fmt.Println("No networks found")
				return nil
			}
			
			for _, network := range networks {
				if netMap, ok := network.(map[string]interface{}); ok {
					fmt.Printf("%v [%v]\n", 
						netMap["name"], 
						netMap["status"])
				}
			}
			
			return nil
		},
	}
}

// newNetworkConnectCmd creates the network connect command
func newNetworkConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect [network_name] [agent_id]",
		Short: "Connect an agent to a network",
		Long:  `Connect an existing agent to a specified network`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			networkName := args[0]
			agentID := args[1]
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// Connect agent to network
			if err := sp.ConnectAgentToNetwork(ctx, networkName, agentID); err != nil {
				return fmt.Errorf("failed to connect agent to network: %w", err)
			}
			
			fmt.Printf("Agent '%s' successfully connected to network '%s'\n", agentID, networkName)
			return nil
		},
	}
}

// newNetworkDisconnectCmd creates the network disconnect command
func newNetworkDisconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect [network_name] [agent_id]",
		Short: "Disconnect an agent from a network",
		Long:  `Disconnect an agent from a specified network`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			networkName := args[0]
			agentID := args[1]
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// Disconnect agent from network
			if err := sp.DisconnectAgentFromNetwork(ctx, networkName, agentID); err != nil {
				return fmt.Errorf("failed to disconnect agent from network: %w", err)
			}
			
			fmt.Printf("Agent '%s' successfully disconnected from network '%s'\n", agentID, networkName)
			return nil
		},
	}
}

// newNetworkRemoveCmd creates the network remove command
func newNetworkRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "rm [network_name]",
		Aliases: []string{"remove"},
		Short:   "Remove a network",
		Long:    `Remove a specified network`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			networkName := args[0]
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// Remove the network
			if err := sp.RemoveNetwork(ctx, networkName, force); err != nil {
				return fmt.Errorf("failed to remove network: %w", err)
			}
			
			fmt.Printf("Network '%s' successfully removed\n", networkName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force removal even if network has connected agents")
	return cmd
}

// newNetworkInspectCmd creates the network inspect command
func newNetworkInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [network_name]",
		Short: "Display detailed information on a network",
		Long:  `Display detailed information about a network, including connected agents`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			networkName := args[0]
			
			fmt.Printf("Inspecting network '%s'\n\n", networkName)
			
			// Get the service provider from context
			sp := sentinel.GetServiceProvider(ctx)
			
			// Get network details
			network, err := sp.InspectNetwork(ctx, networkName)
			if err != nil {
				return fmt.Errorf("failed to inspect network: %w", err)
			}
			
			// Display network details
			if netMap, ok := network.(map[string]interface{}); ok {
				fmt.Printf("Network: %v\n", netMap["name"])
				
				if v, ok := netMap["created"]; ok {
					fmt.Printf("  Created: %v\n", v)
				}
				
				if v, ok := netMap["status"]; ok {
					fmt.Printf("  Status: %v\n", v)
				}
				
				if v, ok := netMap["driver"]; ok {
					fmt.Printf("  Driver: %v\n", v)
				}
				
				if agents, ok := netMap["agents"].([]string); ok {
					fmt.Printf("  Connected Agents: %d\n", len(agents))
					for _, agent := range agents {
						fmt.Printf("    - %s\n", agent)
					}
				}
			}
			
			return nil
		},
	}
}
