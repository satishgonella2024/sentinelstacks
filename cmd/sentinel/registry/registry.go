// Package registry contains commands for interacting with the registry
package registry

import (
	"github.com/spf13/cobra"
)

// NewRegistryCmd creates a new registry command
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Manage agent registries",
		Long:  `Commands for managing agent registries, including push, pull, login, and logout.`,
	}
	// Add subcommands
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewPushCmd())
	cmd.AddCommand(NewPullCmd())
	cmd.AddCommand(NewSearchCmd())
	cmd.AddCommand(NewTagsCmd())
	return cmd
}
