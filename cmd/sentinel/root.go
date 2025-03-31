package sentinel

import (
	"github.com/spf13/cobra"

	"github.com/sentinelstacks/sentinel/cmd/sentinel/build"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/config"
	initCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/init"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/run"
)

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "SentinelStacks - AI Agent Management System",
	Long: `SentinelStacks is a comprehensive system for creating, managing,
and distributing AI agents using natural language definitions.

It provides a Docker-like workflow for AI agents:
- Define agents using Sentinelfiles
- Build agent images from Sentinelfiles
- Run agents from images locally or from registries
- Share agents through registries`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Add commands
	rootCmd.AddCommand(initCmd.NewInitCmd())  // Init command
	rootCmd.AddCommand(build.NewBuildCmd())   // Build command
	rootCmd.AddCommand(run.NewRunCmd())       // Run command
	rootCmd.AddCommand(config.NewConfigCmd()) // Config command

	// TODO: Add more commands as they are implemented
	// rootCmd.AddCommand(NewPushCmd())
	// rootCmd.AddCommand(NewPullCmd())
	// rootCmd.AddCommand(NewPsCmd())
	// rootCmd.AddCommand(NewStopCmd())
	// rootCmd.AddCommand(NewLogsCmd())
}
