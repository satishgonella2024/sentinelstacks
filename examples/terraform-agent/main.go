package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: terraform-agent <command> [args]")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "plan":
		if len(args) < 1 {
			fmt.Println("Error: path argument required")
			os.Exit(1)
		}
		runTerraform("plan", args[0])

	case "apply":
		if len(args) < 1 {
			fmt.Println("Error: path argument required")
			os.Exit(1)
		}
		runTerraform("apply", args[0])

	case "destroy":
		if len(args) < 1 {
			fmt.Println("Error: path argument required")
			os.Exit(1)
		}
		runTerraform("destroy", args[0])

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runTerraform(command, path string) {
	// Change to the Terraform directory
	if err := os.Chdir(path); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(1)
	}

	// Run Terraform command
	cmd := exec.Command("terraform", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running terraform %s: %v\n", command, err)
		os.Exit(1)
	}
}
