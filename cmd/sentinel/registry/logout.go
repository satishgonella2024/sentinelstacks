package registry

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
)

// NewLogoutCmd creates a new logout command
func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out from a registry",
		Long:  `Log out from a registry, invalidating the current authentication token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get registry URL from flag or env var
			registryURL, _ := cmd.Flags().GetString("registry")
			if registryURL == "" {
				registryURL = os.Getenv("SENTINEL_REGISTRY_URL")
			}
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Create auth provider
			authConfig := auth.AuthConfig{
				RegistryURL: registryURL,
			}
			authProvider, err := auth.NewFileTokenProvider(authConfig)
			if err != nil {
				return fmt.Errorf("failed to create auth provider: %w", err)
			}

			// Check if already logged in
			ctx := context.Background()
			if !authProvider.IsAuthenticated(ctx) {
				fmt.Printf("Not logged in to registry %s\n", registryURL)
				return nil
			}

			// Logout
			if err := authProvider.Logout(ctx); err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}

			fmt.Printf("Successfully logged out from registry %s\n", registryURL)
			return nil
		},
	}

	// Add flags
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")

	return cmd
}
