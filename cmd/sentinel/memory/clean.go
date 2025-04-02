package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	
	"github.com/satishgonella2024/sentinelstacks/internal/memory"
)

// newMemoryCleanCmd creates a command to clean memory
func newMemoryCleanCmd() *cobra.Command {
	var force bool
	var olderThan string
	
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean memory storage",
		Long:  `Remove old or unused memory data`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get memory path
			memoryPath, err := memory.DefaultMemoryPath()
			if err != nil {
				return fmt.Errorf("failed to get memory path: %w", err)
			}
			
			// Check if memory path exists
			if _, err := os.Stat(memoryPath); os.IsNotExist(err) {
				fmt.Println("Memory directory does not exist")
				fmt.Println("Run 'sentinel memory init' to initialize memory subsystem")
				return nil
			}
			
			// Parse older-than parameter
			var cutoffTime time.Time
			if olderThan != "" {
				duration, err := parseDuration(olderThan)
				if err != nil {
					return fmt.Errorf("invalid duration format: %w", err)
				}
				cutoffTime = time.Now().Add(-duration)
			}
			
			// Confirmation
			if !force {
				fmt.Println("This operation will remove memory data that may be in use by agents or stacks.")
				fmt.Println("It is recommended to stop all running agents before cleaning memory.")
				fmt.Print("Do you want to continue? [y/N] ")
				
				var response string
				fmt.Scanln(&response)
				
				if strings.ToLower(response) != "y" {
					fmt.Println("Operation cancelled")
					return nil
				}
			}
			
			// Clean memory
			fmt.Println("Cleaning memory storage...")
			
			var removedFiles, removedDirs int
			var removedSize int64
			
			err = filepath.Walk(memoryPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				
				// Skip root directory and direct subdirectories (store types)
				if path == memoryPath {
					return nil
				}
				
				relPath, _ := filepath.Rel(memoryPath, path)
				parts := strings.Split(relPath, string(os.PathSeparator))
				if len(parts) == 1 {
					return nil
				}
				
				// Check if file is older than cutoff
				if olderThan != "" && !info.ModTime().Before(cutoffTime) {
					return nil
				}
				
				// Remove file or directory
				if !info.IsDir() {
					removedSize += info.Size()
					removedFiles++
					return os.Remove(path)
				} else {
					// Only remove empty directories
					entries, err := os.ReadDir(path)
					if err != nil {
						return err
					}
					
					if len(entries) == 0 {
						removedDirs++
						return os.Remove(path)
					}
				}
				
				return nil
			})
			
			if err != nil {
				return fmt.Errorf("failed to clean memory: %w", err)
			}
			
			// Print summary
			fmt.Printf("\nMemory cleaning complete:\n")
			fmt.Printf("  Removed files: %d\n", removedFiles)
			fmt.Printf("  Removed directories: %d\n", removedDirs)
			fmt.Printf("  Freed space: %s\n", formatSize(removedSize))
			
			return nil
		},
	}
	
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force cleaning without confirmation")
	cmd.Flags().StringVarP(&olderThan, "older-than", "o", "", "Only clean files older than specified duration (e.g., 7d, 24h, 30d)")
	
	return cmd
}

// parseDuration parses a human-readable duration like "7d" or "24h"
func parseDuration(durationStr string) (time.Duration, error) {
	// Check for day suffix
	if strings.HasSuffix(durationStr, "d") {
		days, err := parseInt(strings.TrimSuffix(durationStr, "d"))
		if err != nil {
			return 0, err
		}
		return time.Hour * 24 * time.Duration(days), nil
	}
	
	// Check for week suffix
	if strings.HasSuffix(durationStr, "w") {
		weeks, err := parseInt(strings.TrimSuffix(durationStr, "w"))
		if err != nil {
			return 0, err
		}
		return time.Hour * 24 * 7 * time.Duration(weeks), nil
	}
	
	// Parse using standard time.ParseDuration
	return time.ParseDuration(durationStr)
}

// parseInt parses an integer with error handling
func parseInt(s string) (int, error) {
	var value int
	_, err := fmt.Sscanf(s, "%d", &value)
	return value, err
}
