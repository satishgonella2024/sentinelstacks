package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/satishgonella2024/sentinelstacks/internal/memory"
	"github.com/satishgonella2024/sentinelstacks/internal/parser"
	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

var (
	inputFile  string
	inputJson  string
	inputYaml  string
	inputNl    string
	verbose    bool
	timeoutSec int
)

// NewRunCommand creates a new command for running stacks
func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a stack of agents",
		Long:  `Execute a multi-agent stack defined in a Stackfile or natural language description`,
		RunE:  runStack,
	}

	// Add flags
	cmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to Stackfile (YAML or JSON)")
	cmd.Flags().StringVarP(&inputJson, "json", "j", "", "Stack definition as JSON string")
	cmd.Flags().StringVarP(&inputYaml, "yaml", "y", "", "Stack definition as YAML string")
	cmd.Flags().StringVarP(&inputNl, "nl", "n", "", "Natural language description of the stack")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().IntVarP(&timeoutSec, "timeout", "t", 0, "Execution timeout in seconds (0 for no timeout)")

	return cmd
}

// runStack is the main function for executing a stack
// runStack is the main function for executing a stack
func runStack(cmd *cobra.Command, args []string) error {
	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt, cancelling execution...")
		cancel()
	}()

	// Configure logging
	if verbose {
		log.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	// Parse input and create stack spec
	stackSpec, err := parseInput()
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	// Parse input JSON if provided
	var inputData map[string]interface{}
	if inputJson != "" {
		if err := json.Unmarshal([]byte(inputJson), &inputData); err != nil {
			return fmt.Errorf("failed to parse input JSON: %w", err)
		}
	}

	// Create memory factory
	memoryFactory, err := memory.NewMemoryStoreFactory("")
	if err != nil {
		return fmt.Errorf("failed to create memory factory: %w", err)
	}

	// Create and execute stack engine
	if verbose {
		fmt.Printf("Creating stack engine for '%s'...\n", stackSpec.Name)
	}

	// Create options
	engineOptions := []stack.EngineOption{
		stack.WithVerbose(verbose),
		stack.WithMemoryFactory(memoryFactory),
	}

	// Create engine
	engine, err := stack.NewStackEngine(stackSpec, engineOptions...)
	if err != nil {
		return fmt.Errorf("failed to create stack engine: %w", err)
	}

	// Display execution plan
	dag, err := engine.BuildExecutionGraph()
	if err != nil {
		return fmt.Errorf("failed to build execution graph: %w", err)
	}

	executionOrder, err := dag.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	fmt.Printf("Executing stack: %s\n", stackSpec.Name)
	fmt.Printf("Agents: %d\n", len(stackSpec.Agents))
	fmt.Printf("Execution order: %v\n", executionOrder)

	// Create execution options
	executeOptions := []stack.ExecuteOption{
		stack.WithTimeout(timeoutSec),
	}

	// Add input data if provided
	if inputData != nil {
		executeOptions = append(executeOptions, stack.WithInput(inputData))
	}

	// Execute the stack
	fmt.Println("Starting execution...")
	startTime := time.Now()
	
	err = engine.Execute(ctx, executeOptions...)
	
	duration := time.Since(startTime)
	if err != nil {
		fmt.Printf("Stack execution failed after %v: %v\n", duration, err)
		return err
	}

	// Get execution summary
	summary := engine.GetState()
	
	// Display results
	fmt.Printf("\nStack execution completed in %v\n", duration)
	fmt.Printf("Total agents: %d\n", summary.TotalAgents)
	fmt.Printf("Completed: %d\n", summary.CompletedCount)
	fmt.Printf("Failed: %d\n", summary.FailedCount)
	fmt.Printf("Blocked: %d\n", summary.BlockedCount)
	
	// If verbose, show detailed agent states
	if verbose {
		fmt.Println("\nAgent details:")
		for id, state := range summary.AgentStates {
			fmt.Printf("  - %s: %s\n", id, state.Status)
			if state.Status == stack.AgentStatusFailed && state.ErrorMessage != "" {
				fmt.Printf("    Error: %s\n", state.ErrorMessage)
			}
		}
	}

	return nil
}

// parseInput parses the stack definition from various input sources
func parseInput() (stack.StackSpec, error) {
	p := parser.NewStackParser()
	
	// Check input sources in order of precedence
	if inputFile != "" {
		// Load from file
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return stack.StackSpec{}, fmt.Errorf("failed to read input file: %w", err)
		}
		
		// Determine file type from extension
		switch filepath.Ext(inputFile) {
		case ".json":
			return p.ParseFromJSON(string(content))
		case ".yaml", ".yml":
			return p.ParseFromYAML(string(content))
		default:
			return stack.StackSpec{}, fmt.Errorf("unknown file format: %s", inputFile)
		}
	} else if inputJson != "" {
		// Parse JSON input
		return p.ParseFromJSON(inputJson)
	} else if inputYaml != "" {
		// Parse YAML input
		return p.ParseFromYAML(inputYaml)
	} else if inputNl != "" {
		// Parse natural language description
		return p.ParseFromNaturalLanguage(inputNl)
	} else {
		// Try to read from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Data is being piped to stdin
			content, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return stack.StackSpec{}, fmt.Errorf("failed to read from stdin: %w", err)
			}
			
			// Try to parse as JSON first, then YAML
			spec, err := p.ParseFromJSON(string(content))
			if err == nil {
				return spec, nil
			}
			
			spec, err = p.ParseFromYAML(string(content))
			if err == nil {
				return spec, nil
			}
			
			// Finally try as natural language
			return p.ParseFromNaturalLanguage(string(content))
		}
	}
	
	// No valid input provided
	return stack.StackSpec{}, fmt.Errorf("no stack definition provided - use --file, --json, --yaml, or --nl flags")
}

// ParseStackFile parses a stack file (exposed for testing)
func ParseStackFile(filename string) (stack.StackSpec, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return stack.StackSpec{}, fmt.Errorf("failed to read stack file: %w", err)
	}
	
	var spec stack.StackSpec
	
	// Try to parse as YAML first, then JSON if that fails
	err = yaml.Unmarshal(content, &spec)
	if err != nil {
		// Try JSON
		err = json.Unmarshal(content, &spec)
		if err != nil {
			return stack.StackSpec{}, fmt.Errorf("failed to parse stack file: %w", err)
		}
	}
	
	return spec, nil
}
