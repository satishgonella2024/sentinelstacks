package stack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// NewInitCommand creates a new command for initializing stacks
func NewInitCommand() *cobra.Command {
	var (
		templateName string
		description  string
		outputFile   string
	)

	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new stack",
		Long:  `Create a new stack definition file with specified template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			
			// Validate stack name (must be valid filename)
			if strings.ContainsAny(stackName, "/\\:*?\"<>|") {
				return fmt.Errorf("invalid stack name: must not contain /\\:*?\"<>|")
			}
			
			// Set default output file if not specified
			if outputFile == "" {
				outputFile = "Stackfile.yaml"
			}
			
			// Check if file already exists
			if _, err := os.Stat(outputFile); err == nil {
				return fmt.Errorf("file already exists: %s", outputFile)
			}
			
			// Generate template based on name
			var template string
			switch templateName {
			case "basic":
				template = getBasicTemplate(stackName, description)
			case "analyzer":
				template = getAnalyzerTemplate(stackName, description)
			case "pipeline":
				template = getPipelineTemplate(stackName, description)
			default:
				return fmt.Errorf("unknown template: %s", templateName)
			}
			
			// Create directory if needed
			dir := filepath.Dir(outputFile)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			}
			
			// Write file
			if err := ioutil.WriteFile(outputFile, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			
			fmt.Printf("Stack initialized: %s\n", stackName)
			fmt.Printf("Stack definition written to: %s\n", outputFile)
			
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&templateName, "template", "t", "basic", "Stack template (basic, analyzer, pipeline)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Stack description")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: Stackfile.yaml)")

	return cmd
}
