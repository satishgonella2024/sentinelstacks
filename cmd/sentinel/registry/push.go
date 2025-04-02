package registry

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/pkg/registry"
)

// NewPushCmd creates a new push command
func NewPushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [agent_image:tag]",
		Short: "Push an agent image to a registry",
		Long:  `Push an agent image to a registry for sharing or deployment.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPush(args[0])
		},
	}

	return cmd
}

func runPush(imageTag string) error {
	// Validate image tag format
	if !strings.Contains(imageTag, ":") {
		imageTag += ":latest" // Use latest tag if none specified
	}

	parts := strings.Split(imageTag, ":")
	imageName := parts[0]
	tag := parts[1]

	// Push the image
	if err := registry.PushImage(imageName, tag); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("Successfully pushed image %s\n", imageTag)
	return nil
}
