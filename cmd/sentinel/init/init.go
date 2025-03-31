package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Template for a basic Sentinelfile
const defaultSentinelfileTemplate = `# Sentinelfile for MyAgent

Create an agent that [describe your agent's purpose].

The agent should be able to:
- [Capability 1]
- [Capability 2]
- [Capability 3]

The agent should use claude-3.7-sonnet as its base model.

It should maintain state about [state description].

When the conversation starts, the agent should [initialization behavior].

Allow the agent to access the following tools:
- [Tool 1]
- [Tool 2]

Set [parameter name] to [value].
`

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var name string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Sentinelfile",
		Long:  `Create a new Sentinelfile in the current directory or specified directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine the filename and directory
			filename := "Sentinelfile"
			if name != "" {
				// Create a directory for the named agent if it doesn't exist
				if err := os.MkdirAll(name, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
				filename = filepath.Join(name, filename)
			}

			// Check if file already exists
			if _, err := os.Stat(filename); err == nil {
				return fmt.Errorf("file %s already exists", filename)
			}

			// Create the Sentinelfile
			if err := os.WriteFile(filename, []byte(defaultSentinelfileTemplate), 0644); err != nil {
				return fmt.Errorf("failed to write Sentinelfile: %w", err)
			}

			fmt.Printf("Created Sentinelfile at %s\n", filename)
			return nil
		},
	}

	initCmd.Flags().StringVar(&name, "name", "", "Name for the agent directory")

	return initCmd
}
