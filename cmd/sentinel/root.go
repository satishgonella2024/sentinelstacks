package sentinel

import (
	"github.com/spf13/cobra"

	apiCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/api"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/build"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/config"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/images"
	initCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/init"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/logs"
	multimodalCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/multimodal"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/ps"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/run"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/stop"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/version"
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
	rootCmd.AddCommand(initCmd.NewInitCmd())             // Init command
	rootCmd.AddCommand(build.NewBuildCmd())              // Build command
	rootCmd.AddCommand(run.NewRunCmd())                  // Run command
	rootCmd.AddCommand(ps.NewPsCmd())                    // PS command
	rootCmd.AddCommand(stop.NewStopCmd())                // Stop command
	rootCmd.AddCommand(logs.NewLogsCmd())                // Logs command
	rootCmd.AddCommand(images.NewImagesCmd())            // Images command
	rootCmd.AddCommand(config.NewConfigCmd())            // Config command
	rootCmd.AddCommand(version.NewVersionCmd())          // Version command
	rootCmd.AddCommand(apiCmd.NewAPICmd())               // API command
	rootCmd.AddCommand(multimodalCmd.NewMultimodalCmd()) // Multimodal command

	// TODO: Add more commands as they are implemented
	// rootCmd.AddCommand(NewPushCmd())
	// rootCmd.AddCommand(NewPullCmd())
}
