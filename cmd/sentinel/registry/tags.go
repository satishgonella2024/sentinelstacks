package registry

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
)

// NewTagsCmd creates a new tags command
func NewTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags [package-name]",
		Short: "List available tags for a package",
		Long:  `List all available tags (versions) for a package in the registry.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get package name
			packageName := args[0]

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

			// Create registry client
			clientConfig := client.Config{
				BaseURL:      registryURL,
				AuthProvider: authProvider,
			}
			registryClient := client.NewClient(clientConfig)

			// Get package tags
			fmt.Printf("Getting tags for package %s from registry %s...\n", packageName, registryURL)
			tags, err := registryClient.ListTags(context.Background(), packageName)
			if err != nil {
				return fmt.Errorf("failed to get tags: %w", err)
			}

			// Display tags
			if len(tags) == 0 {
				fmt.Println("No tags found for package:", packageName)
				return nil
			}

			// Sort tags
			sort.Strings(tags)

			// Print tags
			fmt.Printf("Tags for package %s:\n", packageName)
			for _, tag := range tags {
				fmt.Printf("  %s\n", tag)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")

	return cmd
}
