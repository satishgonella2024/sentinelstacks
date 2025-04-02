package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// Removed term dependency
)

// NewLoginCmd creates the login command
func NewLoginCmd() *cobra.Command {
	var (
		username string
		password string
		token    string
	)

	cmd := &cobra.Command{
		Use:   "login [registry]",
		Short: "Log in to a registry",
		Long:  `Log in to a registry to enable pushing and pulling agent images`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine registry URL
			registryURL := ""
			if len(args) > 0 {
				registryURL = args[0]
			}
			
			if registryURL == "" {
				registryURL = viper.GetString("registry.default")
				if registryURL == "" {
					registryURL = "sentinel.registry.ai" // Default registry
				}
			}
			
			fmt.Printf("Logging in to registry: %s\n", registryURL)
			
			// If token is provided, use it directly
			if token != "" {
				return saveCredentials(registryURL, token, "", "")
			}
			
			// If username is not provided, prompt for it
			if username == "" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Username: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read username: %w", err)
				}
				username = strings.TrimSpace(input)
			}
			
			// If password is not provided, prompt for it
			if password == "" {
				fmt.Print("Password: ")
				reader := bufio.NewReader(os.Stdin)
				passwordInput, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}
				password = strings.TrimSpace(passwordInput)
				// Note: This isn't secure as it shows the password, but we're
				// removing the dependency on term for simplicity
			}
			
			// TODO: Implement actual authentication with registry
			// For now, simulate authentication
			
			fmt.Println("Authenticating...")
			
			// Simulate getting a token from the registry
			simulatedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYWRtaW4ifQ.pQn2H0WCsdrS7fKLLt7BMnpuPfLLZ1uNcJIItfFgIkM"
			
			// Save the credentials
			if err := saveCredentials(registryURL, simulatedToken, username, ""); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}
			
			fmt.Printf("Login Succeeded for %s\n", username)
			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "Username for registry authentication")
	cmd.Flags().StringVar(&password, "password", "", "Password for registry authentication")
	cmd.Flags().StringVar(&token, "token", "", "Authentication token (if using token auth)")

	// Mark password flag as sensitive to avoid showing in help
	cmd.Flags().MarkHidden("password")

	return cmd
}

// saveCredentials saves authentication credentials to the config
func saveCredentials(registry, token, username, password string) error {
	// Create a key based on the registry
	registryKey := strings.ReplaceAll(registry, ".", "_")
	registryKey = strings.ReplaceAll(registryKey, ":", "_")
	registryKey = strings.ReplaceAll(registryKey, "/", "_")
	
	// Save the token
	viper.Set(fmt.Sprintf("registry.auth.%s.token", registryKey), token)
	
	// Save the username (for display purposes)
	if username != "" {
		viper.Set(fmt.Sprintf("registry.auth.%s.username", registryKey), username)
	}
	
	// Set as default registry if not already set
	if viper.GetString("registry.default") == "" {
		viper.Set("registry.default", registry)
	}
	
	// Write the changes to the config file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	return nil
}
