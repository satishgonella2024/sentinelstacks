package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "0.1.0"
	BuildDate = "2024-03-31"
	GitCommit = "development"
)

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print detailed version information about SentinelStacks.`,
		Run: func(cmd *cobra.Command, args []string) {
			short, _ := cmd.Flags().GetBool("short")
			runVersion(short)
		},
	}

	// Add flags
	cmd.Flags().BoolP("short", "s", false, "Print just the version number")

	return cmd
}

// runVersion executes the version command
func runVersion(short bool) {
	if short {
		fmt.Println(Version)
		return
	}

	fmt.Println("SentinelStacks AI Agent Management System")
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Git commit: %s\n", GitCommit)
	fmt.Printf("Built:      %s\n", BuildDate)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
