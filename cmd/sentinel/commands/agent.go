package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/pkg/runtime"
)

// AgentCmd returns the agent command
func AgentCmd() *cobra.Command {
	var modelEndpoint string

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage and run agents",
		Long:  `Run, monitor, and manage your AI agents.`,
	}

	runCmd := &cobra.Command{
		Use:   "run [agent-name]",
		Short: "Run an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			agentName := args[0]
			fmt.Printf("Running agent: %s\n", agentName)
			
			// Check if agent directory exists
			agentDir := agentName
			if _, err := os.Stat(agentDir); os.IsNotExist(err) {
				fmt.Printf("Error: Agent '%s' does not exist\n", agentName)
				os.Exit(1)
			}
			
			// Check if agentfile.yaml exists
			agentfilePath := filepath.Join(agentDir, "agentfile.yaml")
			if _, err := os.Stat(agentfilePath); os.IsNotExist(err) {
				fmt.Printf("Error: Agentfile for '%s' does not exist\n", agentName)
				os.Exit(1)
			}
			
			// Create agent runtime
			runtime := runtime.NewAgentRuntime()
			
			// Load agentfile
			err := runtime.LoadAgentfile(agentfilePath)
			if err != nil {
				fmt.Printf("Error loading agentfile: %v\n", err)
				os.Exit(1)
			}
			
			// Get the model endpoint from flag or use default
			endpoint, _ := cmd.Flags().GetString("endpoint")
			if endpoint != "" {
				// Override the endpoint if specified
				runtime.ModelEndpoint = endpoint
			}
			
			// Initialize the runtime
			err = runtime.Initialize()
			if err != nil {
				fmt.Printf("Error initializing agent: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Agent '%s' is ready. Type 'exit' to quit.\n\n", agentName)
			
			// Start interactive loop
			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print("> ")
				if !scanner.Scan() {
					break
				}
				
				input := scanner.Text()
				if input == "exit" {
					break
				}
				
				// Process the input
				response, err := runtime.Run(input)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
				
				fmt.Println(response)
			}
		},
	}

	runCmd.Flags().StringVarP(&modelEndpoint, "endpoint", "e", "", "Override the model endpoint URL")
	
	cmd.AddCommand(runCmd)
	return cmd
}
