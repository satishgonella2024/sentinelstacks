// Package auth provides authentication mechanisms for the registry
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// AuthConfig contains configuration for authentication
type AuthConfig struct {
	RegistryURL string
	TokenFile   string
}

// FileTokenProvider implements types.AuthProvider with file-based token storage
type FileTokenProvider struct {
	registryURL string
	tokenFile   string
	token       string
}

// TokenData represents the stored token data
type TokenData struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// NewFileTokenProvider creates a new file-based token provider
func NewFileTokenProvider(config AuthConfig) (*FileTokenProvider, error) {
	// Set default token file if not provided
	tokenFile := config.TokenFile
	if tokenFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		tokenFile = filepath.Join(homeDir, ".sentinel", "registry", "auth.json")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(tokenFile), 0700); err != nil {
		return nil, fmt.Errorf("failed to create token directory: %w", err)
	}

	// Create provider
	provider := &FileTokenProvider{
		registryURL: config.RegistryURL,
		tokenFile:   tokenFile,
	}

	// Load token if exists
	provider.loadToken()

	return provider, nil
}

// GetToken returns a valid authentication token
func (p *FileTokenProvider) GetToken(ctx context.Context) (string, error) {
	// Check if token is loaded and valid
	if p.isTokenValid() {
		return p.token, nil
	}

	// If token is not valid, try to refresh from file
	if err := p.loadToken(); err != nil {
		return "", fmt.Errorf("authentication required: %w", err)
	}

	// Check if refreshed token is valid
	if p.isTokenValid() {
		return p.token, nil
	}

	return "", fmt.Errorf("authentication required")
}

// Login performs authentication and returns a token
func (p *FileTokenProvider) Login(ctx context.Context, username, password string) (string, error) {
	// TODO: Implement actual API call to authenticate
	// For now, generate a demo token

	// Create token claims
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("demo-secret-key"))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token
	p.token = tokenString
	if err := p.saveToken(tokenString, time.Now().Add(24*time.Hour)); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return tokenString, nil
}

// Logout invalidates the current token
func (p *FileTokenProvider) Logout(ctx context.Context) error {
	// Clear token
	p.token = ""

	// Remove token file
	if err := os.Remove(p.tokenFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the client is authenticated
func (p *FileTokenProvider) IsAuthenticated(ctx context.Context) bool {
	return p.isTokenValid()
}

// isTokenValid checks if the current token is valid
func (p *FileTokenProvider) isTokenValid() bool {
	if p.token == "" {
		return false
	}

	// Parse token
	token, err := jwt.Parse(p.token, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return signing key
		return []byte("demo-secret-key"), nil
	})

	// Check if token is valid
	if err != nil {
		return false
	}

	// Check claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			return time.Now().Unix() < int64(exp)
		}
	}

	return false
}

// loadToken loads the token from the token file
func (p *FileTokenProvider) loadToken() error {
	// Check if token file exists
	if _, err := os.Stat(p.tokenFile); os.IsNotExist(err) {
		return fmt.Errorf("token file does not exist")
	}

	// Read token file
	data, err := ioutil.ReadFile(p.tokenFile)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	// Parse token data
	var tokenData TokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return fmt.Errorf("failed to parse token data: %w", err)
	}

	// Check if token is expired
	if tokenData.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token is expired")
	}

	// Set token
	p.token = tokenData.Token

	return nil
}

// saveToken saves the token to the token file
func (p *FileTokenProvider) saveToken(token string, expiresAt time.Time) error {
	// Create token data
	tokenData := TokenData{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	// Marshal token data
	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Write token file
	if err := ioutil.WriteFile(p.tokenFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
