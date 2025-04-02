package shell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/satishgonella2024/sentinelstacks/internal/runtime"
)

// NewShellCmd creates the shell command
func NewShellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shell [agent_id]",
		Short: "Start an interactive shell with an agent",
		Long:  `Start an interactive shell session with a running agent`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentID := args[0]
			
			// Check if the agent exists
			rt, err := runtime.GetRuntime()
			if err != nil {
				return fmt.Errorf("failed to get runtime: %w", err)
			}
			
			agent, err := rt.GetAgent(agentID)
			if err != nil {
				return fmt.Errorf("agent not found: %s", agentID)
			}
			
			fmt.Printf("Starting shell session with agent '%s' (%s)\n", agent.ID, agent.Name)
			fmt.Println("Type 'exit' to end the session, 'help' for available commands")
			
			// Create context with cancellation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			
			// Handle termination signals
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
			
			go func() {
				<-signalCh
				fmt.Println("\nReceived termination signal, ending shell session...")
				cancel()
			}()
			
			// Start shell loop
			return runShellLoop(ctx, agentID, agent.Name)
		},
	}
}

// runShellLoop runs the interactive shell loop
func runShellLoop(ctx context.Context, agentID, agentName string) error {
	reader := bufio.NewReader(os.Stdin)
	
	// Print the shell prompt with agent name
	fmt.Printf("\n[%s]> ", agentName)
	
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shell session terminated")
			return nil
		default:
			// Read input
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			
			// Trim whitespace
			input = strings.TrimSpace(input)
			
			// Handle empty input
			if input == "" {
				fmt.Printf("[%s]> ", agentName)
				continue
			}
			
			// Handle special commands
			if input == "exit" {
				fmt.Println("Exiting shell session")
				return nil
			}
			
			if input == "help" {
				printShellHelp()
				fmt.Printf("[%s]> ", agentName)
				continue
			}
			
			// TODO: In a real implementation, send the input to the agent and get a response
			
			// Simulate agent response
			fmt.Printf("\nAgent: I received your message: \"%s\"\n", input)
			fmt.Printf("As this is a simulation, I'm not actually processing your request.\n")
			fmt.Printf("In a complete implementation, I would use my capabilities to respond.\n")
			
			// Print prompt for next input
			fmt.Printf("\n[%s]> ", agentName)
		}
	}
}

// printShellHelp prints help information for shell commands
func printShellHelp() {
	fmt.Println("\nAvailable shell commands:")
	fmt.Println("  exit       - Exit the shell session")
	fmt.Println("  help       - Show this help message")
	fmt.Println("  clear      - Clear the screen")
	fmt.Println("  memory     - Display agent's memory state")
	fmt.Println("  capabilities - List agent's capabilities")
	fmt.Println("  tools      - List available tools")
	fmt.Println("  status     - Show agent status information")
	fmt.Println("  history    - Show conversation history")
	fmt.Println()
}
