package main

import (
	"fmt"
	"os"

	"github.com/sentinelstacks/sentinel/cmd/sentinel"
)

func main() {
	if err := sentinel.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
