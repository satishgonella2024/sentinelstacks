package tools

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sentinelstacks/sentinel/internal/runtime"
	"github.com/sentinelstacks/sentinel/internal/tools"
)

// NewToolsCmd creates the tools command
func NewToolsCmd() *cobra.Command {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage agent tools",
		Long:  `Manage tools and permissions for Sentinel agents`,
	}
	
	// Add subcommands
	toolsCmd.AddCommand(NewToolsListCmd())
	toolsCmd.AddCommand(NewToolsPermsCmd())
	toolsCmd.AddCommand(NewToolsGrantCmd())
	toolsCmd.AddCommand(NewToolsRevokeCmd())
	
	return toolsCmd
}

// NewToolsListCmd creates the tools list command
func NewToolsListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available tools",
		Long:  `List all available tools that can be used by Sentinel agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get registry
			registry := tools.GetRegistry()
			
			// Get all tools
			allTools := registry.ListTools()
			
			if len(allTools) == 0 {
				fmt.Println("No tools registered.")
				return nil
			}
			
			// Group tools by permission
			toolsByPerm := make(map[tools.Permission][]tools.Tool)
			for _, tool := range allTools {
				perm := tool.RequiredPermission()
				toolsByPerm[perm] = append(toolsByPerm[perm], tool)
			}
			
			// Print tools grouped by permission
			fmt.Println("Available tools:")
			fmt.Println()
			
			// Order of permissions to display
			permOrder := []tools.Permission{
				tools.PermissionNone,
				tools.PermissionFile,
				tools.PermissionNetwork,
				tools.PermissionAPI,
				tools.PermissionShell,
			}
			
			for _, perm := range permOrder {
				tools, ok := toolsByPerm[perm]
				if !ok {
					continue
				}
				
				fmt.Printf("Permission: %s\n", perm)
				fmt.Println(strings.Repeat("-", 40))
				
				for _, tool := range tools {
					fmt.Printf("  %s:\n", tool.GetName())
					fmt.Printf("    Description: %s\n", tool.GetDescription())
					
					// Print parameters
					if params := tool.GetParameters(); len(params) > 0 {
						fmt.Println("    Parameters:")
						for _, param := range params {
							reqStr := ""
							if param.Required {
								reqStr = " (required)"
							}
							fmt.Printf("      %s: %s%s\n", param.Name, param.Description, reqStr)
						}
					}
					
					fmt.Println()
				}
			}
			
			return nil
		},
	}
	
	return listCmd
}

// NewToolsPermsCmd creates the tools permissions command
func NewToolsPermsCmd() *cobra.Command {
	var all bool
	
	permsCmd := &cobra.Command{
		Use:   "perms [agent-id]",
		Short: "Show agent tool permissions",
		Long:  `Show tool permissions for a specific agent or all agents`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create permission manager
			permManager, err := tools.NewPermissionManager("")
			if err != nil {
				return fmt.Errorf("failed to create permission manager: %w", err)
			}
			
			// Show permissions for all agents
			if all || len(args) == 0 {
				// Get all agent IDs
				agentIDs := permManager.GetAgentIDs()
				
				if len(agentIDs) == 0 {
					fmt.Println("No agents have tool permissions.")
					return nil
				}
				
				fmt.Println("Agent tool permissions:")
				fmt.Println()
				
				for _, agentID := range agentIDs {
					perms := permManager.GetPermissions(agentID)
					
					// Get agent info
					rt, err := runtime.GetRuntime()
					if err == nil {
						agentInfo, err := rt.GetAgent(agentID)
						if err == nil {
							fmt.Printf("Agent: %s (%s)\n", agentInfo.Name, agentID)
						} else {
							fmt.Printf("Agent: %s\n", agentID)
						}
					} else {
						fmt.Printf("Agent: %s\n", agentID)
					}
					
					if len(perms) == 0 {
						fmt.Println("  No permissions")
					} else {
						for _, perm := range perms {
							fmt.Printf("  - %s\n", perm)
						}
					}
					
					fmt.Println()
				}
				
				return nil
			}
			
			// Show permissions for specific agent
			agentID := args[0]
			
			// Get agent info
			rt, err := runtime.GetRuntime()
			if err == nil {
				agentInfo, err := rt.GetAgent(agentID)
				if err == nil {
					fmt.Printf("Agent: %s (%s)\n", agentInfo.Name, agentID)
				} else {
					fmt.Printf("Agent: %s\n", agentID)
				}
			} else {
				fmt.Printf("Agent: %s\n", agentID)
			}
			
			// Get permissions
			perms := permManager.GetPermissions(agentID)
			
			if len(perms) == 0 {
				fmt.Println("  No permissions")
			} else {
				for _, perm := range perms {
					fmt.Printf("  - %s\n", perm)
				}
			}
			
			return nil
		},
	}
	
	// Add flags
	permsCmd.Flags().BoolVarP(&all, "all", "a", false, "Show permissions for all agents")
	
	return permsCmd
}

// NewToolsGrantCmd creates the tools grant command
func NewToolsGrantCmd() *cobra.Command {
	grantCmd := &cobra.Command{
		Use:   "grant [agent-id] [permission]",
		Short: "Grant tool permission to an agent",
		Long:  `Grant a tool permission to a specific agent`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			permStr := args[1]
			
			// Validate permission
			perm := tools.Permission(permStr)
			validPerms := map[tools.Permission]bool{
				tools.PermissionNone:    true,
				tools.PermissionFile:    true,
				tools.PermissionNetwork: true,
				tools.PermissionAPI:     true,
				tools.PermissionShell:   true,
				tools.PermissionAll:     true,
			}
			
			if !validPerms[perm] {
				return fmt.Errorf("invalid permission: %s", permStr)
			}
			
			// Create permission manager
			permManager, err := tools.NewPermissionManager("")
			if err != nil {
				return fmt.Errorf("failed to create permission manager: %w", err)
			}
			
			// Grant permission
			if err := permManager.Grant(agentID, perm); err != nil {
				return fmt.Errorf("failed to grant permission: %w", err)
			}
			
			fmt.Printf("Granted %s permission to agent %s\n", perm, agentID)
			
			return nil
		},
	}
	
	return grantCmd
}

// NewToolsRevokeCmd creates the tools revoke command
func NewToolsRevokeCmd() *cobra.Command {
	revokeCmd := &cobra.Command{
		Use:   "revoke [agent-id] [permission]",
		Short: "Revoke tool permission from an agent",
		Long:  `Revoke a tool permission from a specific agent`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			permStr := args[1]
			
			// Validate permission
			perm := tools.Permission(permStr)
			validPerms := map[tools.Permission]bool{
				tools.PermissionNone:    true,
				tools.PermissionFile:    true,
				tools.PermissionNetwork: true,
				tools.PermissionAPI:     true,
				tools.PermissionShell:   true,
				tools.PermissionAll:     true,
			}
			
			if !validPerms[perm] {
				return fmt.Errorf("invalid permission: %s", permStr)
			}
			
			// Create permission manager
			permManager, err := tools.NewPermissionManager("")
			if err != nil {
				return fmt.Errorf("failed to create permission manager: %w", err)
			}
			
			// Revoke permission
			if err := permManager.Revoke(agentID, perm); err != nil {
				return fmt.Errorf("failed to revoke permission: %w", err)
			}
			
			fmt.Printf("Revoked %s permission from agent %s\n", perm, agentID)
			
			return nil
		},
	}
	
	return revokeCmd
}