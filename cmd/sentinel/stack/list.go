package stack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

// StackInfo holds meta information about a stack
type StackInfo struct {
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Version     string    `json:"version" yaml:"version"`
	FilePath    string    `json:"filePath" yaml:"filePath"`
	AgentCount  int       `json:"agentCount" yaml:"agentCount"`
	Modified    time.Time `json:"modified" yaml:"modified"`
}

// NewListCommand creates a 'stack list' command
func NewListCommand() *cobra.Command {
	var (
		format       string
		listAll      bool
		nameFilter   string
		sortBy       string
		stacksFolder string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available stacks",
		Long:  `Display a list of available stacks in the local registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set directories to search
			searchDirs := []string{"."}
			
			// Add home directory .sentinel/stacks if listAll
			if listAll || stacksFolder == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				
				// Get stacks dir from config or use default
				if stacksFolder == "" {
					if viper.IsSet("stacks.dir") {
						stacksFolder = viper.GetString("stacks.dir")
					} else {
						stacksFolder = filepath.Join(home, ".sentinel", "stacks")
					}
				}
				
				// Create directory if it doesn't exist
				if _, err := os.Stat(stacksFolder); os.IsNotExist(err) {
					if err := os.MkdirAll(stacksFolder, 0755); err != nil {
						return fmt.Errorf("failed to create stacks directory: %w", err)
					}
				}
				
				searchDirs = append(searchDirs, stacksFolder)
			}
			
			// Find and parse stack files
			stackInfos, err := findStackFiles(searchDirs, nameFilter)
			if err != nil {
				return fmt.Errorf("failed to find stack files: %w", err)
			}
			
			// Sort by selected field
			switch sortBy {
			case "name":
				sort.Slice(stackInfos, func(i, j int) bool {
					return stackInfos[i].Name < stackInfos[j].Name
				})
			case "modified":
				sort.Slice(stackInfos, func(i, j int) bool {
					return stackInfos[i].Modified.After(stackInfos[j].Modified)
				})
			case "agents":
				sort.Slice(stackInfos, func(i, j int) bool {
					return stackInfos[i].AgentCount > stackInfos[j].AgentCount
				})
			}
			
			// Display based on format
			switch format {
			case "wide":
				printWideFormat(stackInfos)
			case "json":
				printJSONFormat(stackInfos)
			case "yaml":
				printYAMLFormat(stackInfos)
			default:
				printDefaultFormat(stackInfos)
			}
			
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&format, "format", "f", "default", "Output format (default, wide, json, yaml)")
	cmd.Flags().BoolVarP(&listAll, "all", "a", false, "List all stacks (including those in home directory)")
	cmd.Flags().StringVarP(&nameFilter, "filter", "n", "", "Filter stacks by name")
	cmd.Flags().StringVarP(&sortBy, "sort", "s", "name", "Sort by (name, modified, agents)")
	cmd.Flags().StringVarP(&stacksFolder, "dir", "d", "", "Directory to search for stacks")

	return cmd
}

// findStackFiles searches directories for stack files
func findStackFiles(dirs []string, nameFilter string) ([]StackInfo, error) {
	var stackInfos []StackInfo
	
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip files we can't access
			}
			
			// Check if this is a stack file
			if !info.IsDir() && isStackFile(path) {
				// Parse the stack file
				stackInfo, err := parseStackFile(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to parse %s: %v\n", path, err)
					return nil
				}
				
				// Apply name filter if any
				if nameFilter != "" && !strings.Contains(strings.ToLower(stackInfo.Name), strings.ToLower(nameFilter)) {
					return nil
				}
				
				// Set modified time
				stackInfo.Modified = info.ModTime()
				
				// Add to results
				stackInfos = append(stackInfos, stackInfo)
			}
			
			return nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
		}
	}
	
	return stackInfos, nil
}

// isStackFile checks if a file is a stack definition
func isStackFile(path string) bool {
	ext := filepath.Ext(path)
	baseName := filepath.Base(path)
	
	return (ext == ".yaml" || ext == ".yml") &&
		(baseName == "Stackfile.yaml" || baseName == "Stackfile.yml" ||
		 strings.HasPrefix(baseName, "stack_") ||
		 strings.HasPrefix(baseName, "Stack_"))
}

// parseStackFile reads and parses a stack file
func parseStackFile(path string) (StackInfo, error) {
	// Read the file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return StackInfo{}, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Parse YAML
	var spec stack.StackSpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		return StackInfo{}, fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Create stack info
	return StackInfo{
		Name:        spec.Name,
		Description: spec.Description,
		Version:     spec.Version,
		FilePath:    path,
		AgentCount:  len(spec.Agents),
	}, nil
}

// printDefaultFormat prints stacks in a simple tabular format
func printDefaultFormat(stacks []StackInfo) {
	if len(stacks) == 0 {
		fmt.Println("No stacks found")
		return
	}
	
	// Print header
	fmt.Printf("%-20s %-10s %-6s\n", "NAME", "VERSION", "AGENTS")
	fmt.Println(strings.Repeat("-", 40))
	
	// Print rows
	for _, stack := range stacks {
		fmt.Printf("%-20s %-10s %-6d\n", 
			truncateString(stack.Name, 20),
			truncateString(stack.Version, 10),
			stack.AgentCount)
	}
	
	fmt.Printf("\nFound %d stacks\n", len(stacks))
}

// printWideFormat prints stacks in a wide tabular format
func printWideFormat(stacks []StackInfo) {
	if len(stacks) == 0 {
		fmt.Println("No stacks found")
		return
	}
	
	// Print header
	fmt.Printf("%-20s %-10s %-6s %-40s %-19s %-30s\n", 
		"NAME", "VERSION", "AGENTS", "DESCRIPTION", "MODIFIED", "PATH")
	fmt.Println(strings.Repeat("-", 120))
	
	// Print rows
	for _, stack := range stacks {
		fmt.Printf("%-20s %-10s %-6d %-40s %-19s %-30s\n", 
			truncateString(stack.Name, 20),
			truncateString(stack.Version, 10),
			stack.AgentCount,
			truncateString(stack.Description, 40),
			stack.Modified.Format("2006-01-02 15:04:05"),
			truncateString(stack.FilePath, 30))
	}
	
	fmt.Printf("\nFound %d stacks\n", len(stacks))
}

// printJSONFormat prints stacks in JSON format
func printJSONFormat(stacks []StackInfo) {
	// This would use json.Marshal in a real implementation
	fmt.Println("{")
	for i, stack := range stacks {
		fmt.Printf("  \"%s\": {\n", stack.Name)
		fmt.Printf("    \"name\": \"%s\",\n", stack.Name)
		fmt.Printf("    \"version\": \"%s\",\n", stack.Version)
		fmt.Printf("    \"agents\": %d,\n", stack.AgentCount)
		fmt.Printf("    \"description\": \"%s\",\n", stack.Description)
		fmt.Printf("    \"modified\": \"%s\",\n", stack.Modified.Format(time.RFC3339))
		fmt.Printf("    \"path\": \"%s\"\n", stack.FilePath)
		if i < len(stacks)-1 {
			fmt.Println("  },")
		} else {
			fmt.Println("  }")
		}
	}
	fmt.Println("}")
}

// printYAMLFormat prints stacks in YAML format
func printYAMLFormat(stacks []StackInfo) {
	// This would use yaml.Marshal in a real implementation
	for _, stack := range stacks {
		fmt.Printf("- name: %s\n", stack.Name)
		fmt.Printf("  version: %s\n", stack.Version)
		fmt.Printf("  agents: %d\n", stack.AgentCount)
		fmt.Printf("  description: %s\n", stack.Description)
		fmt.Printf("  modified: %s\n", stack.Modified.Format(time.RFC3339))
		fmt.Printf("  path: %s\n", stack.FilePath)
		fmt.Println()
	}
}

// truncateString truncates a string to a specific length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
