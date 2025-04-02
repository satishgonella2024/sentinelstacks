package stack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewStackCommand creates a new 'stack' command that groups related subcommands
func NewStackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Manage agent stacks",
		Long:  `Create and run multi-agent stacks to accomplish complex tasks`,
	}

	// Add subcommands
	cmd.AddCommand(NewRunCommand())    // Run stacks
	cmd.AddCommand(NewInitCommand())   // Initialize new stacks
	cmd.AddCommand(NewListCommand())   // List available stacks
	cmd.AddCommand(NewInspectCommand())// Inspect stack details
	cmd.AddCommand(NewPushCommand())   // Push stack to registry
	cmd.AddCommand(NewPullCommand())   // Pull stack from registry
	cmd.AddCommand(NewSearchCommand()) // Search stacks in registry

	return cmd
}

// getBasicTemplate returns a simple stack template with a single agent
func getBasicTemplate(name, description string) string {
	if description == "" {
		description = "A simple stack with a single agent"
	}
	
	return fmt.Sprintf(`name: %s
description: %s
version: 1.0.0
agents:
  - id: main-agent
    uses: generic-agent
    params:
      prompt: "Process the input and provide insights"
`, name, description)
}

// getAnalyzerTemplate returns a template for a data analysis stack
func getAnalyzerTemplate(name, description string) string {
	if description == "" {
		description = "A stack for analyzing data with multiple specialized agents"
	}
	
	return fmt.Sprintf(`name: %s
description: %s
version: 1.0.0
agents:
  - id: data-processor
    uses: processor
    params:
      format: "json"
  - id: analyzer
    uses: analyzer
    inputFrom:
      - data-processor
    params:
      analysis_type: "comprehensive"
  - id: summarizer
    uses: summarizer
    inputFrom:
      - analyzer
    params:
      format: "bullet_points"
`, name, description)
}

// getPipelineTemplate returns a template for a processing pipeline
func getPipelineTemplate(name, description string) string {
	if description == "" {
		description = "A data processing pipeline with multiple stages"
	}
	
	return fmt.Sprintf(`name: %s
description: %s
version: 1.0.0
agents:
  - id: extractor
    uses: extractor
    params:
      source: "database"
  - id: transformer
    uses: transformer
    inputFrom:
      - extractor
    params:
      operations: ["filter", "normalize"]
  - id: enricher
    uses: enricher
    inputFrom:
      - transformer
    params:
      enrichment_sources: ["knowledge_base", "external_api"]
  - id: loader
    uses: loader
    inputFrom:
      - enricher
    params:
      destination: "warehouse"
      format: "parquet"
`, name, description)
}
