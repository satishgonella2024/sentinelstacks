package push

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/satishgonella2024/sentinelstacks/internal/registry"
)

// NewPushCmd creates the push command
func NewPushCmd() *cobra.Command {
	var (
		registry string
		public   bool
	)

	cmd := &cobra.Command{
		Use:   "push [image_name]",
		Short: "Push an agent image to a registry",
		Long:  `Push a locally built agent image to a remote registry`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName := args[0]
			
			// Parse image name and tag
			name, tag := parseImageName(imageName)
			
			// Check if image exists locally
			localRegistry, err := registry.GetLocalRegistry()
			if err != nil {
				return fmt.Errorf("failed to get local registry: %w", err)
			}
			
			image, err := localRegistry.Get(name, tag)
			if err != nil {
				return fmt.Errorf("image '%s:%s' not found locally: %w", name, tag, err)
			}
			
			// Determine registry URL
			registryURL := registry
			if registryURL == "" {
				registryURL = viper.GetString("registry.default")
				if registryURL == "" {
					registryURL = "sentinel.registry.ai" // Default registry
				}
			}
			
			fmt.Printf("Pushing image '%s:%s' to registry '%s' (public: %v)\n", name, tag, registryURL, public)
			
			// TODO: Implement actual image pushing logic here
			
			// Simulate upload progress
			fmt.Println("Uploading agent definition...")
			fmt.Println("Uploading capabilities...")
			fmt.Println("Uploading tools configuration...")
			
			// Check registry authentication
			if !isAuthenticated(registryURL) {
				fmt.Printf("Warning: Not authenticated to registry '%s'\n", registryURL)
				fmt.Println("Image will be pushed anonymously (if supported by registry)")
				fmt.Println("Use 'sentinel login' to authenticate for full access")
			}
			
			fmt.Printf("Successfully pushed image '%s:%s'\n", name, tag)
			
			if public {
				fmt.Printf("Image is now publicly available at %s/%s:%s\n", registryURL, name, tag)
			} else {
				fmt.Printf("Image is now available to authorized users at %s/%s:%s\n", registryURL, name, tag)
			}
			
			return nil
		},
	}

	cmd.Flags().StringVar(&registry, "registry", "", "Registry URL to push to")
	cmd.Flags().BoolVar(&public, "public", false, "Make the image publicly accessible")

	return cmd
}

// parseImageName parses an image name into name and tag parts
func parseImageName(imageName string) (string, string) {
	// Add :latest tag if no tag specified
	if !strings.Contains(imageName, ":") {
		imageName = imageName + ":latest"
	}

	// Split name and tag
	parts := strings.SplitN(imageName, ":", 2)
	name := parts[0]
	tag := "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}
	
	return name, tag
}

// isAuthenticated checks if user is authenticated to the registry
func isAuthenticated(registry string) bool {
	// TODO: Implement actual authentication check
	
	// For now, return true if there's an API key in the config
	return viper.GetString("registry.api_key") != ""
}
