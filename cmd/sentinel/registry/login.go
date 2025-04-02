package registry

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/sentinelstacks/sentinel/pkg/registry"
)

// NewLoginCmd creates a new login command
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [registry_url]",
		Short: "Log in to an agent registry",
		Long:  `Log in to an agent registry to enable pushing and pulling private agents.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var registryURL string
			if len(args) > 0 {
				registryURL = args[0]
			} else {
				registryURL = registry.DefaultRegistryURL
			}
			
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			token, _ := cmd.Flags().GetString("token")
			
			return runLogin(registryURL, username, password, token)
		},
	}

	// Add flags
	cmd.Flags().StringP("username", "u", "", "Username for registry authentication")
	cmd.Flags().StringP("password", "p", "", "Password for registry authentication (not recommended, use interactive prompt)")
	cmd.Flags().String("token", "", "Authentication token (alternative to username/password)")

	return cmd
}

func runLogin(registryURL, username, password, token string) error {
	// If token is provided, use that for authentication
	if token != "" {
		if err := registry.LoginWithToken(registryURL, token); err != nil {
			return fmt.Errorf("failed to log in to registry: %w", err)
		}
		fmt.Printf("Successfully logged in to %s\n", registryURL)
		return nil
	}

	// If username is not provided, prompt for it
	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	// If password is not provided, prompt for it
	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Add a newline after password input
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = string(passwordBytes)
	}

	// Validate input
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Attempt login
	if err := registry.Login(registryURL, username, password); err != nil {
		return fmt.Errorf("failed to log in to registry: %w", err)
	}

	fmt.Printf("Successfully logged in to %s as %s\n", registryURL, username)
	return nil
}
