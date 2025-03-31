package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sentinelstacks/sentinel/cmd/sentinel"
	versionCmd "github.com/sentinelstacks/sentinel/cmd/sentinel/version"
)

func main() {
	// Print sentinel logo as ASCII art for non-help commands
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") && os.Args[1] != "help" && os.Args[1] != "completion" {
		printLogo()
	}

	// Execute the command
	if err := sentinel.Execute(); err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), err)
		os.Exit(1)
	}
}

func printLogo() {
	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Println()
	fmt.Println(cyan("  ____             _   _            _  ____  _             _        "))
	fmt.Println(cyan(" / ___|  ___ _ __ | |_(_)_ __   ___| |/ ___|| |_ __ _  ___| | _____ "))
	fmt.Println(cyan(" \\___ \\ / _ \\ '_ \\| __| | '_ \\ / _ \\ |\\___ \\| __/ _` |/ __| |/ / __|"))
	fmt.Println(cyan("  ___) |  __/ | | | |_| | | | |  __/ | ___) | || (_| | (__|   <\\__ \\"))
	fmt.Println(cyan(" |____/ \\___|_| |_|\\__|_|_| |_|\\___|_||____/ \\__\\__,_|\\___|_|\\_\\___/"))
	fmt.Println()
	fmt.Printf("%s %s\n", blue("Version:"), versionCmd.Version)
	fmt.Println()
}
