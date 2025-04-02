package registry

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/pkg/registry"
)

// NewPullCmd creates a new pull command
func NewPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [agent_image:tag]",
		Short: "Pull an agent image from a registry",
		Long:  `Pull an agent image from a registry to the local environment.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(args[0])
		},
	}

	return cmd
}

func runPull(imageTag string) error {
	// Validate image tag format
	if !strings.Contains(imageTag, ":") {
		imageTag += ":latest" // Use latest tag if none specified
	}

	parts := strings.Split(imageTag, ":")
	imageName := parts[0]
	tag := parts[1]

	// Pull the image
	if err := registry.PullImage(imageName, tag); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	fmt.Printf("Successfully pulled image %s\n", imageTag)
	return nil
}
