// Package main provides the main entry point for the Sentinel Stacks application
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/satishgonella2024/sentinelstacks/pkg/api"
)

var (
	// Command line flags
	configDir = flag.String("config", "config", "Path to configuration directory")
	dataDir   = flag.String("data", "data", "Path to data directory")
	verbose   = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if *verbose {
		log.Println("Verbose logging enabled")
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		fmt.Println("Received termination signal. Shutting down...")
		cancel()
	}()

	// Initialize API
	sentinel, err := initializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer sentinel.Close()

	// Print banner
	printBanner()

	// Print status
	printStatus(sentinel)

	// Wait for termination signal
	fmt.Println("Sentinel Stacks is running. Press Ctrl+C to exit.")
	<-ctx.Done()
	fmt.Println("Shutting down...")
}

// initializeAPI initializes the API with configuration
func initializeAPI() (api.API, error) {
	// Create data directories
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create API configuration
	config := api.APIConfig{
		StackConfig: api.StackServiceConfig{
			StoragePath: filepath.Join(*dataDir, "stacks"),
			Verbose:     *verbose,
		},
		MemoryConfig: api.MemoryServiceConfig{
			StoragePath:         filepath.Join(*dataDir, "memory"),
			EmbeddingProvider:   "local",
			EmbeddingModel:      "local",
			EmbeddingDimensions: 1536,
		},
		RegistryConfig: api.RegistryServiceConfig{
			RegistryURL: "https://registry.sentinelstacks.example",
			CachePath:   filepath.Join(*dataDir, "registry-cache"),
			Username:    os.Getenv("REGISTRY_USERNAME"),
			AccessToken: os.Getenv("REGISTRY_TOKEN"),
		},
	}

	// Create API instance
	sentinel, err := api.NewAPI(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize API: %w", err)
	}

	return sentinel, nil
}

// printBanner prints the application banner
func printBanner() {
	banner := `
 _____            _   _            _   _____ _             _        
/  ___|          | | (_)          | | /  ___| |           | |       
\ \____ ___  _ __| |_ _ _ __   ___| | \ \__ | |_ __ _  ___| | _____ 
 \____ / _ \| '__| __| | '_ \ / _ | |  \__ \| __/ _\` + "`" + ` |/ __| |/ / __|
/\__/ | (_) | |  | |_| | | | |  __| | /\__/ | || (_| | (__|   <\__ \
\____/ \___/|_|   \__|_|_| |_|\___|_| \____/ \__\__,_|\___|_|\_|___/
                                                                     
`
	fmt.Println(banner)
	fmt.Println("Sentinel Stacks - AI Agent Orchestration Platform")
	fmt.Println("Version: 1.0.0")
	fmt.Println("--------------------------------------------")
}

// printStatus prints the application status
func printStatus(sentinel api.API) {
	// Get stack count
	ctx := context.Background()
	stacks, err := sentinel.Stack().ListStacks(ctx)
	if err != nil {
		fmt.Printf("Error getting stacks: %v\n", err)
		return
	}

	fmt.Printf("Configuration directory: %s\n", *configDir)
	fmt.Printf("Data directory: %s\n", *dataDir)
	fmt.Printf("Loaded stacks: %d\n", len(stacks))
	fmt.Println("--------------------------------------------")
}
