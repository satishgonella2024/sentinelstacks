package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config file paths
const (
	configDir  = ".sentinel"
	configFile = "config.yaml"
)

// NewConfigCmd creates the config command and its subcommands
func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage SentinelStacks configuration",
		Long:  `View and modify the configuration settings for SentinelStacks`,
	}

	// Add subcommands
	configCmd.AddCommand(newConfigSetCmd())
	configCmd.AddCommand(newConfigGetCmd())
	configCmd.AddCommand(newConfigListCmd())

	return configCmd
}

// newConfigSetCmd creates the 'config set' command
func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Long:  `Set a configuration value in the SentinelStacks configuration file`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// Initialize configuration
			if err := initConfig(); err != nil {
				return err
			}

			// Set the configuration value
			viper.Set(key, value)

			// Save the configuration
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("error writing config: %w", err)
			}

			fmt.Printf("Set %s to %s\n", key, value)
			return nil
		},
	}
}

// newConfigGetCmd creates the 'config get' command
func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Long:  `Get a configuration value from the SentinelStacks configuration file`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			// Initialize configuration
			if err := initConfig(); err != nil {
				return err
			}

			// Get the configuration value
			value := viper.Get(key)
			if value == nil {
				return fmt.Errorf("key %s not found in configuration", key)
			}

			fmt.Printf("%v\n", value)
			return nil
		},
	}
}

// newConfigListCmd creates the 'config list' command
func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		Long:  `List all configuration values in the SentinelStacks configuration file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize configuration
			if err := initConfig(); err != nil {
				return err
			}

			// List all settings
			settings := viper.AllSettings()
			if len(settings) == 0 {
				fmt.Println("No configuration settings found")
				return nil
			}

			fmt.Println("Configuration settings:")
			printSettings(settings, "")
			return nil
		},
	}
}

// printSettings prints a nested map of settings with proper indentation
func printSettings(settings map[string]interface{}, prefix string) {
	for k, v := range settings {
		key := prefix + k
		if nested, ok := v.(map[string]interface{}); ok {
			printSettings(nested, key+".")
		} else {
			fmt.Printf("  %s: %v\n", key, v)
		}
	}
}

// initConfig initializes the configuration
func initConfig() error {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	// Set configuration file path
	configPath := filepath.Join(home, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Set config name and path
	viper.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))
	viper.SetConfigType(strings.TrimPrefix(filepath.Ext(configFile), "."))
	viper.AddConfigPath(configPath)

	// Set default values
	viper.SetDefault("llm.provider", "claude")
	viper.SetDefault("registry.url", "https://registry.sentinelstacks.com")

	// Try to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// If the config file doesn't exist, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFilePath := filepath.Join(configPath, configFile)
			if err := viper.WriteConfigAs(configFilePath); err != nil {
				return fmt.Errorf("error creating config file: %w", err)
			}
		} else {
			return fmt.Errorf("error reading config: %w", err)
		}
	}

	return nil
}
