package commands

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new root command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sentinel",
		Short: "SentinelStacks - Docker for AI Agents",
		Long: `SentinelStacks is a platform for creating, running, sharing, and orchestrating AI agents.
It provides a natural language way to define agents that can leverage different model backends.`,
	}

	// Add commands
	cmd.AddCommand(AgentfileCmd())
	cmd.AddCommand(AgentCmd())
	cmd.AddCommand(RegistryCmd())
	cmd.AddCommand(VersionCmd())

	return cmd
}
