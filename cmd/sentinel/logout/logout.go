package logout

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewLogoutCmd creates the logout command
func NewLogoutCmd() *cobra.Command {
	var allRegistries bool

	cmd := &cobra.Command{
		Use:   "logout [registry]",
		Short: "Log out from a registry",
		Long:  `Log out from a registry to remove saved credentials`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If --all flag is set, log out from all registries
			if allRegistries {
				return logoutAllRegistries()
			}
			
			// Determine registry URL
			registryURL := ""
			if len(args) > 0 {
				registryURL = args[0]
			}
			
			if registryURL == "" {
				registryURL = viper.GetString("registry.default")
				if registryURL == "" {
					return fmt.Errorf("no registry specified and no default registry set")
				}
			}
			
			fmt.Printf("Logging out from registry: %s\n", registryURL)
			
			// Get the registry key
			registryKey := getRegistryKey(registryURL)
			
			// Get the username for display purposes
			username := viper.GetString(fmt.Sprintf("registry.auth.%s.username", registryKey))
			if username == "" {
				username = "current user"
			}
			
			// Remove the credentials
			viper.Set(fmt.Sprintf("registry.auth.%s", registryKey), nil)
			
			// Write the changes to the config file
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
			
			fmt.Printf("Removed login credentials for %s from %s\n", username, registryURL)
			return nil
		},
	}

	cmd.Flags().BoolVar(&allRegistries, "all", false, "Log out from all registries")

	return cmd
}

// logoutAllRegistries logs out from all registries
func logoutAllRegistries() error {
	// Get all registry entries
	authEntries := viper.GetStringMap("registry.auth")
	
	if len(authEntries) == 0 {
		fmt.Println("No registry logins found")
		return nil
	}
	
	// Remove all entries
	viper.Set("registry.auth", nil)
	
	// Write the changes to the config file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	fmt.Printf("Logged out from %d registries\n", len(authEntries))
	return nil
}

// getRegistryKey converts a registry URL to a config key
func getRegistryKey(registry string) string {
	registryKey := strings.ReplaceAll(registry, ".", "_")
	registryKey = strings.ReplaceAll(registryKey, ":", "_")
	registryKey = strings.ReplaceAll(registryKey, "/", "_")
	return registryKey
}
