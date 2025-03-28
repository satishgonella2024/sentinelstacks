package main

import (
	"flag"
	"log"

	"github.com/satishgonella2024/sentinelstacks/internal/api"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	if err := api.StartServer(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
