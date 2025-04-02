package stack

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/client"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/format"
	packages "github.com/satishgonella2024/sentinelstacks/internal/registry/package"
	"github.com/satishgonella2024/sentinelstacks/internal/registry/security"
)

// NewPullCommand creates a 'stack pull' command
func NewPullCommand() *cobra.Command {
	var (
		outputDir    string
		verify       bool
		savePackage  bool
		extractAgent bool
	)

	cmd := &cobra.Command{
		Use:   "pull [name:version]",
		Short: "Pull a stack from the registry",
		Long:  `Download a stack from the registry for local use`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse the reference
			reference := args[0]
			name, version, err := parseReference(reference)
			if err != nil {
				return err
			}

			// Set default version if not provided
			if version == "" {
				version = "latest"
				fmt.Printf("No version specified, using latest\n")
			}

			// Determine output directory
			if outputDir == "" {
				outputDir = "."
			}

			// Ensure output directory exists
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Get registry URL from config
			registryURL := viper.GetString("registry.url")
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Get auth token from config
			authToken := viper.GetString("registry.auth_token")

			// Create registry client
			registryClient := client.NewRegistryClient(registryURL, authToken)

			fmt.Printf("Pulling stack '%s' (version %s) from registry...\n", name, version)

			// Pull the stack
			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			defer cancel()

			// Create a temporary directory for the package
			tempDir, err := os.MkdirTemp("", "sentinel-pull-")
			if err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}
			defer func() {
				if !savePackage {
					os.RemoveAll(tempDir)
				}
			}()

			// Determine package path
			packagePath := filepath.Join(tempDir, format.GetDefaultFilename(name, version, "stack"))

			// Pull the package
			if err := registryClient.PullPackage(ctx, name, version, packagePath); err != nil {
				return fmt.Errorf("failed to pull package: %w", err)
			}

			// Verify the package if requested
			if verify {
				fmt.Printf("Verifying package integrity and signatures...\n")
				
				// Get keys directory from config
				keysDir := viper.GetString("security.keys_dir")
				if keysDir == "" {
					home, _ := os.UserHomeDir()
					keysDir = filepath.Join(home, ".sentinel", "keys")
				}

				// Create key manager
				keyManager, err := security.NewKeyManager(keysDir)
				if err != nil {
					return fmt.Errorf("failed to create key manager: %w", err)
				}

				// Verify package
				valid, failures, err := registryClient.VerifyPackage(packagePath)
				if err != nil {
					return fmt.Errorf("package verification failed: %w", err)
				}

				if !valid {
					return fmt.Errorf("package integrity check failed: %v", failures)
				}

				fmt.Printf("Package verification successful\n")
			}

			// Extract the package
			extractDir := filepath.Join(tempDir, "extract")
			if err := os.MkdirAll(extractDir, 0755); err != nil {
				return fmt.Errorf("failed to create extraction directory: %w", err)
			}

			fmt.Printf("Extracting package contents...\n")
			pkg := &packages.SentinelPackage{}
			if err := pkg.Unpackage(packagePath, extractDir); err != nil {
				return fmt.Errorf("failed to extract package: %w", err)
			}

			// Find the stack definition file
			var stackFilePath string
			for _, file := range pkg.Manifest.Files {
				if file.IsMain && (strings.HasSuffix(file.Path, format.StackDefinitionExtension) || 
				                    strings.HasSuffix(file.Path, ".yaml") ||
									strings.HasSuffix(file.Path, ".yml")) {
					stackFilePath = filepath.Join(extractDir, file.Path)
					break
				}
			}

			if stackFilePath == "" {
				return fmt.Errorf("no stack definition found in package")
			}

			// Copy stack file to output directory
			outputStackFile := filepath.Join(outputDir, format.GetDefaultFilename(name, "", "stack-def"))
			if err := copyFile(stackFilePath, outputStackFile); err != nil {
				return fmt.Errorf("failed to copy stack definition: %w", err)
			}

			fmt.Printf("Stack definition saved to: %s\n", outputStackFile)

			// Copy README if present
			for _, file := range pkg.Manifest.Files {
				if strings.ToLower(filepath.Base(file.Path)) == "readme.md" {
					readmePath := filepath.Join(extractDir, file.Path)
					outputReadme := filepath.Join(outputDir, "README.md")
					if err := copyFile(readmePath, outputReadme); err != nil {
						fmt.Printf("Warning: Failed to copy README: %v\n", err)
					} else {
						fmt.Printf("Copied README to: %s\n", outputReadme)
					}
					break
				}
			}

			// Extract required agents if requested
			if extractAgent && len(pkg.Manifest.Dependencies) > 0 {
				fmt.Printf("Pulling required agents...\n")
				
				for _, dep := range pkg.Manifest.Dependencies {
					if dep.Type == packages.PackageTypeAgent {
						fmt.Printf("  - Pulling agent: %s:%s\n", dep.Name, dep.Version)
						
						// Pull agent
						// In a real implementation, this would call the sentinel CLI
						agentRef := fmt.Sprintf("%s:%s", dep.Name, dep.Version)
						cmd := exec.Command("sentinel", "pull", agentRef)
						output, err := cmd.CombinedOutput()
						if err != nil {
							fmt.Printf("    Warning: Failed to pull agent: %v\n", err)
							fmt.Printf("    Output: %s\n", string(output))
							
							if dep.Required {
								return fmt.Errorf("failed to pull required agent %s: %w", agentRef, err)
							}
						} else {
							fmt.Printf("    Successfully pulled agent: %s\n", agentRef)
						}
					}
				}
			}

			// Save package if requested
			if savePackage {
				outputPackagePath := filepath.Join(outputDir, format.GetDefaultFilename(name, version, "stack"))
				if err := copyFile(packagePath, outputPackagePath); err != nil {
					return fmt.Errorf("failed to save package: %w", err)
				}
				fmt.Printf("Package saved to: %s\n", outputPackagePath)
			}

			fmt.Printf("Successfully pulled stack '%s' (version %s)\n", name, version)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	cmd.Flags().BoolVarP(&verify, "verify", "v", true, "Verify package integrity and signatures")
	cmd.Flags().BoolVarP(&savePackage, "save-package", "s", false, "Save the package file")
	cmd.Flags().BoolVarP(&extractAgent, "extract-agents", "e", false, "Extract required agents")

	return cmd
}

// parseReference parses a stack reference (name:version)
func parseReference(reference string) (string, string, error) {
	parts := strings.SplitN(reference, ":", 2)
	
	name := parts[0]
	if name == "" {
		return "", "", fmt.Errorf("stack name cannot be empty")
	}
	
	version := ""
	if len(parts) > 1 {
		version = parts[1]
	}
	
	return name, version, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Read source file
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	
	// Write to destination file
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}
	
	return nil
}
