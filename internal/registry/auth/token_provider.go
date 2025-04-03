// Package auth provides authentication for the registry
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// TokenProvider implements the types.AuthProvider interface for registry authentication
type TokenProvider struct {
	baseURL      string
	tokenPath    string
	clientID     string
	clientSecret string
	token        string
	tokenExpiry  time.Time
	httpClient   *http.Client
	mu           sync.RWMutex
}

// TokenInfo represents a stored authentication token
type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Username  string    `json:"username"`
}

// NewTokenProvider creates a new token provider
func NewTokenProvider(baseURL, tokenPath, clientID, clientSecret string) *TokenProvider {
	// If token path is not specified, use default path
	if tokenPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			tokenPath = filepath.Join(homeDir, ".sentinel", "auth", "token.json")
		}
	}

	return &TokenProvider{
		baseURL:      baseURL,
		tokenPath:    tokenPath,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetToken returns an existing token or attempts to load one
func (p *TokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	token := p.token
	expiry := p.tokenExpiry
	p.mu.RUnlock()

	// If token exists and is not expired, return it
	if token != "" && expiry.After(time.Now()) {
		return token, nil
	}

	// Try to load token from file
	if err := p.loadToken(); err == nil {
		p.mu.RLock()
		token = p.token
		expiry = p.tokenExpiry
		p.mu.RUnlock()

		if token != "" && expiry.After(time.Now()) {
			return token, nil
		}
	}

	// Token doesn't exist or is expired
	return "", fmt.Errorf("authentication required")
}

// Login performs authentication and stores the token
func (p *TokenProvider) Login(ctx context.Context, username, password string) (string, error) {
	// Prepare request
	url := fmt.Sprintf("%s/auth/login", p.baseURL)
	reqData := map[string]string{
		"username": username,
		"password": password,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if p.clientID != "" {
		req.Header.Set("X-Client-ID", p.clientID)
	}
	if p.clientSecret != "" {
		req.Header.Set("X-Client-Secret", p.clientSecret)
	}

	// Send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: %s", string(body))
	}

	// Parse response
	var tokenResp struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expiresIn"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Store token
	p.mu.Lock()
	p.token = tokenResp.Token
	p.tokenExpiry = expiresAt
	p.mu.Unlock()

	// Save token to file
	if err := p.saveToken(username); err != nil {
		// Just log the error, but don't fail the login
		fmt.Printf("Warning: Failed to save token: %v\n", err)
	}

	return tokenResp.Token, nil
}

// Logout invalidates the current token
func (p *TokenProvider) Logout(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Clear token
	p.token = ""
	p.tokenExpiry = time.Time{}

	// Delete token file
	if p.tokenPath != "" {
		if err := os.Remove(p.tokenPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete token file: %w", err)
		}
	}

	return nil
}

// IsAuthenticated checks if the client is authenticated
func (p *TokenProvider) IsAuthenticated(ctx context.Context) bool {
	token, err := p.GetToken(ctx)
	return err == nil && token != ""
}

// loadToken loads a token from the token file
func (p *TokenProvider) loadToken() error {
	// Check if token path is set
	if p.tokenPath == "" {
		return fmt.Errorf("token path not set")
	}

	// Check if token file exists
	if _, err := os.Stat(p.tokenPath); os.IsNotExist(err) {
		return fmt.Errorf("token file not found")
	}

	// Read token file
	data, err := ioutil.ReadFile(p.tokenPath)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	// Parse token data
	var tokenInfo TokenInfo
	if err := json.Unmarshal(data, &tokenInfo); err != nil {
		return fmt.Errorf("failed to parse token file: %w", err)
	}

	// Store token
	p.mu.Lock()
	p.token = tokenInfo.Token
	p.tokenExpiry = tokenInfo.ExpiresAt
	p.mu.Unlock()

	return nil
}

// saveToken saves a token to the token file
func (p *TokenProvider) saveToken(username string) error {
	// Check if token path is set
	if p.tokenPath == "" {
		return fmt.Errorf("token path not set")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(p.tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Create token data
	tokenInfo := TokenInfo{
		Token:     p.token,
		ExpiresAt: p.tokenExpiry,
		Username:  username,
	}

	// Marshal token data
	data, err := json.Marshal(tokenInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Write token file
	if err := ioutil.WriteFile(p.tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
