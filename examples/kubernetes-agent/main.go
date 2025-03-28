package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kubernetes-agent <command> [args]")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "deploy":
		if len(args) < 1 {
			fmt.Println("Error: manifest path required")
			os.Exit(1)
		}
		runKubectl("apply", "-f", args[0])

	case "scale":
		if len(args) < 2 {
			fmt.Println("Error: deployment name and replicas required")
			os.Exit(1)
		}
		replicas, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: invalid replicas value: %v\n", err)
			os.Exit(1)
		}
		runKubectl("scale", "deployment", args[0], "--replicas="+strconv.Itoa(replicas))

	case "status":
		if len(args) > 0 {
			runKubectl("get", args[0])
		} else {
			runKubectl("get", "all")
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runKubectl(args ...string) {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running kubectl: %v\n", err)
		os.Exit(1)
	}
}
