package main

import (
	"fmt"
	"os"

	"github.com/satishgonella2024/sentinelstacks/cmd/sentinel"
)

func main() {
	if err := sentinel.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
