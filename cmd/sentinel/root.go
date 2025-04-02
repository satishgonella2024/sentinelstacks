package sentinel

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/api"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/build"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/chat"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/compose"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/config"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/exec"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/history"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/images"
	initCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/init"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/login"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/logout"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/logs"
	multimodalCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/multimodal"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/network"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/ps"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/pull"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/push"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/run"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/search"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/shell"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/stop"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/system"
	toolsCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/tools"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/version"
	"github.com/sentinelstacks/sentinel/cmd/sentinel/volume"
)

// rootCmd is the root command for the sentinel CLI
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

// ServiceProviderKey is the key for the service provider in the context
type contextKey int
const serviceProviderKey contextKey = 0

// Config file paths
const (
	configDir  = ".sentinel"
	configFile = "config.yaml"
)

// initConfig initializes the configuration
func initConfig() {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Set configuration file path
	configPath := filepath.Join(home, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Set config name and path
	viper.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))
	viper.SetConfigType(strings.TrimPrefix(filepath.Ext(configFile), "."))
	viper.AddConfigPath(configPath)

	// Add project directory as fallback
	viper.AddConfigPath("./")
	
	// Set environment variable prefix
	viper.SetEnvPrefix("SENTINEL")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("llm.provider", "claude")
	viper.SetDefault("llm.model", "claude-3.5-sonnet-20240627")
	viper.SetDefault("ollama.endpoint", "http://localhost:11434/api/generate")
	viper.SetDefault("registry.enable_cache", true)
	viper.SetDefault("home", filepath.Join(home, configDir))

	// Try to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Creating default config file at: %s/%s\n", configPath, configFile)
			configFilePath := filepath.Join(configPath, configFile)
			if err := viper.WriteConfigAs(configFilePath); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}
}

// Execute runs the root command
func Execute() error {
	// Initialize configuration
	initConfig()
	
	// Create a base context
	ctx := context.Background()
	
	// Initialize service provider
	dataDir := filepath.Join(viper.GetString("home"), "data")
	serviceProvider := NewBasicServiceProvider(dataDir)
	
	// Add service provider to context
	ctx = context.WithValue(ctx, serviceProviderKey, serviceProvider)
	
	// Execute with context
	return rootCmd.ExecuteContext(ctx)
}

// GetServiceProvider retrieves the service provider from the context
func GetServiceProvider(ctx context.Context) ServiceProvider {
	if ctx == nil {
		return nil
	}
	
	if sp, ok := ctx.Value(serviceProviderKey).(ServiceProvider); ok {
		return sp
	}
	
	// Create a new service provider if none exists in the context
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".sentinel", "data")
	return NewBasicServiceProvider(dataDir)
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
	rootCmd.AddCommand(chat.NewChatCmd())                // Chat command
	rootCmd.AddCommand(toolsCmd.NewToolsCmd())           // Tools command
	
	// Add new commands
	rootCmd.AddCommand(exec.NewExecCmd())                // Exec command
	rootCmd.AddCommand(shell.NewShellCmd())              // Shell command
	rootCmd.AddCommand(pull.NewPullCmd())                // Pull command
	rootCmd.AddCommand(push.NewPushCmd())                // Push command
	rootCmd.AddCommand(login.NewLoginCmd())              // Login command
	rootCmd.AddCommand(logout.NewLogoutCmd())            // Logout command
	rootCmd.AddCommand(search.NewSearchCmd())            // Search command
	rootCmd.AddCommand(history.NewHistoryCmd())          // History command
	rootCmd.AddCommand(network.NewNetworkCmd())          // Network command
	rootCmd.AddCommand(volume.NewVolumeCmd())            // Volume command
	rootCmd.AddCommand(compose.NewComposeCmd())          // Compose command
	rootCmd.AddCommand(system.NewSystemCmd())            // System command
}
