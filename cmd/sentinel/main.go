package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/satishgonella2024/sentinelstacks/internal/agent"
	"github.com/satishgonella2024/sentinelstacks/internal/registry"
	"github.com/satishgonella2024/sentinelstacks/pkg/ui"
	"github.com/spf13/cobra"
)

const asciiLogo = `
    ╔══════════════════════════════════════════════════════════════╗
    ║  _____ _____ _   _ _____ _____ _   _ _____ _      _____ ___║
    ║ /  ___/  ___| \ | |_   _|_   _| \ | |  ___| |    /  ___/ __║
    ║ \ '--.\ '--.|  \| | | |   | | |  \| | |__ | |    \ '--.\ '-║
    ║  '--. \--. \ . ' | | |   | | | . ' |  __|| |     '--. \--. ║
    ║ /\__/ /\__/ | |\  | | |  _| |_| |\  | |___| |____/\__/ \__/║
    ║ \____/\____/\_| \_/ \_/  \___/\_| \_\____/\_____/\____/____/║
    ║                                                              ║
    ║     [S]ecure • [T]rusted • [A]utonomous • [C]onfigurable   ║
    ║              [K]nowledgeable • [S]calable                   ║
    ╚══════════════════════════════════════════════════════════════╝`

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func getColoredLogo() string {
	lines := strings.Split(asciiLogo, "\n")
	var coloredLines []string

	// Border color
	borderColor := color.New(color.FgHiCyan).SprintFunc()

	// Text colors
	mainColor := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	accentColor := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	taglineColor := color.New(color.FgHiGreen).SprintFunc()

	for i, line := range lines {
		if i == 0 || i == len(lines)-1 {
			// Empty lines
			coloredLines = append(coloredLines, line)
			continue
		}

		if strings.Contains(line, "║") {
			// Lines with borders
			content := line[4 : len(line)-1] // Remove border chars
			if strings.Contains(line, "[") {
				// Tagline lines
				coloredContent := taglineColor(content)
				coloredLines = append(coloredLines, fmt.Sprintf("    %s %s %s",
					borderColor("║"), coloredContent, borderColor("║")))
			} else if strings.Contains(line, "SENTINEL") || strings.Contains(line, "STACKS") {
				// Main text lines
				parts := strings.Split(content, "")
				for i, p := range parts {
					if i%2 == 0 {
						parts[i] = mainColor(p)
					} else {
						parts[i] = accentColor(p)
					}
				}
				coloredContent := strings.Join(parts, "")
				coloredLines = append(coloredLines, fmt.Sprintf("    %s %s %s",
					borderColor("║"), coloredContent, borderColor("║")))
			} else {
				// Border lines
				coloredLines = append(coloredLines, fmt.Sprintf("    %s %s %s",
					borderColor("║"), content, borderColor("║")))
			}
		} else {
			// Top/bottom border
			coloredLines = append(coloredLines, borderColor(line))
		}
	}

	return strings.Join(coloredLines, "\n")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("SentinelStacks CLI %s (%s) built on %s\n", version, commit, date)
		return
	}

	fmt.Println("SentinelStacks CLI - AI-powered infrastructure management")

	var rootCmd = &cobra.Command{
		Use:   "sentinel",
		Short: "SentinelStacks - AI Agent Management Platform",
		Long: getColoredLogo() + "\n\n" +
			color.HiWhiteString("Welcome to SentinelStacks - Your AI Agent Management Platform") + "\n" +
			color.HiCyanString("Version: 1.0.0") + "\n\n" +
			"A powerful platform for creating, running, and managing AI agents.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var registryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Manage agent registry",
		Long:  "Interact with the agent registry - login, push, pull, and list agents.",
	}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to registry",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			server, _ := cmd.Flags().GetString("server")

			if username == "" || password == "" {
				color.Red("Error: username and password are required")
				os.Exit(1)
			}

			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle(fmt.Sprintf("Logging in to %s...", server), "dots")
			spinner.Start()

			// Simulate network delay (in a real app, this would be an actual API call)
			err := registry.Login(server, username, password)

			if err != nil {
				spinner.Error(fmt.Sprintf("Login failed: %v", err))
				os.Exit(1)
			}

			spinner.Success("Login successful!")
		},
	}

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push agent to registry",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			version, _ := cmd.Flags().GetString("version")

			if name == "" {
				color.Red("Error: agent name is required")
				os.Exit(1)
			}

			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle(fmt.Sprintf("Pushing %s:%s...", name, version), "smooth")
			spinner.Start()

			// Update the message after a delay to show progress
			go func() {
				time.Sleep(1 * time.Second)
				spinner.UpdateMessage(fmt.Sprintf("Preparing %s:%s for upload...", name, version))

				time.Sleep(1 * time.Second)
				spinner.UpdateMessage(fmt.Sprintf("Uploading %s:%s to registry...", name, version))
			}()

			// Simulate network delay (in a real app, this would be an actual API call)
			err := registry.Push(name, version)

			if err != nil {
				spinner.Error(fmt.Sprintf("Push failed: %v", err))
				os.Exit(1)
			}

			spinner.Success(fmt.Sprintf("Successfully pushed %s:%s", name, version))
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull agent from registry",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			version, _ := cmd.Flags().GetString("version")

			if name == "" {
				color.Red("Error: agent name is required")
				os.Exit(1)
			}

			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle(fmt.Sprintf("Pulling %s:%s...", name, version), "bounce")
			spinner.Start()

			// Update the message after a delay to show progress
			go func() {
				time.Sleep(800 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Downloading %s:%s from registry...", name, version))

				time.Sleep(1200 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Verifying package integrity...", name, version))

				time.Sleep(700 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Installing %s:%s locally...", name, version))
			}()

			// Simulate network delay (in a real app, this would be an actual API call)
			err := registry.Pull(name, version)

			if err != nil {
				spinner.Error(fmt.Sprintf("Pull failed: %v", err))
				os.Exit(1)
			}

			spinner.Success(fmt.Sprintf("Successfully pulled %s:%s", name, version))
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available agents",
		Run: func(cmd *cobra.Command, args []string) {
			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle("Fetching agents...", "classic")
			spinner.Start()

			// Simulate network delay (in a real app, this would be an actual API call)
			agents, err := registry.List()

			if err != nil {
				spinner.Error(fmt.Sprintf("List failed: %v", err))
				os.Exit(1)
			}

			spinner.Success("Successfully fetched agents")

			color.HiWhite("\nAvailable Agents:\n")
			for _, a := range agents {
				fmt.Printf("  %s %s:%s\n",
					color.HiGreenString("•"),
					color.HiWhiteString(a.Name),
					color.HiBlueString(a.Version))
			}
		},
	}

	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Manage agents",
		Long:  "Create, run, and manage AI agents.",
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run an agent",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			version, _ := cmd.Flags().GetString("version")
			interactive, _ := cmd.Flags().GetBool("interactive")

			if name == "" {
				color.Red("Error: agent name is required")
				os.Exit(1)
			}

			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle(fmt.Sprintf("Starting agent %s:%s...", name, version), "arrow")
			spinner.Start()

			// Update messages to show progress
			go func() {
				time.Sleep(800 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Loading agent configuration..."))

				time.Sleep(1200 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Initializing model..."))

				time.Sleep(700 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Preparing runtime environment..."))
			}()

			if interactive {
				// Load the agent
				ag, err := agent.LoadAgent(name, version)
				if err != nil {
					spinner.Error(fmt.Sprintf("Failed to load agent: %v", err))
					os.Exit(1)
				}

				spinner.Success(fmt.Sprintf("Agent %s:%s is ready", name, version))

				color.HiWhite("\nRunning %s in interactive mode. Type 'exit' to quit.\n", name)
				color.HiGreen("Agent is ready. What would you like to do?\n")

				// Interactive loop
				reader := bufio.NewReader(os.Stdin)
				for {
					// Print prompt
					fmt.Print(color.HiCyanString("You: "))

					// Read user input (handle multiline)
					input, err := reader.ReadString('\n')
					if err != nil {
						color.Red("✗ Error reading input: %v\n", err)
						continue
					}

					// Trim whitespace
					input = strings.TrimSpace(input)

					// Check for exit command
					if input == "exit" || input == "quit" {
						break
					}

					// Show thinking spinner
					thinkingSpinner := ui.NewSpinnerWithStyle("Agent is thinking...", "dots")
					thinkingSpinner.Start()

					// Execute agent
					response, err := ag.Execute(input)

					if err != nil {
						thinkingSpinner.Error(fmt.Sprintf("Error: %v", err))
						continue
					}

					thinkingSpinner.Stop()

					// Print response
					fmt.Print(color.HiMagentaString("Agent: "))
					fmt.Println(response)
					fmt.Println()
				}

				color.Green("✓ Agent session completed successfully")
			} else {
				// Run in non-interactive mode
				err := agent.Run(name, version)
				if err != nil {
					spinner.Error(fmt.Sprintf("Run failed: %v", err))
					os.Exit(1)
				}
				spinner.Success("Agent completed successfully")
			}
		},
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			model, _ := cmd.Flags().GetString("model")
			memoryType, _ := cmd.Flags().GetString("memory")

			if name == "" {
				color.Red("Error: agent name is required")
				os.Exit(1)
			}

			// Create and start spinner
			spinner := ui.NewSpinnerWithStyle(fmt.Sprintf("Creating agent %s...", name), "smooth")
			spinner.Start()

			// Update messages to show progress
			go func() {
				time.Sleep(500 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Creating directory structure..."))

				time.Sleep(700 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Generating Agentfile..."))

				time.Sleep(800 * time.Millisecond)
				spinner.UpdateMessage(fmt.Sprintf("Creating example script..."))
			}()

			// Create agent directory
			agentDir := filepath.Join("agents", name)
			if err := os.MkdirAll(agentDir, 0755); err != nil {
				spinner.Error(fmt.Sprintf("Failed to create agent directory: %v", err))
				os.Exit(1)
			}

			// Create Agentfile
			agentfile := fmt.Sprintf(`name: %s
version: "1.0.0"
description: "%s"
model:
  provider: "ollama"
  name: "%s"
capabilities:
  - basic
  - file-io
memory:
  type: "%s"
  persistence: true
`, name, description, model, memoryType)

			if err := os.WriteFile(filepath.Join(agentDir, "Agentfile"), []byte(agentfile), 0644); err != nil {
				spinner.Error(fmt.Sprintf("Failed to create Agentfile: %v", err))
				os.Exit(1)
			}

			// Create example script
			script := fmt.Sprintf(`# %s
# 
# This is an example agent that demonstrates basic capabilities.
# You can modify this file to add your own functionality.

def main():
    print("Hello from %s!")
    # Add your agent logic here

if __name__ == "__main__":
    main()
`, name, name)

			if err := os.WriteFile(filepath.Join(agentDir, "agent.py"), []byte(script), 0644); err != nil {
				spinner.Error(fmt.Sprintf("Failed to create agent script: %v", err))
				os.Exit(1)
			}

			spinner.Success(fmt.Sprintf("Successfully created agent %s", name))

			color.HiWhite("\nNext steps:\n")
			fmt.Printf("1. Edit %s/Agentfile to configure your agent\n", agentDir)
			fmt.Printf("2. Edit %s/agent.py to add your agent's logic\n", agentDir)
			fmt.Printf("3. Run 'sentinel registry push --name %s' to publish your agent\n", name)
		},
	}

	// Add flags
	loginCmd.Flags().String("username", "", "Registry username")
	loginCmd.Flags().String("password", "", "Registry password")
	loginCmd.Flags().String("server", "https://localhost", "Registry server URL")

	pushCmd.Flags().String("name", "", "Agent name")
	pushCmd.Flags().String("version", "latest", "Agent version")

	pullCmd.Flags().String("name", "", "Agent name")
	pullCmd.Flags().String("version", "latest", "Agent version")

	runCmd.Flags().String("name", "", "Agent name")
	runCmd.Flags().String("version", "latest", "Agent version")
	runCmd.Flags().Bool("interactive", false, "Run in interactive mode")

	// Add flags for create command
	createCmd.Flags().String("name", "", "Agent name")
	createCmd.Flags().String("description", "A new SentinelStacks agent", "Agent description")
	createCmd.Flags().String("model", "llama2", "Model to use (e.g. llama2, gpt-4)")
	createCmd.Flags().String("memory", "simple", "Memory type (simple or vector)")

	// Add commands to registry
	registryCmd.AddCommand(loginCmd, pushCmd, pullCmd, listCmd)
	agentCmd.AddCommand(runCmd, createCmd)
	rootCmd.AddCommand(registryCmd, agentCmd)

	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
