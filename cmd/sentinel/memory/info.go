package memory

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"github.com/satishgonella2024/sentinelstacks/internal/memory"
)

// newMemoryInfoCmd creates a command to show memory information
func newMemoryInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show memory information",
		Long:  `Display information about memory subsystem configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get memory path
			memoryPath, err := memory.DefaultMemoryPath()
			if err != nil {
				return fmt.Errorf("failed to get memory path: %w", err)
			}
			
			// Check if memory path exists
			memoryExists := true
			if _, err := os.Stat(memoryPath); os.IsNotExist(err) {
				memoryExists = false
			}
			
			// Print information
			fmt.Println("Memory Subsystem Information:")
			fmt.Printf("  Base path: %s\n", memoryPath)
			fmt.Printf("  Initialized: %v\n", memoryExists)
			
			if memoryExists {
				// Get store types
				storeTypes, err := os.ReadDir(memoryPath)
				if err != nil {
					return fmt.Errorf("failed to read memory directory: %w", err)
				}
				
				fmt.Println("\nAvailable store types:")
				for _, storeType := range storeTypes {
					if storeType.IsDir() {
						fmt.Printf("  - %s\n", storeType.Name())
					}
				}
			} else {
				fmt.Println("\nMemory subsystem is not initialized")
				fmt.Println("Run 'sentinel memory init' to initialize")
			}
			
			return nil
		},
	}
	
	return cmd
}
