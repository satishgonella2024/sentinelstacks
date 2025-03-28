package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/satishgonella2024/sentinelstacks/internal/agent"
	"github.com/satishgonella2024/sentinelstacks/internal/registry"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'registry' or 'agent' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "registry":
		if len(os.Args) < 3 {
			fmt.Println("expected 'login', 'push', 'pull', or 'list' subcommands")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "login":
			loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
			username := loginCmd.String("username", "", "Registry username")
			password := loginCmd.String("password", "", "Registry password")
			server := loginCmd.String("server", "https://localhost", "Registry server URL")
			loginCmd.Parse(os.Args[3:])

			if *username == "" || *password == "" {
				fmt.Println("username and password are required")
				os.Exit(1)
			}

			err := registry.Login(*server, *username, *password)
			if err != nil {
				fmt.Printf("Login failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Login successful!")

		case "push":
			pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
			agentName := pushCmd.String("name", "", "Agent name")
			version := pushCmd.String("version", "latest", "Agent version")
			pushCmd.Parse(os.Args[3:])

			if *agentName == "" {
				fmt.Println("agent name is required")
				os.Exit(1)
			}

			err := registry.Push(*agentName, *version)
			if err != nil {
				fmt.Printf("Push failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully pushed %s:%s\n", *agentName, *version)

		case "pull":
			pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
			agentName := pullCmd.String("name", "", "Agent name")
			version := pullCmd.String("version", "latest", "Agent version")
			pullCmd.Parse(os.Args[3:])

			if *agentName == "" {
				fmt.Println("agent name is required")
				os.Exit(1)
			}

			err := registry.Pull(*agentName, *version)
			if err != nil {
				fmt.Printf("Pull failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully pulled %s:%s\n", *agentName, *version)

		case "list":
			agents, err := registry.List()
			if err != nil {
				fmt.Printf("List failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Available agents:")
			for _, a := range agents {
				fmt.Printf("- %s:%s\n", a.Name, a.Version)
			}

		default:
			fmt.Println("expected 'login', 'push', 'pull', or 'list' subcommands")
			os.Exit(1)
		}

	case "agent":
		if len(os.Args) < 3 {
			fmt.Println("expected 'run' subcommand")
			os.Exit(1)
		}

		switch os.Args[2] {
		case "run":
			runCmd := flag.NewFlagSet("run", flag.ExitOnError)
			agentName := runCmd.String("name", "", "Agent name")
			version := runCmd.String("version", "latest", "Agent version")
			runCmd.Parse(os.Args[3:])

			if *agentName == "" {
				fmt.Println("agent name is required")
				os.Exit(1)
			}

			err := agent.Run(*agentName, *version)
			if err != nil {
				fmt.Printf("Run failed: %v\n", err)
				os.Exit(1)
			}

		default:
			fmt.Println("expected 'run' subcommand")
			os.Exit(1)
		}

	default:
		fmt.Println("expected 'registry' or 'agent' subcommands")
		os.Exit(1)
	}
}
