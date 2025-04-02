package system

import (
	"fmt"
	"time"

	"github.com/satishgonella2024/sentinelstacks/pkg/app"
	"github.com/satishgonella2024/sentinelstacks/pkg/models"
	"github.com/spf13/cobra"
)

// NewSystemCmd creates the system command group
func NewSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "View and manage system-wide resources",
		Long:  `View and manage system-wide resources and settings`,
	}

	// Add subcommands
	cmd.AddCommand(newSystemInfoCmd())
	cmd.AddCommand(newSystemDfCmd())
	cmd.AddCommand(newSystemPruneCmd())
	cmd.AddCommand(newSystemEventsCmd())

	return cmd
}

// newSystemInfoCmd creates the system info command
func newSystemInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Display system information",
		Long:  `Display detailed information about the SentinelStacks system`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			serviceRegistry := app.FromContext(ctx)
			
			// Get services
			networkService := serviceRegistry.NetworkService()
			volumeService := serviceRegistry.VolumeService()
			composeService := serviceRegistry.ComposeService()
			
			// Get resource counts
			networks, err := networkService.ListNetworks(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to list networks: %v\n", err)
			}
			
			volumes, err := volumeService.ListVolumes(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to list volumes: %v\n", err)
			}
			
			systems, err := composeService.ListSystems(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to list systems: %v\n", err)
			}
			
			// Count running systems
			runningCount := 0
			for _, system := range systems {
				if system.Status == "running" {
					runningCount++
				}
			}
			
			fmt.Println("SentinelStacks System Information:")
			
			fmt.Println("\nRuntime:")
			fmt.Printf("  Version:       %s\n", "v0.5.0")
			fmt.Printf("  API Version:   %s\n", "v1")
			fmt.Printf("  Build Date:    %s\n", "2023-10-25")
			fmt.Printf("  Go Version:    %s\n", "go1.21.2")
			fmt.Printf("  OS/Arch:       %s\n", "linux/amd64")
			
			fmt.Println("\nResources:")
			fmt.Printf("  Systems:       %d running / %d total\n", runningCount, len(systems))
			fmt.Printf("  Networks:      %d\n", len(networks))
			fmt.Printf("  Volumes:       %d\n", len(volumes))
			
			// Calculate total volume size
			totalVolumeSize := calculateTotalVolumeSize(volumes)
			
			fmt.Println("\nStorage:")
			fmt.Printf("  Volumes:       %s\n", formatSize(totalVolumeSize))
			
			fmt.Println("\nLLM Providers:")
			fmt.Printf("  Default:       %s\n", "claude (claude-3-5-sonnet-20240627)")
			fmt.Printf("  Configured:    %s\n", "claude, openai, ollama, google")
			
			return nil
		},
	}
}

// newSystemDfCmd creates the system df (disk free) command
func newSystemDfCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "df",
		Short: "Show disk usage",
		Long:  `Show disk usage by SentinelStacks components`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			serviceRegistry := app.FromContext(ctx)

			// Get services
			volumeService := serviceRegistry.VolumeService()
			
			// Get volumes
			volumes, err := volumeService.ListVolumes(ctx)
			if err != nil {
				return fmt.Errorf("failed to list volumes: %w", err)
			}
			
			fmt.Println("SentinelStacks Disk Usage:")
			
			fmt.Println("\nVolumes:")
			fmt.Printf("%-20s %-10s %-10s %-10s %s\n", "VOLUME NAME", "SIZE", "USED", "AVAILABLE", "MOUNTED BY")
			fmt.Printf("%-20s %-10s %-10s %-10s %s\n", "-----------", "----", "----", "---------", "----------")
			
			var totalSize, totalUsed int64
			
			for _, volume := range volumes {
				size := parseSize(volume.Size)
				used := parseSize(volume.Used)
				available := size - used
				
				mountedBy := "None"
				if volume.MountedBy != "" {
					mountedBy = volume.MountedBy
				}
				
				fmt.Printf("%-20s %-10s %-10s %-10s %s\n", 
					truncateString(volume.Name, 18),
					formatSize(size),
					formatSize(used),
					formatSize(available),
					mountedBy)
				
				totalSize += size
				totalUsed += used
			}
			
			fmt.Printf("\nTotal: %s of %s used\n", formatSize(totalUsed), formatSize(totalSize))
			
			return nil
		},
	}
}

// newSystemPruneCmd creates the system prune command
func newSystemPruneCmd() *cobra.Command {
	var (
		all     bool
		force   bool
		volumes bool
	)

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove unused data",
		Long:  `Remove unused data such as volumes and caches`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			serviceRegistry := app.FromContext(ctx)
			
			// Get services
			volumeService := serviceRegistry.VolumeService()
			
			if !force {
				fmt.Println("WARNING! This will remove:")
				if all {
					fmt.Println("  - all unused resources")
				} else {
					fmt.Println("  - all dangling resources")
				}
				
				if volumes {
					fmt.Println("  - all unused volumes")
				}
				
				fmt.Printf("\nThis may result in loss of data. Are you sure? [y/N] ")
				var response string
				fmt.Scanln(&response)
				
				if response != "y" && response != "Y" {
					fmt.Println("Aborting")
					return nil
				}
			}
			
			fmt.Println("Removing unused resources...")
			
			// Find and remove unused volumes
			if volumes {
				volumeList, err := volumeService.ListVolumes(ctx)
				if err != nil {
					return fmt.Errorf("failed to list volumes: %w", err)
				}
				
				fmt.Println("Removing unused volumes:")
				removedCount := 0
				
				for _, volume := range volumeList {
					// If volume is not mounted, consider it unused
					if volume.MountedBy == "" {
						fmt.Printf("  - %s (%s)\n", volume.Name, volume.Size)
						
						// Delete the volume
						if err := volumeService.DeleteVolume(ctx, volume.ID); err != nil {
							fmt.Printf("    Warning: Failed to delete volume: %v\n", err)
						} else {
							removedCount++
						}
					}
				}
				
				fmt.Printf("\nRemoved %d unused volumes\n", removedCount)
			}
			
			fmt.Println("Pruning completed successfully")
			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Remove all unused resources, not just dangling ones")
	cmd.Flags().BoolVar(&force, "force", false, "Do not prompt for confirmation")
	cmd.Flags().BoolVar(&volumes, "volumes", false, "Remove unused volumes as well")

	return cmd
}

// newSystemEventsCmd creates the system events command
func newSystemEventsCmd() *cobra.Command {
	var (
		since   string
		until   string
		filter  string
		limit   int
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "events",
		Short: "Get real-time events from the system",
		Long:  `Get real-time events from the SentinelStacks system with optional filtering`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			
			// Parse time filters
			sinceTime := parseTimeFilter(since)
			untilTime := parseTimeFilter(until)
			
			// For now, we'll create some sample events
			// In a real implementation, you'd get these from a system event log
			events := []systemEvent{
				{
					timestamp: time.Now().Add(-5 * time.Minute),
					eventType: "agent",
					subject:   "research-assistant",
					action:    "start",
					status:    "success",
					details: map[string]string{
						"id":     "agent123",
						"image":  "research-assistant:latest",
						"memory": "256MB",
					},
				},
				{
					timestamp: time.Now().Add(-10 * time.Minute),
					eventType: "volume",
					subject:   "research-memory",
					action:    "mount",
					status:    "success",
					details: map[string]string{
						"id":    "vol789",
						"agent": "agent123",
						"path":  "/memory",
					},
				},
			}
			
			// Filter events by time
			filtered := filterEventsByTime(events, sinceTime, untilTime)
			
			// Filter by type if specified
			if filter != "" {
				filtered = filterEventsByType(filtered, filter)
			}
			
			// Limit the number of events
			if len(filtered) > limit {
				filtered = filtered[:limit]
			}
			
			// Display events
			fmt.Println("System events:")
			if verbose {
				displayVerboseEvents(filtered)
			} else {
				displayEvents(filtered)
			}
			
			return nil
		},
	}

	cmd.Flags().StringVar(&since, "since", "1h", "Show events created since timestamp (e.g. 2h, 10m)")
	cmd.Flags().StringVar(&until, "until", "", "Show events created until timestamp")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter events by type (create, start, stop, etc.)")
	cmd.Flags().IntVar(&limit, "limit", 20, "Maximum number of events to show")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose event details")

	return cmd
}

// systemEvent represents a system event
type systemEvent struct {
	timestamp time.Time
	eventType string
	subject   string
	action    string
	status    string
	details   map[string]string
}

// filterEventsByTime filters events by time range
func filterEventsByTime(events []systemEvent, since, until time.Time) []systemEvent {
	filtered := []systemEvent{}
	for _, event := range events {
		if (since.IsZero() || event.timestamp.After(since)) &&
		   (until.IsZero() || event.timestamp.Before(until)) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// filterEventsByType filters events by type
func filterEventsByType(events []systemEvent, filter string) []systemEvent {
	filtered := []systemEvent{}
	for _, event := range events {
		if event.eventType == filter || event.action == filter {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// displayEvents displays events in a compact format
func displayEvents(events []systemEvent) {
	fmt.Printf("%-20s %-10s %-20s %-10s %-10s\n", "TIME", "TYPE", "SUBJECT", "ACTION", "STATUS")
	fmt.Println("-------------------------------------------------------------------------")
	
	for _, event := range events {
		fmt.Printf("%-20s %-10s %-20s %-10s %-10s\n",
			event.timestamp.Format("2006-01-02 15:04:05"),
			event.eventType,
			truncateString(event.subject, 18),
			event.action,
			event.status)
	}
}

// displayVerboseEvents displays events with full details
func displayVerboseEvents(events []systemEvent) {
	for i, event := range events {
		fmt.Printf("Event %d:\n", i+1)
		fmt.Printf("  Time:      %s\n", event.timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Type:      %s\n", event.eventType)
		fmt.Printf("  Subject:   %s\n", event.subject)
		fmt.Printf("  Action:    %s\n", event.action)
		fmt.Printf("  Status:    %s\n", event.status)
		
		fmt.Println("  Details:")
		for k, v := range event.details {
			fmt.Printf("    %-12s %s\n", k+":", v)
		}
		fmt.Println()
	}
}

// Helper functions

// calculateTotalVolumeSize calculates the total size of all volumes
func calculateTotalVolumeSize(volumes []*models.Volume) int64 {
	total := int64(0)
	for _, vol := range volumes {
		// Parse size (e.g., "1GB", "500MB")
		size := parseSize(vol.Size)
		total += size
	}
	return total
}

// parseSize parses a size string (e.g., "1GB", "500MB") and returns bytes
func parseSize(size string) int64 {
	// This is a simplified implementation
	// In a real implementation, you would parse the string properly
	return 1024 * 1024 * 1024 // Default to 1GB for now
}

// formatSize formats a byte size into a human-readable string
func formatSize(bytes int64) string {
	const (
		GB = 1024 * 1024 * 1024
		MB = 1024 * 1024
		KB = 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// parseTimeFilter parses a time filter (e.g., "2h", "30m") into a time.Time
func parseTimeFilter(filter string) time.Time {
	if filter == "" {
		return time.Time{} // Zero time means no filter
	}
	
	var duration time.Duration
	var err error
	
	// Try to parse as a duration
	if duration, err = time.ParseDuration(filter); err == nil {
		return time.Now().Add(-duration)
	}
	
	// Try to parse as a timestamp
	if t, err := time.Parse(time.RFC3339, filter); err == nil {
		return t
	}
	
	// Try to parse as a simple date
	if t, err := time.Parse("2006-01-02", filter); err == nil {
		return t
	}
	
	// If all parsing fails, return now (no filter)
	fmt.Printf("Warning: Could not parse time filter '%s', ignoring\n", filter)
	return time.Time{}
}

// truncateString truncates a string to the specified length and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
