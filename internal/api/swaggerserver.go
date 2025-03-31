package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// StartSwaggerServer starts a simple HTTP server to serve the Swagger documentation
func StartSwaggerServer(port int) {
	currentDir, _ := os.Getwd()
	docsDir := filepath.Join(currentDir, "docs")

	// Create a file server to serve the Swagger documentation
	fs := http.FileServer(http.Dir(docsDir))

	// Create a simple HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting Swagger documentation server on http://localhost%s", addr)
	log.Printf("Serving documentation from: %s", docsDir)

	// Start the server in a goroutine
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("Error starting Swagger server: %v", err)
		}
	}()
}
