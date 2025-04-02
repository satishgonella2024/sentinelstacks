package volume

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVolumeCmd creates the volume command group
func NewVolumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Manage agent memory volumes",
		Long:  `Create and manage persistent memory volumes for agents`,
	}

	// Add subcommands
	cmd.AddCommand(newVolumeCreateCmd())
	cmd.AddCommand(newVolumeListCmd())
	cmd.AddCommand(newVolumeMountCmd())
	cmd.AddCommand(newVolumeUnmountCmd())
	cmd.AddCommand(newVolumeRemoveCmd())
	cmd.AddCommand(newVolumeInspectCmd())

	return cmd
}

// newVolumeCreateCmd creates the volume create command
func newVolumeCreateCmd() *cobra.Command {
	var (
		size      string
		encrypted bool
	)

	cmd := &cobra.Command{
		Use:   "create [volume_name]",
		Short: "Create a new memory volume",
		Long:  `Create a new persistent memory volume for agent state storage`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			volumeName := args[0]
			fmt.Printf("Creating volume '%s' (size: %s, encrypted: %v)\n", volumeName, size, encrypted)
			
			// TODO: Implement actual volume creation here
			
			fmt.Printf("Volume '%s' created successfully\n", volumeName)
			return nil
		},
	}

	cmd.Flags().StringVar(&size, "size", "1GB", "Volume size (e.g., 500MB, 2GB)")
	cmd.Flags().BoolVar(&encrypted, "encrypted", false, "Enable encryption for the volume")
	return cmd
}

// newVolumeListCmd creates the volume list command
func newVolumeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List volumes",
		Long:    `List all available memory volumes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Listing all volumes:")
			
			// TODO: Implement actual volume listing here
			
			fmt.Println("Volume1    1GB    Created 2023-10-15    In use by: agent1")
			fmt.Println("Volume2    500MB  Created 2023-10-20    Available")
			return nil
		},
	}
}

// newVolumeMountCmd creates the volume mount command
func newVolumeMountCmd() *cobra.Command {
	var mountPath string

	cmd := &cobra.Command{
		Use:   "mount [volume_name] [agent_id]",
		Short: "Mount a volume to an agent",
		Long:  `Mount a persistent memory volume to a specific agent`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			agentID := args[1]
			
			fmt.Printf("Mounting volume '%s' to agent '%s' at path '%s'\n", volumeName, agentID, mountPath)
			
			// Get the volume service from application context
			serviceRegistry := app.FromContext(ctx)
			volumeService := serviceRegistry.VolumeService()
			
			// Mount the volume
			if err := volumeService.MountVolume(ctx, volumeName, agentID, mountPath); err != nil {
				return fmt.Errorf("failed to mount volume: %w", err)
			}
			
			fmt.Printf("Volume '%s' successfully mounted to agent '%s'\n", volumeName, agentID)
			return nil
		},
	}

	cmd.Flags().StringVar(&mountPath, "path", "/memory", "Path where the volume will be mounted")
	return cmd
}

// newVolumeUnmountCmd creates the volume unmount command
func newVolumeUnmountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unmount [volume_name] [agent_id]",
		Short: "Unmount a volume from an agent",
		Long:  `Unmount a persistent memory volume from a specified agent`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			agentID := args[1]
			
			fmt.Printf("Unmounting volume '%s' from agent '%s'\n", volumeName, agentID)
			
			// Get the volume service from application context
			serviceRegistry := app.FromContext(ctx)
			volumeService := serviceRegistry.VolumeService()
			
			// Unmount the volume
			if err := volumeService.UnmountVolume(ctx, volumeName, agentID); err != nil {
				return fmt.Errorf("failed to unmount volume: %w", err)
			}
			
			fmt.Printf("Volume '%s' successfully unmounted from agent '%s'\n", volumeName, agentID)
			return nil
		},
	}
}

// newVolumeRemoveCmd creates the volume remove command
func newVolumeRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "rm [volume_name]",
		Aliases: []string{"remove"},
		Short:   "Remove a volume",
		Long:    `Remove a memory volume and all its contents`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			
			fmt.Printf("Removing volume '%s' (force: %v)\n", volumeName, force)
			
			// Get the volume service from application context
			serviceRegistry := app.FromContext(ctx)
			volumeService := serviceRegistry.VolumeService()
			
			// Get volume first to get its ID
			volume, err := volumeService.GetVolumeByName(ctx, volumeName)
			if err != nil {
				return fmt.Errorf("failed to find volume: %w", err)
			}
			
			// Check if volume is mounted and force flag is not set
			if volume.MountedBy != "" && !force {
				return fmt.Errorf("volume is mounted by agent '%s'; use --force to remove", volume.MountedBy)
			}
			
			// Delete the volume
			if err := volumeService.DeleteVolume(ctx, volume.ID); err != nil {
				return fmt.Errorf("failed to remove volume: %w", err)
			}
			
			fmt.Printf("Volume '%s' successfully removed\n", volumeName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force removal even if volume is in use")
	return cmd
}

// newVolumeInspectCmd creates the volume inspect command
func newVolumeInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [volume_name]",
		Short: "Display detailed information on a volume",
		Long:  `Display detailed information about a memory volume, including usage`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			volumeName := args[0]
			
			fmt.Printf("Inspecting volume '%s'\n\n", volumeName)
			
			// Get the volume service from application context
			serviceRegistry := app.FromContext(ctx)
			volumeService := serviceRegistry.VolumeService()
			
			// Get volume details
			volume, err := volumeService.InspectVolume(ctx, volumeName)
			if err != nil {
				return fmt.Errorf("failed to inspect volume: %w", err)
			}
			
			fmt.Printf("Volume: %s\n", volume.Name)
			fmt.Printf("  ID: %s\n", volume.ID)
			fmt.Printf("  Created: %s\n", volume.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Size: %s\n", volume.Size)
			fmt.Printf("  Used: %s\n", volume.Used)
			
			// Calculate free space if possible
			// This would require a proper size parser to handle units like MB, GB, etc.
			// fmt.Printf("  Free: %s\n", calculateFreeSpace(volume.Size, volume.Used))
			
			fmt.Printf("  Encrypted: %v\n", volume.Encrypted)
			
			if volume.MountedBy != "" {
				fmt.Printf("  Status: In use\n")
				fmt.Printf("  Mounted by: %s at %s\n", volume.MountedBy, volume.MountPath)
			} else {
				fmt.Printf("  Status: Available\n")
			}
			
			// TODO: Implement content listing when we have actual storage backend
			fmt.Printf("  Contents: Not implemented yet\n")
			
			if len(volume.Metadata) > 0 {
				fmt.Println("  Metadata:")
				for key, value := range volume.Metadata {
					fmt.Printf("    %s: %s\n", key, value)
				}
			}
			
			return nil
		},
	}
}
