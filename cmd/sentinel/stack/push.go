package stack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/client"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/format"
	packages "github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/package"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/security"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// NewPushCommand creates a 'stack push' command
func NewPushCommand() *cobra.Command {
	var (
		author        string
		buildPackage  bool
		sign          bool
		keyID         string
		outputPackage string
	)

	cmd := &cobra.Command{
		Use:   "push [stackfile_path]",
		Short: "Push a stack to the registry",
		Long:  `Package and push a stack to the registry for sharing with others`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the stack file path
			stackFilePath := args[0]

			// Check if file exists
			if _, err := os.Stat(stackFilePath); os.IsNotExist(err) {
				return fmt.Errorf("stack file not found: %s", stackFilePath)
			}

			// Set default author if not provided
			if author == "" {
				// Try to get from git config
				author = getGitAuthor()
				if author == "" {
					author = "unknown"
				}
			}

			// Parse the stack file
			stackSpec, err := parseStackFile(stackFilePath)
			if err != nil {
				return fmt.Errorf("failed to parse stack file: %w", err)
			}

			fmt.Printf("Pushing stack '%s' (version %s) to registry...\n", stackSpec.Name, stackSpec.Version)

			// Get registry URL from config
			registryURL := viper.GetString("registry.url")
			if registryURL == "" {
				registryURL = "https://registry.sentinelstacks.io"
			}

			// Get auth token from config
			authToken := viper.GetString("registry.auth_token")

			// Create registry client
			registryClient := client.NewRegistryClient(registryURL, authToken)

			// Check if we need to build and sign the package
			if buildPackage {
				// Create package path
				packagePath := outputPackage
				if packagePath == "" {
					// Default to stack name in current directory
					packagePath = format.GetDefaultFilename(stackSpec.Name, stackSpec.Version, "stack")
				}

				// Create package builder
				builder := packages.NewPackageBuilder(packages.PackageTypeStack, stackSpec.Name, stackSpec.Version, stackSpec.Description, author)

				// Add main stack file
				if err := builder.AddFile(stackFilePath, format.GetDefaultFilename(stackSpec.Name, "", "stack-def"), true, packages.FileTypeManifest); err != nil {
					return fmt.Errorf("failed to add stack file: %w", err)
				}

				// Add any additional files in the same directory
				stackDir := filepath.Dir(stackFilePath)
				readmePath := filepath.Join(stackDir, "README.md")
				if _, err := os.Stat(readmePath); err == nil {
					if err := builder.AddFile(readmePath, "README.md", false, packages.FileTypeDoc); err != nil {
						return fmt.Errorf("failed to add README: %w", err)
					}
				}

				// Add examples directory if it exists
				examplesDir := filepath.Join(stackDir, "examples")
				if _, err := os.Stat(examplesDir); err == nil {
					if err := builder.AddDirectory(examplesDir, "examples"); err != nil {
						return fmt.Errorf("failed to add examples directory: %w", err)
					}
				}

				// Sign the package if requested
				if sign {
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

					// Set key manager for signing
					builder.SetKeyManager(keyManager)

					// Use default key ID if not specified
					if keyID == "" {
						keyID = "default"
					}

					// Check if key exists, generate if not
					keyPath := filepath.Join(keysDir, keyID+".key")
					if _, err := os.Stat(keyPath); os.IsNotExist(err) {
						fmt.Printf("Generating new signing key: %s\n", keyID)
						if err := keyManager.GenerateKeyPair(keyID, 2048); err != nil {
							return fmt.Errorf("failed to generate key pair: %w", err)
						}
					}
				}

				// Build the package
				fmt.Printf("Building package: %s\n", packagePath)
				if err := builder.Build(packagePath); err != nil {
					return fmt.Errorf("failed to build package: %w", err)
				}

				// Push the package
				fmt.Printf("Pushing package to registry: %s\n", registryURL)
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				if err := registryClient.PushPackage(ctx, packagePath); err != nil {
					return fmt.Errorf("failed to push package: %w", err)
				}
			} else {
				// Push the stack directly
				fmt.Printf("Pushing stack directly to registry: %s\n", registryURL)
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				if err := registryClient.PushStack(ctx, stackSpec, author); err != nil {
					return fmt.Errorf("failed to push stack: %w", err)
				}
			}

			fmt.Printf("Successfully pushed stack '%s' (version %s) to registry\n", stackSpec.Name, stackSpec.Version)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&author, "author", "a", "", "Author of the stack")
	cmd.Flags().BoolVarP(&buildPackage, "build", "b", false, "Build a package file")
	cmd.Flags().BoolVarP(&sign, "sign", "s", false, "Sign the package")
	cmd.Flags().StringVarP(&keyID, "key", "k", "", "Key ID for signing")
	cmd.Flags().StringVarP(&outputPackage, "output", "o", "", "Output package path")

	return cmd
}

// parseStackFile parses a stack file (YAML or JSON)
func parseStackFile(filePath string) (stack.StackSpec, error) {
	var stackSpec stack.StackSpec

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return stackSpec, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".json" {
		if err := json.Unmarshal(content, &stackSpec); err != nil {
			return stackSpec, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		// Assume YAML
		if err := yaml.Unmarshal(content, &stackSpec); err != nil {
			return stackSpec, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	// Set default version if not specified
	if stackSpec.Version == "" {
		stackSpec.Version = "1.0.0"
	}

	return stackSpec, nil
}

// getGitAuthor attempts to get author info from git config
func getGitAuthor() string {
	cmd := exec.Command("git", "config", "user.name")
	name, err := cmd.Output()
	if err != nil {
		return ""
	}

	cmd = exec.Command("git", "config", "user.email")
	email, err := cmd.Output()
	if err != nil {
		return strings.TrimSpace(string(name))
	}

	return fmt.Sprintf("%s <%s>", strings.TrimSpace(string(name)), strings.TrimSpace(string(email)))
}
