package registry

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/auth"
)

// NewLoginCmd creates a new login command
func NewLoginCmd() *cobra.Command {
	var (
		username string
		password string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a registry",
		Long:  `Log in to a registry to push or pull packages.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get registry URL from flag or env var
			registryURL, _ := cmd.Flags().GetString("registry")
			if registryURL == "" {
				registryURL = os.Getenv("SENTINEL_REGISTRY_URL")
			}
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Check if username and password are provided
			if username == "" {
				// Prompt for username
				fmt.Print("Username: ")
				fmt.Scanln(&username)
			}

			if password == "" {
				// Prompt for password
				fmt.Print("Password: ")
				passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}
				fmt.Println() // Add newline after password input
				password = string(passwordBytes)
			}

			// Create auth provider
			authConfig := auth.AuthConfig{
				RegistryURL: registryURL,
			}
			authProvider, err := auth.NewFileTokenProvider(authConfig)
			if err != nil {
				return fmt.Errorf("failed to create auth provider: %w", err)
			}

			// Login
			ctx := context.Background()
			token, err := authProvider.Login(ctx, username, password)
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			fmt.Printf("Successfully logged in to registry %s\n", registryURL)

			// Print token expiration info
			if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
				fmt.Printf("Token: %s\n", token)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&username, "username", "", "Registry username")
	cmd.Flags().StringVar(&password, "password", "", "Registry password")
	cmd.Flags().String("registry", "", "Registry URL (default: https://registry.sentinelstacks.io)")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	return cmd
}
