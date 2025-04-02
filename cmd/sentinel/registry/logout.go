package registry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/registry"
)

// NewLogoutCmd creates a new logout command
func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout [registry_url]",
		Short: "Log out from an agent registry",
		Long:  `Log out from an agent registry, removing stored credentials.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var registryURL string
			if len(args) > 0 {
				registryURL = args[0]
			} else {
				registryURL = registry.DefaultRegistryURL
			}
			
			return runLogout(registryURL)
		},
	}

	return cmd
}

func runLogout(registryURL string) error {
	// Check if logged in
	if !registry.IsLoggedIn(registryURL) {
		return fmt.Errorf("not logged in to %s", registryURL)
	}

	// Log out
	if err := registry.Logout(registryURL); err != nil {
		return fmt.Errorf("failed to log out from registry: %w", err)
	}

	fmt.Printf("Successfully logged out from %s\n", registryURL)
	return nil
}
