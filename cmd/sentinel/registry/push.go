package registry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
)

// NewPushCmd creates a new push command
func NewPushCmd() *cobra.Command {
	var (
		username string
		password string
		token    string
	)

	cmd := &cobra.Command{
		Use:   "push [package-path]",
		Short: "Push a package to the registry",
		Long:  `Push a package to the registry. Authentication is required.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get package path
			packagePath := args[0]

			// Check if file exists
			info, err := os.Stat(packagePath)
			if err != nil {
				return fmt.Errorf("invalid package path: %w", err)
			}

			// Ensure it's a file
			if info.IsDir() {
				return errors.New("package path must be a file, not a directory")
			}

			// Check file extension
			if !strings.HasSuffix(packagePath, ".sentinel-pkg") && !strings.HasSuffix(packagePath, ".tar.gz") {
				return errors.New("package file must have .sentinel-pkg or .tar.gz extension")
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

			// Check if user needs to authenticate
			ctx := context.Background()
			authenticated := authProvider.IsAuthenticated(ctx)

			if !authenticated {
				// Try to authenticate
				if username != "" && password != "" {
					// Login with username/password
					fmt.Println("Authenticating with username and password...")
					if _, err := authProvider.Login(ctx, username, password); err != nil {
						return fmt.Errorf("authentication failed: %w", err)
					}
				} else if token != "" {
					// TODO: Implement token-based auth directly
					return errors.New("token-based authentication not implemented yet")
				} else {
					// Prompt for credentials
					return errors.New("authentication required. Use --username and --password flags or login first with 'sentinel registry login'")
				}
			}

			// Create registry client
			clientConfig := client.Config{
				BaseURL:      registryURL,
				AuthProvider: authProvider,
			}
			registryClient := client.NewClient(clientConfig)

			// Push package
			fmt.Printf("Pushing package %s to registry %s...\n", filepath.Base(packagePath), registryURL)
			if err := registryClient.Push(ctx, packagePath); err != nil {
				return fmt.Errorf("failed to push package: %w", err)
			}

			fmt.Println("Package pushed successfully!")
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&username, "username", "", "Registry username")
	cmd.Flags().StringVar(&password, "password", "", "Registry password")
	cmd.Flags().StringVar(&token, "token", "", "Auth token (alternative to username/password)")
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")

	return cmd
}
