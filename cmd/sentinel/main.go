package main

import (
	"fmt"
	"os"

	"github.com/satishgonella2024/sentinelstacks/cmd/sentinel/commands"
)

func main() {
	rootCmd := commands.NewRootCommand()
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
