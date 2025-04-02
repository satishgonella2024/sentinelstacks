package stack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/parser"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// NewInspectCommand creates a 'stack inspect' command
func NewInspectCommand() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "inspect [stackfile_path]",
		Short: "Inspect a stack",
		Long:  `Display detailed information about a specific stack`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the stack file path
			stackFilePath := args[0]

			// If no extension provided, try with .yaml and .yml
			if filepath.Ext(stackFilePath) == "" {
				yamlPath := stackFilePath + ".yaml"
				ymlPath := stackFilePath + ".yml"

				if _, err := os.Stat(yamlPath); err == nil {
					stackFilePath = yamlPath
				} else if _, err := os.Stat(ymlPath); err == nil {
					stackFilePath = ymlPath
				} else {
					// Try with Stackfile.yaml in the specified directory
					dirStackfile := filepath.Join(stackFilePath, "Stackfile.yaml")
					if _, err := os.Stat(dirStackfile); err == nil {
						stackFilePath = dirStackfile
					}
				}
			}

			// Read the stack file
			content, err := ioutil.ReadFile(stackFilePath)
			if err != nil {
				return fmt.Errorf("failed to read stack file: %w", err)
			}

			// Parse the stack
			var spec stack.StackSpec
			p := parser.NewStackParser()
			
			if filepath.Ext(stackFilePath) == ".json" {
				spec, err = p.ParseFromJSON(string(content))
			} else {
				spec, err = p.ParseFromYAML(string(content))
			}
			
			if err != nil {
				return fmt.Errorf("failed to parse stack file: %w", err)
			}

			// Display the stack information based on format
			switch format {
			case "json":
				// Output as JSON
				jsonData, err := json.MarshalIndent(spec, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to convert to JSON: %w", err)
				}
				fmt.Println(string(jsonData))
				
			case "yaml":
				// Output as YAML
				yamlData, err := yaml.Marshal(spec)
				if err != nil {
					return fmt.Errorf("failed to convert to YAML: %w", err)
				}
				fmt.Println(string(yamlData))
				
			case "dot":
				// Output as GraphViz DOT format for visualization
				dotContent := generateDotGraph(spec)
				fmt.Println(dotContent)
				
			default:
				// Default to human-readable format
				printHumanReadable(spec)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&format, "format", "f", "human", "Output format (human, json, yaml, dot)")

	return cmd
}

// printHumanReadable outputs stack information in a human-readable format
func printHumanReadable(spec stack.StackSpec) {
	fmt.Printf("Stack: %s (version %s)\n", spec.Name, spec.Version)
	fmt.Printf("Description: %s\n", spec.Description)
	fmt.Printf("\nAgents: %d\n", len(spec.Agents))
	
	// Build a dependency map for better visualization
	dependencyMap := make(map[string][]string)
	for _, agent := range spec.Agents {
		var dependencies []string
		dependencies = append(dependencies, agent.InputFrom...)
		dependencies = append(dependencies, agent.Depends...)
		dependencyMap[agent.ID] = dependencies
	}
	
	// Print agent details
	for i, agent := range spec.Agents {
		fmt.Printf("\n%d. %s\n", i+1, agent.ID)
		fmt.Printf("   - Uses: %s\n", agent.Uses)
		
		// Show dependencies
		deps := dependencyMap[agent.ID]
		if len(deps) > 0 {
			fmt.Printf("   - Dependencies: %v\n", deps)
		}
		
		// Show parameters
		if len(agent.Params) > 0 {
			fmt.Printf("   - Parameters:\n")
			for k, v := range agent.Params {
				fmt.Printf("     * %s: %v\n", k, v)
			}
		}
		
		// Show input/output info
		if agent.InputKey != "" {
			fmt.Printf("   - Input key: %s\n", agent.InputKey)
		}
		if agent.OutputKey != "" {
			fmt.Printf("   - Output key: %s\n", agent.OutputKey)
		}
	}
	
	// Print execution flow
	fmt.Printf("\nExecution Flow:\n")
	
	// Create a DAG and get topological sort
	dag, err := stack.NewDAG(spec)
	if err != nil {
		fmt.Printf("Error creating execution graph: %v\n", err)
		return
	}
	
	executionOrder, err := dag.TopologicalSort()
	if err != nil {
		fmt.Printf("Error determining execution order: %v\n", err)
		return
	}
	
	fmt.Printf("The agents will be executed in the following order:\n")
	for i, agentID := range executionOrder {
		fmt.Printf("%d. %s\n", i+1, agentID)
	}
}

// generateDotGraph creates a GraphViz DOT representation of the stack
func generateDotGraph(spec stack.StackSpec) string {
	dot := "digraph G {\n"
	dot += "  rankdir=LR;\n"
	dot += "  node [shape=box, style=filled, fillcolor=lightblue];\n\n"
	
	// Add nodes
	for _, agent := range spec.Agents {
		label := fmt.Sprintf("%s\\n(%s)", agent.ID, agent.Uses)
		dot += fmt.Sprintf("  \"%s\" [label=\"%s\"];\n", agent.ID, label)
	}
	
	dot += "\n"
	
	// Add edges
	for _, agent := range spec.Agents {
		// Add inputFrom edges
		for _, input := range agent.InputFrom {
			dot += fmt.Sprintf("  \"%s\" -> \"%s\";\n", input, agent.ID)
		}
		
		// Add depends edges
		for _, dep := range agent.Depends {
			dot += fmt.Sprintf("  \"%s\" -> \"%s\" [style=dashed];\n", dep, agent.ID)
		}
	}
	
	dot += "}\n"
	return dot
}
