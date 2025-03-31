package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sentinelstacks/sentinel/internal/runtime"
)

// Server represents the API server
type Server struct {
	router  *mux.Router
	server  *http.Server
	runtime *runtime.Runtime
	config  *Config
	log     *log.Logger
	once    sync.Once
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

	s := &Server{
		router:  mux.NewRouter(),
		runtime: r,
		config:  config,
		log:     log.New(os.Stdout, "[API] ", log.LstdFlags),
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
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	s.log.Printf("API server starting on %s", addr)

	var err error
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
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
				s.log.Printf("Panic recovered: %v", err)
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
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.log.Printf("Error encoding JSON response: %v", err)
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

// Route handlers (to be implemented in separate files)
func (s *Server) listAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := s.runtime.GetRunningAgents()
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to get agents")
		return
	}
	s.sendJSON(w, http.StatusOK, map[string]interface{}{"agents": agents})
}

func (s *Server) getAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	agent, err := s.runtime.GetAgent(id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	s.sendJSON(w, http.StatusOK, agent)
}

func (s *Server) createAgent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Not implemented yet")
}

func (s *Server) deleteAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.runtime.DeleteAgent(id)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete agent: %v", err))
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]string{"status": "Agent deleted"})
}

func (s *Server) getAgentLogs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Not implemented yet")
}

func (s *Server) listImages(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Not implemented yet")
}

func (s *Server) getImage(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Not implemented yet")
}

// Route handlers for registry operations (to be implemented in detail later)
func (s *Server) searchRegistry(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Registry search not implemented yet")
}

func (s *Server) pushImage(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Image push not implemented yet")
}

func (s *Server) pullImage(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	s.sendError(w, http.StatusNotImplemented, "Image pull not implemented yet")
}
