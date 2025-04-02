package main

import (
	"fmt"
	"os"
)

func main() {
	// Create a new command
	cmd := NewImagesCmd()
	
	// Execute the command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
