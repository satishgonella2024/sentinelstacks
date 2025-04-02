package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	// Import our local network package
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/network-module/pkg/network"
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "sentinel",
		Short: "SentinelStacks - AI Agent Management System",
		Long: `SentinelStacks is a comprehensive system for creating, managing,
and distributing AI agents using natural language definitions.`,
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Set data directory
	dataDir := filepath.Join(home, ".sentinel", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Add network command to root command
	rootCmd.AddCommand(network.NewNetworkCmd(dataDir))

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
