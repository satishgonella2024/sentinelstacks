package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/registry"
	"github.com/satishgonella2024/sentinelstacks/pkg/agent"
)

// NewImagesCmd creates a new images command
func NewImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List Sentinel Images",
		Long:  `List all Sentinel Images available locally.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read directly from the file system to debug
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}

			imagesDir := filepath.Join(homeDir, ".sentinel/images")
			
			// Check if directory exists
			if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
				fmt.Println("No images found. Registry directory does not exist.")
				return nil
			}
			
			// Read the directory
			entries, err := os.ReadDir(imagesDir)
			if err != nil {
				return fmt.Errorf("failed to read images directory: %w", err)
			}
			
			if len(entries) == 0 {
				fmt.Println("No images found. Registry directory is empty.")
				return nil
			}
			
			// Print the table header
			fmt.Println("IMAGE ID\tNAME\tTAG\tCREATED\tSIZE\tBASE MODEL")
			
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
				
				// Get file info for size
				fileInfo, err := entry.Info()
				if err != nil {
					fmt.Printf("Warning: Failed to get file info for %s: %v\n", filePath, err)
					continue
				}
				
				// Print a dummy line for each file to confirm we're finding files
				parts := strings.Split(entry.Name(), "_")
				if len(parts) >= 2 {
					name := strings.ReplaceAll(parts[0], "-", "/")
					tag := strings.TrimSuffix(parts[1], ".json")
					
					fmt.Printf("img_%x\t%s\t%s\t%s\t%.2f KB\t%s\n", 
						time.Now().UnixNano(), 
						name, 
						tag, 
						"Recent", 
						float64(fileInfo.Size())/1024,
						"unknown")
				}
			}
			
			return nil
		},
	}

	return cmd
}
