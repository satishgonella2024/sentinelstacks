package registry

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
)

// NewPullCmd creates a new pull command
func NewPullCmd() *cobra.Command {
	var (
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "pull [package-name]:[version]",
		Short: "Pull a package from the registry",
		Long:  `Pull a package from the registry. If no version is specified, latest will be used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse package name and version
			packageRef := args[0]
			packageName := packageRef
			packageVersion := "latest"

			// Extract version if provided
			if strings.Contains(packageRef, ":") {
				parts := strings.Split(packageRef, ":")
				packageName = parts[0]
				packageVersion = parts[1]
			}

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

			// Pull package
			fmt.Printf("Pulling package %s:%s from registry %s...\n", packageName, packageVersion, registryURL)
			downloadedPath, err := registryClient.Pull(context.Background(), packageName, packageVersion)
			if err != nil {
				return fmt.Errorf("failed to pull package: %w", err)
			}

			// Move to output path if specified
			if outputPath != "" {
				// TODO: Implement moving to output path
				fmt.Printf("Downloaded to %s (would move to %s)\n", downloadedPath, outputPath)
			} else {
				fmt.Printf("Downloaded to %s\n", downloadedPath)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for the downloaded package")
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")

	return cmd
}
