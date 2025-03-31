package api

import (
	"fmt"
	"os"
	"time"

	"github.com/sentinelstacks/sentinel/internal/api"
	"github.com/spf13/cobra"
)

// NewAPICmd creates a new API command
func NewAPICmd() *cobra.Command {
	var (
		port            int
		host            string
		tlsCertFile     string
		tlsKeyFile      string
		tokenAuthSecret string
		enableCORS      bool
		logRequests     bool
	)

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the SentinelStacks API server",
		Long:  `Start the SentinelStacks API server, which exposes a RESTful API for managing agents and images.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Configure API server
			config := &api.Config{
				Host:            host,
				Port:            port,
				TLSCertFile:     tlsCertFile,
				TLSKeyFile:      tlsKeyFile,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				ShutdownTimeout: 30 * time.Second,
				TokenAuthSecret: tokenAuthSecret,
				EnableCORS:      enableCORS,
				LogRequests:     logRequests,
			}

			// If no token auth secret is provided, generate a random one
			if config.TokenAuthSecret == "" {
				config.TokenAuthSecret = fmt.Sprintf("sentinel-secret-%d", time.Now().UnixNano())
				fmt.Fprintf(os.Stderr, "Warning: No token auth secret provided. Using generated secret.\n")
				fmt.Fprintf(os.Stderr, "For production use, set a fixed secret using --token-auth-secret.\n")
			}

			// Create and start server
			server, err := api.NewServer(config)
			if err != nil {
				return fmt.Errorf("failed to create API server: %w", err)
			}

			fmt.Printf("Starting SentinelStacks API server on %s:%d\n", host, port)
			if enableCORS {
				fmt.Println("CORS is enabled")
			}

			return server.RunWithGracefulShutdown()
		},
	}

	// Add flags
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	cmd.Flags().StringVar(&host, "host", "localhost", "Host address to listen on")
	cmd.Flags().StringVar(&tlsCertFile, "tls-cert", "", "TLS certificate file")
	cmd.Flags().StringVar(&tlsKeyFile, "tls-key", "", "TLS key file")
	cmd.Flags().StringVar(&tokenAuthSecret, "token-auth-secret", "", "Secret for JWT token authentication")
	cmd.Flags().BoolVar(&enableCORS, "cors", true, "Enable CORS")
	cmd.Flags().BoolVar(&logRequests, "log-requests", true, "Log API requests")

	return cmd
}
