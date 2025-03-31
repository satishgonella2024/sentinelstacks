package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 8081, "Port to listen on")
	flag.Parse()

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Determine the docs directory path
	docsDir := filepath.Join(cwd, "docs")
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		log.Fatalf("Docs directory not found at: %s", docsDir)
	}

	// Create a file server for the docs directory
	fileServer := http.FileServer(http.Dir(docsDir))

	// Set up HTTP handlers
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	}))

	// Start the HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting Swagger documentation server on http://localhost%s", addr)
	log.Printf("Serving documentation from: %s", docsDir)
	log.Fatal(http.ListenAndServe(addr, nil))
}
