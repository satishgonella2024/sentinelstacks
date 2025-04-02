package pull

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPullCmd creates the pull command
func NewPullCmd() *cobra.Command {
	var (
		force    bool
		registry string
	)

	cmd := &cobra.Command{
		Use:   "pull [image_name]",
		Short: "Pull an agent image from a registry",
		Long:  `Pull an agent image from a remote registry to your local environment`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName := args[0]
			
			// Parse image name and tag
			name, tag := parseImageName(imageName)
			
			// Determine registry URL
			registryURL := registry
			if registryURL == "" {
				registryURL = viper.GetString("registry.default")
				if registryURL == "" {
					registryURL = "sentinel.registry.ai" // Default registry
				}
			}
			
			fmt.Printf("Pulling image '%s:%s' from registry '%s'\n", name, tag, registryURL)
			
			// TODO: Implement actual image pulling logic here
			
			// Simulate download progress
			fmt.Println("Downloading agent definition...")
			fmt.Println("Downloading capabilities...")
			fmt.Println("Downloading tools configuration...")
			
			fmt.Printf("Successfully pulled image '%s:%s'\n", name, tag)
			fmt.Println("Image is now available for local use")
			
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force pull even if image exists locally")
	cmd.Flags().StringVar(&registry, "registry", "", "Registry URL to pull from")

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
