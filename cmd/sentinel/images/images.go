package images

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry"
)

// NewImagesCmd creates the images command
func NewImagesCmd() *cobra.Command {
	var (
		format string
	)

	imagesCmd := &cobra.Command{
		Use:   "images [filter]",
		Short: "List Sentinel Agent images",
		Long:  `List all Sentinel Agent images in the local registry`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get filter if provided
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}

			// Get local registry
			reg, err := registry.GetLocalRegistry()
			if err != nil {
				return fmt.Errorf("failed to get registry: %w", err)
			}

			// Get images from registry
			imageList, err := reg.List()
			if err != nil {
				return fmt.Errorf("failed to list images: %w", err)
			}

			// Filter images if filter is provided
			if filter != "" {
				var filtered []registry.ImageSummary
				for _, img := range imageList {
					if strings.Contains(img.Name, filter) || strings.Contains(img.Tag, filter) {
						filtered = append(filtered, img)
					}
				}
				imageList = filtered
			}

			// Print the images based on format
			switch format {
			case "json":
				return printImagesJSON(imageList)
			case "wide":
				return printImagesWide(imageList)
			default:
				return printImagesTable(imageList)
			}
		},
	}

	imagesCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, wide, json)")

	return imagesCmd
}

// printImagesTable prints images in a table format
func printImagesTable(images []registry.ImageSummary) error {
	// Print header
	fmt.Printf("%-40s %-15s %-20s %-10s\n", "IMAGE", "TAG", "CREATED", "SIZE")
	fmt.Println(strings.Repeat("-", 90))

	// Print images
	for _, img := range images {
		// Format the creation time
		created := formatTime(img.Created)

		// Format size
		size := formatSize(img.Size)

		fmt.Printf("%-40s %-15s %-20s %-10s\n",
			img.Name,
			img.Tag,
			created,
			size)
	}

	// Print total
	fmt.Printf("\nTotal: %d images\n", len(images))
	return nil
}

// printImagesWide prints images in a wide format with more details
func printImagesWide(images []registry.ImageSummary) error {
	// Print header
	fmt.Printf("%-40s %-15s %-20s %-10s %-12s %-30s\n",
		"IMAGE", "TAG", "CREATED", "SIZE", "BASE MODEL", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 130))

	// Print images
	for _, img := range images {
		// Format the creation time
		created := formatTime(img.Created)

		// Format size
		size := formatSize(img.Size)

		// Get base model
		baseModel := img.BaseModel
		if baseModel == "" {
			baseModel = "N/A"
		}

		// Truncate description if too long
		desc := img.Description
		if len(desc) > 30 {
			desc = desc[:27] + "..."
		}

		fmt.Printf("%-40s %-15s %-20s %-10s %-12s %-30s\n",
			img.Name,
			img.Tag,
			created,
			size,
			baseModel,
			desc)
	}

	// Print total
	fmt.Printf("\nTotal: %d images\n", len(images))
	return nil
}

// printImagesJSON prints images in JSON format
func printImagesJSON(images []registry.ImageSummary) error {
	// Use registry's builtin JSON formatter
	jsonData, err := registry.FormatImagesAsJSON(images)
	if err != nil {
		return fmt.Errorf("failed to format images as JSON: %w", err)
	}

	fmt.Println(jsonData)
	return nil
}

// formatTime formats a time.Time as a human-readable string
func formatTime(t time.Time) string {
	// If time is zero, return "N/A"
	if t.IsZero() {
		return "N/A"
	}

	// Get time difference
	diff := time.Since(t)

	// Format based on how long ago it was
	if diff < time.Hour {
		// Less than an hour ago
		mins := int(diff.Minutes())
		if mins < 1 {
			return "Just now"
		}
		return fmt.Sprintf("%d minute(s) ago", mins)
	} else if diff < 24*time.Hour {
		// Less than a day ago
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour(s) ago", hours)
	} else if diff < 7*24*time.Hour {
		// Less than a week ago
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day(s) ago", days)
	} else if diff < 30*24*time.Hour {
		// Less than a month ago
		weeks := int(diff.Hours() / 24 / 7)
		return fmt.Sprintf("%d week(s) ago", weeks)
	} else {
		// More than a month ago, use exact date
		return t.Format("2006-01-02 15:04:05")
	}
}

// formatSize formats a size in bytes as a human-readable string
func formatSize(sizeBytes int64) string {
	// If size is zero, return "N/A"
	if sizeBytes == 0 {
		return "N/A"
	}

	// Format based on size
	const unit = 1024
	if sizeBytes < unit {
		return fmt.Sprintf("%d B", sizeBytes)
	}
	div, exp := int64(unit), 0
	for n := sizeBytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(sizeBytes)/float64(div), "KMGTPE"[exp])
}
