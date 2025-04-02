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

// NewMemoryCmd creates a new memory command
func NewMemoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Manage memory storage",
		Long:  `Manage memory storage for agents and stacks`,
	}
	
	// Add subcommands
	cmd.AddCommand(newMemoryInitCmd())
	cmd.AddCommand(newMemoryListCmd())
	cmd.AddCommand(newMemoryCleanCmd())
	cmd.AddCommand(newMemoryInfoCmd())
	
	return cmd
}

// newMemoryInitCmd creates a command to initialize memory
func newMemoryInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize memory subsystem",
		Long:  `Initialize memory subsystem by creating required directories`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing memory subsystem...")
			
			err := memory.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize memory: %w", err)
			}
			
			fmt.Println("Memory subsystem initialized successfully")
			return nil
		},
	}
	
	return cmd
}

// newMemoryListCmd creates a command to list memory content
func newMemoryListCmd() *cobra.Command {
	var detailed bool
	
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List memory content",
		Long:  `List memory content for agents and stacks`,
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
			
			// Walk memory directory
			fmt.Printf("Memory storage located at: %s\n\n", memoryPath)
			
			// Count files and directories
			var totalFiles, totalDirs int
			var totalSize int64
			
			err = filepath.Walk(memoryPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				
				// Skip root directory
				if path == memoryPath {
					return nil
				}
				
				// Get relative path
				relPath, err := filepath.Rel(memoryPath, path)
				if err != nil {
					return err
				}
				
				// Count files and directories
				if info.IsDir() {
					totalDirs++
					if detailed {
						fmt.Printf("[DIR] %s\n", relPath)
					}
				} else {
					totalFiles++
					totalSize += info.Size()
					if detailed {
						// Format file size
						sizeStr := formatSize(info.Size())
						// Format time
						timeStr := info.ModTime().Format("2006-01-02 15:04:05")
						fmt.Printf("[FILE] %s (%.30s, %s)\n", relPath, sizeStr, timeStr)
					}
				}
				
				return nil
			})
			
			if err != nil {
				return fmt.Errorf("failed to walk memory directory: %w", err)
			}
			
			// Print summary
			fmt.Printf("\nMemory storage summary:\n")
			fmt.Printf("  Directories: %d\n", totalDirs)
			fmt.Printf("  Files: %d\n", totalFiles)
			fmt.Printf("  Total size: %s\n", formatSize(totalSize))
			
			return nil
		},
	}
	
	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed information")
	
	return cmd
}
