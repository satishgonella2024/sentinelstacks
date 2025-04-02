package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Image represents a Sentinel Image with minimal fields
type Image struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Tag        string     `json:"tag"`
	CreatedAt  int64      `json:"createdAt"`
	Definition Definition `json:"definition"`
}

// Definition represents a simplified agent definition
type Definition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseModel   string `json:"baseModel"`
}

func main() {
	// Get the registry directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	imagesDir := filepath.Join(homeDir, ".sentinel/images")

	// Check if directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		fmt.Println("No images found. Registry directory does not exist.")
		os.Exit(0)
	}

	// Read the directory
	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read images directory: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No images found. Registry directory is empty.")
		os.Exit(0)
	}

	// Print the table header
	fmt.Printf("%-12s %-30s %-10s %-20s %-10s %-20s\n", "IMAGE ID", "NAME", "TAG", "CREATED", "SIZE", "BASE MODEL")
	fmt.Println(strings.Repeat("-", 100))

	// Process each file
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Read the file
		filePath := filepath.Join(imagesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: Failed to read %s: %v\n", filePath, err)
			continue
		}

		// Parse the JSON
		var image Image
		if err := json.Unmarshal(data, &image); err != nil {
			// Fall back to filename parsing if JSON parsing fails
			parts := strings.Split(entry.Name(), "_")
			if len(parts) >= 2 {
				name := strings.ReplaceAll(parts[0], "-", "/")
				tag := strings.TrimSuffix(parts[1], ".json")
				
				// Generate an ID based on the file path
				hasher := md5.New()
				hasher.Write([]byte(filePath))
				id := hex.EncodeToString(hasher.Sum(nil))[:12]
				
				// Get file info for size
				fileInfo, err := entry.Info()
				if err != nil {
					fmt.Printf("Warning: Failed to get file info for %s: %v\n", filePath, err)
					continue
				}
				
				fmt.Printf("%-12s %-30s %-10s %-20s %-10s %-20s\n", 
					id, 
					name, 
					tag, 
					"Unknown", 
					formatSize(fileInfo.Size()),
					"unknown")
			}
			continue
		}

		// Get file info for size
		fileInfo, err := entry.Info()
		if err != nil {
			fmt.Printf("Warning: Failed to get file info for %s: %v\n", filePath, err)
			continue
		}

		// Format the created time
		createdTime := time.Unix(image.CreatedAt, 0)
		timeStr := formatTime(createdTime)

		// Extract ID (use first 12 chars if it's long enough)
		id := image.ID
		if len(id) > 12 {
			id = id[:12]
		}

		// Print the image details
		fmt.Printf("%-12s %-30s %-10s %-20s %-10s %-20s\n",
			id,
			image.Name,
			image.Tag,
			timeStr,
			formatSize(fileInfo.Size()),
			image.Definition.BaseModel)
	}
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
