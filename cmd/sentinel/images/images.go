package images

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/internal/registry"
)

// ImageInfo contains information about a Sentinel Image
type ImageInfo struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Tag        string                 `json:"tag"`
	CreatedAt  time.Time              `json:"createdAt"`
	Size       int64                  `json:"size"`
	BaseModel  string                 `json:"baseModel"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// NewImagesCmd creates a new images command
func NewImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List Sentinel Images",
		Long:  `List all Sentinel Images available locally.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			quiet, _ := cmd.Flags().GetBool("quiet")
			format, _ := cmd.Flags().GetString("format")
			filter, _ := cmd.Flags().GetString("filter")

			return runImages(all, quiet, format, filter)
		},
	}

	// Add flags
	cmd.Flags().BoolP("all", "a", false, "Show all images (default hides intermediate images)")
	cmd.Flags().BoolP("quiet", "q", false, "Only display image IDs")
	cmd.Flags().StringP("format", "f", "", "Format the output using a custom format")
	cmd.Flags().String("filter", "", "Filter output based on conditions")

	return cmd
}

// runImages executes the images command
func runImages(all, quiet bool, format, filter string) error {
	// Get local registry
	_, err := registry.GetLocalRegistry()
	if err != nil {
		return fmt.Errorf("failed to get local registry: %w", err)
	}

	// List images
	var images []ImageInfo

	// For now, create sample images
	// In a real implementation, this would come from the registry.List() method
	images = getSampleImages()

	// Apply filter if specified
	if filter != "" {
		filteredImages := []ImageInfo{}
		filterParts := strings.Split(filter, "=")
		if len(filterParts) != 2 {
			return fmt.Errorf("invalid filter format, expected key=value")
		}

		key := strings.TrimSpace(filterParts[0])
		value := strings.TrimSpace(filterParts[1])

		for _, img := range images {
			switch key {
			case "name":
				if strings.Contains(img.Name, value) {
					filteredImages = append(filteredImages, img)
				}
			case "tag":
				if img.Tag == value {
					filteredImages = append(filteredImages, img)
				}
			case "baseModel":
				if img.BaseModel == value {
					filteredImages = append(filteredImages, img)
				}
			}
		}

		images = filteredImages
	}

	// Handle quiet mode (only show IDs)
	if quiet {
		for _, img := range images {
			fmt.Println(img.ID[:12])
		}
		return nil
	}

	// Handle custom format
	if format != "" {
		return formatOutput(images, format)
	}

	// Default output format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "IMAGE ID\tNAME\tTAG\tCREATED\tSIZE\tBASE MODEL")

	for _, img := range images {
		id := img.ID
		if len(id) > 12 {
			id = id[:12]
		}

		size := formatSize(img.Size)
		createdTime := formatTime(img.CreatedAt)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			id,
			img.Name,
			img.Tag,
			createdTime,
			size,
			img.BaseModel,
		)
	}

	return w.Flush()
}

// formatOutput applies a custom format to the output
func formatOutput(images []ImageInfo, format string) error {
	for _, img := range images {
		line := format

		// Replace placeholders with actual values
		line = strings.ReplaceAll(line, "{{.ID}}", img.ID)
		line = strings.ReplaceAll(line, "{{.Name}}", img.Name)
		line = strings.ReplaceAll(line, "{{.Tag}}", img.Tag)
		line = strings.ReplaceAll(line, "{{.CreatedAt}}", img.CreatedAt.Format(time.RFC3339))
		line = strings.ReplaceAll(line, "{{.Size}}", fmt.Sprintf("%d", img.Size))
		line = strings.ReplaceAll(line, "{{.BaseModel}}", img.BaseModel)

		fmt.Println(line)
	}

	return nil
}

// formatTime formats the time relative to now
func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Less than a minute ago"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	} else if diff < 48*time.Hour {
		return "Yesterday"
	} else {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	}
}

// formatSize formats the size in human-readable format
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}

// plural returns "s" if the number is not 1
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// getSampleImages creates sample images for demonstration purposes
func getSampleImages() []ImageInfo {
	now := time.Now()

	return []ImageInfo{
		{
			ID:        "sha256:abcdef1234567890",
			Name:      "user/chatbot",
			Tag:       "latest",
			CreatedAt: now.Add(-3 * time.Hour),
			Size:      1024 * 1024 * 5, // 5 MB
			BaseModel: "claude-3-haiku-20240307",
			Parameters: map[string]interface{}{
				"temperature": 0.7,
				"memoryDepth": 10,
			},
		},
		{
			ID:        "sha256:9876543210abcdef",
			Name:      "user/research-assistant",
			Tag:       "v1.0",
			CreatedAt: now.Add(-1 * 24 * time.Hour),
			Size:      1024 * 1024 * 8, // 8 MB
			BaseModel: "claude-3-opus-20240229",
			Parameters: map[string]interface{}{
				"temperature": 0.5,
				"memoryDepth": 20,
			},
		},
		{
			ID:        "sha256:fedcba0987654321",
			Name:      "user/translator",
			Tag:       "v2.1",
			CreatedAt: now.Add(-5 * 24 * time.Hour),
			Size:      1024 * 1024 * 3, // 3 MB
			BaseModel: "claude-3-sonnet-20240229",
			Parameters: map[string]interface{}{
				"temperature": 0.3,
				"memoryDepth": 5,
			},
		},
		{
			ID:        "sha256:1234567890abcdef",
			Name:      "user/codehelper",
			Tag:       "latest",
			CreatedAt: now.Add(-12 * time.Hour),
			Size:      1024 * 1024 * 6, // 6 MB
			BaseModel: "claude-3-opus-20240229",
			Parameters: map[string]interface{}{
				"temperature": 0.2,
				"memoryDepth": 15,
			},
		},
	}
}
