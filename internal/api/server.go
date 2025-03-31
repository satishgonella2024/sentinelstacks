package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	goruntime "runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sentinelstacks/sentinel/internal/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SentinelStacks API
// @version 1.0
// @description API for managing AI agents in SentinelStacks

// @contact.name API Support
// @contact.url https://github.com/sentinelstacks/sentinel
// @contact.email support@sentinelstacks.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1
// @schemes http https

// Server represents the API server
type Server struct {
	router    *mux.Router
	server    *http.Server
	runtime   *runtime.Runtime
	config    *Config
	log       *log.Logger
	once      sync.Once
	wsManager *WebSocketManager
}

// Config contains API server configuration
type Config struct {
	Host            string
	Port            int
	TLSCertFile     string
	TLSKeyFile      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	TokenAuthSecret string
	EnableCORS      bool
	LogRequests     bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            8080,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		EnableCORS:      true,
		LogRequests:     true,
	}
}

// NewServer creates a new API server with the given configuration
func NewServer(config *Config) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}

	r, err := runtime.GetRuntime()
	if err != nil {
		return nil, fmt.Errorf("failed to get runtime: %w", err)
	}

	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	s := &Server{
		router:    mux.NewRouter(),
		runtime:   r,
		config:    config,
		log:       logger,
		wsManager: NewWebSocketManager(logger),
	}

	s.setupRoutes()

	return s, nil
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Apply common middleware
	s.router.Use(s.loggingMiddleware)
	if s.config.EnableCORS {
		s.router.Use(s.corsMiddleware)
	}
	s.router.Use(s.recoveryMiddleware)

	// Swagger documentation
	currentDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", currentDir)

	// Use absolute path for docs directory
	docDir := filepath.Join(currentDir, "docs")
	log.Printf("Serving Swagger docs from: %s", docDir)

	s.router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir(docDir))))

	// Swagger documentation endpoint
	// Serve static files for the Swagger UI
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to the API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// API versioning - all routes go under /v1
	api := s.router.PathPrefix("/v1").Subrouter()

	// Authentication endpoint
	api.HandleFunc("/auth/login", s.handleLogin).Methods("POST")

	// Agent routes
	agents := api.PathPrefix("/agents").Subrouter()
	agents.HandleFunc("", s.listAgentsHandler).Methods("GET")
	agents.HandleFunc("", s.createAgentHandler).Methods("POST")
	agents.HandleFunc("/{id}", s.getAgentHandler).Methods("GET")
	agents.HandleFunc("/{id}", s.deleteAgentHandler).Methods("DELETE")
	agents.HandleFunc("/{id}/logs", s.getAgentLogsHandler).Methods("GET")

	// WebSocket routes
	agents.HandleFunc("/{id}/chat", s.HandleAgentChat).Methods("GET")
	agents.HandleFunc("/{id}/events", s.HandleAgentEvents).Methods("GET")

	// Image routes
	images := api.PathPrefix("/images").Subrouter()
	images.HandleFunc("", s.listImagesHandler).Methods("GET")
	images.HandleFunc("/{id}", s.getImageHandler).Methods("GET")

	// Registry routes (protected by auth)
	registry := api.PathPrefix("/registry").Subrouter()
	registry.Use(s.authMiddleware)
	registry.HandleFunc("/search", s.searchRegistry).Methods("GET")
	registry.HandleFunc("/push", s.pushImage).Methods("POST")
	registry.HandleFunc("/pull", s.pullImage).Methods("POST")

	// Health check
	api.HandleFunc("/health", s.healthCheck).Methods("GET")
}

// Start starts the API server
func (s *Server) Start() error {
	var err error

	// Start the Swagger documentation server on port 8081
	StartSwaggerServer(8081)

	// Start the API server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Use the server's router for the main API
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	s.log.Printf("API server starting on %s", addr)

	if s.config.TLSCertFile != "" && s.config.TLSKeyFile != "" {
		err = s.server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
	} else {
		err = s.server.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the API server
func (s *Server) Stop() error {
	s.log.Println("Stopping API server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// RunWithGracefulShutdown runs the server and handles graceful shutdown on SIGINT/SIGTERM
func (s *Server) RunWithGracefulShutdown() error {
	// Channel to listen for errors from Start
	errChan := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Block until we receive a signal or an error
	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		s.log.Println("Received interrupt signal, shutting down...")
		return s.Stop()
	}
}

// Middleware for logging requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.LogRequests {
			s.log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware for CORS
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the origin of the request for debugging
		origin := r.Header.Get("Origin")
		if origin != "" {
			s.log.Printf("Request from origin: %s", origin)
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			s.log.Printf("Handling OPTIONS preflight request from %s", r.RemoteAddr)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware for panic recovery
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.Printf("Panic recovered in API handler: %v", err)
				// Log stack trace to help with debugging
				buf := make([]byte, 4096)
				n := goruntime.Stack(buf, false)
				s.log.Printf("Stack trace: %s", buf[:n])

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// sendJSON sends a JSON response
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		var jsonBytes []byte
		var err error

		// Try to marshal the data
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			s.log.Printf("Error marshaling JSON response: %v (data: %+v)", err, data)
			// If we can't marshal the original data, send a simplified error response
			w.WriteHeader(http.StatusInternalServerError)
			fallbackJSON := []byte(`{"error":"Internal server error: failed to serialize response"}`)
			w.Write(fallbackJSON)
			return
		}

		// Write the JSON data directly
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.log.Printf("Error writing JSON response: %v", err)
		}
	}
}

// sendError sends an error response
func (s *Server) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, map[string]string{"error": message})
}

// Health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// These are implemented in their respective files
// agentsHandler in agents.go
// imagesHandler in images.go
// registryHandler in registry.go
// authHandler in auth.go
// websocketHandler in websocket.go
