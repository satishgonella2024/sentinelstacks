package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "0.1.0"
	GitCommit = "development"
	BuildDate = "unknown"
)

// VersionCmd returns the version command
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SentinelStacks v%s\n", Version)
			fmt.Printf("Git commit: %s\n", GitCommit)
			fmt.Printf("Built on: %s\n", BuildDate)
		},
	}

	return cmd
}
