package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/registry/format"
	packages "github.com/satishgonella2024/sentinelstacks/internal/registry/package"
	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

// RegistryClient handles communication with the registry server
type RegistryClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
	UserAgent  string
	CachePath  string
}

// SearchResult contains search results from the registry
type SearchResult struct {
	TotalCount int               `json:"totalCount"`
	Items      []PackageSummary  `json:"items"`
}

// PackageSummary contains basic information about a package
type PackageSummary struct {
	Name          string                 `json:"name"`
	Type          packages.PackageType   `json:"type"`
	Version       string                 `json:"version"`
	Description   string                 `json:"description"`
	Author        string                 `json:"author"`
	CreatedAt     time.Time              `json:"createdAt"`
	Downloads     int                    `json:"downloads"`
	Labels        map[string]string      `json:"labels,omitempty"`
	Dependencies  []packages.Dependency  `json:"dependencies,omitempty"`
	Verified      bool                   `json:"verified"`
}

// NewRegistryClient creates a new registry client
func NewRegistryClient(baseURL, authToken string) *RegistryClient {
	// Get default cache path
	cachePath := os.Getenv("SENTINEL_CACHE_DIR")
	if cachePath == "" {
		// Use default location
		home, err := os.UserHomeDir()
		if err == nil {
			cachePath = filepath.Join(home, ".sentinel", "cache")
		}
	}

	// Make sure cache path exists
	if cachePath != "" {
		os.MkdirAll(cachePath, 0755)
	}

	return &RegistryClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 120 * time.Second},
		AuthToken:  authToken,
		UserAgent:  "sentinel-registry-client/1.0",
		CachePath:  cachePath,
	}
}

// SetAuthToken sets the authentication token
func (c *RegistryClient) SetAuthToken(token string) {
	c.AuthToken = token
}

// PushPackage pushes a package to the registry
func (c *RegistryClient) PushPackage(ctx context.Context, packagePath string) error {
	// Open package file
	file, err := os.Open(packagePath)
	if err != nil {
		return fmt.Errorf("failed to open package file: %w", err)
	}
	defer file.Close()

	// Create multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add package file to form
	part, err := writer.CreateFormFile("package", filepath.Base(packagePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/packages/publish", c.BaseURL), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", c.UserAgent)
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("push failed: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	return nil
}

// PushStack pushes a stack to the registry
func (c *RegistryClient) PushStack(ctx context.Context, stackSpec stack.StackSpec, author string) error {
	// Create a temporary file for the package
	tempDir, err := os.MkdirTemp("", "sentinel-stack-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Determine package filename
	packageName := format.GetDefaultFilename(stackSpec.Name, stackSpec.Version, "stack")
	packagePath := filepath.Join(tempDir, packageName)

	// Build the package from the stack spec
	if err := packages.BuildFromStackSpec(stackSpec, packagePath, author); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}

	// Push the package
	if err := c.PushPackage(ctx, packagePath); err != nil {
		return fmt.Errorf("failed to push package: %w", err)
	}

	return nil
}

// PullPackage pulls a package from the registry
func (c *RegistryClient) PullPackage(ctx context.Context, name, version, outputPath string) error {
	// If no version specified, use latest
	if version == "" {
		version = "latest"
	}

	// Set up proper content type
	fileType := ""
	if strings.Contains(name, ".agent") {
		fileType = "agent"
	} else if strings.Contains(name, ".stack") {
		fileType = "stack"
	} else {
		// Try to determine from name
		fileType = "package"
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/packages/%s/%s?type=%s", c.BaseURL, name, version, fileType), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.UserAgent)
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pull failed: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy response to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response data: %w", err)
	}

	return nil
}

// PullStack pulls a stack from the registry and returns its specification
func (c *RegistryClient) PullStack(ctx context.Context, name, version string) (*stack.StackSpec, string, error) {
	// Create a temporary directory for the package
	tempDir, err := os.MkdirTemp("", "sentinel-stack-")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Determine package filename
	packageName := format.GetDefaultFilename(name, version, "stack")
	packagePath := filepath.Join(tempDir, packageName)

	// Pull the package
	if err := c.PullPackage(ctx, name, version, packagePath); err != nil {
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("failed to pull package: %w", err)
	}

	// Extract the package
	extractDir := filepath.Join(tempDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("failed to create extraction directory: %w", err)
	}

	// Create a package reader to extract the package
	pkg := &packages.SentinelPackage{}
	if err := pkg.Unpackage(packagePath, extractDir); err != nil {
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("failed to extract package: %w", err)
	}

	// Find the stack definition file
	var stackFilePath string
	for _, file := range pkg.Manifest.Files {
		if file.IsMain && strings.HasSuffix(file.Path, format.StackDefinitionExtension) {
			stackFilePath = filepath.Join(extractDir, file.Path)
			break
		}
	}

	if stackFilePath == "" {
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("no stack definition found in package")
	}

	// Parse the stack file
	stackContent, err := os.ReadFile(stackFilePath)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("failed to read stack file: %w", err)
	}

	var stackSpec stack.StackSpec
	if err := json.Unmarshal(stackContent, &stackSpec); err != nil {
		// Try to parse as YAML
		// In a real implementation, this would use yaml.Unmarshal
		os.RemoveAll(tempDir)
		return nil, "", fmt.Errorf("failed to parse stack definition: %w", err)
	}

	// Don't remove the temp directory, it contains the stack definition
	return &stackSpec, stackFilePath, nil
}

// SearchPackages searches for packages in the registry
func (c *RegistryClient) SearchPackages(ctx context.Context, query string, packageType packages.PackageType, limit int) (*SearchResult, error) {
	// Create URL with query parameters
	params := url.Values{}
	if query != "" {
		params.Set("q", query)
	}
	if packageType != "" {
		params.Set("type", string(packageType))
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	// Build URL
	searchURL := fmt.Sprintf("%s/v1/packages/search", c.BaseURL)
	if len(params) > 0 {
		searchURL += "?" + params.Encode()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.UserAgent)
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Parse response
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetPackageVersions gets all versions of a package
func (c *RegistryClient) GetPackageVersions(ctx context.Context, name string) ([]string, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/packages/%s/versions", c.BaseURL, name), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.UserAgent)
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("version lookup failed: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Parse response
	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return versions, nil
}

// Authenticate authenticates with the registry server
func (c *RegistryClient) Authenticate(ctx context.Context, username, password string) (string, error) {
	// Create request payload
	auth := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}

	// Marshal payload
	authBytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth data: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/auth/login", c.BaseURL), bytes.NewBuffer(authBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	// Parse response
	var result struct {
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Set client token
	c.AuthToken = result.Token

	return result.Token, nil
}

// VerifyPackage verifies a package's signature and integrity
func (c *RegistryClient) VerifyPackage(packagePath string) (bool, []string, error) {
	// Create a package reader
	pkg := &packages.SentinelPackage{}

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "sentinel-verify-")
	if err != nil {
		return false, nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Unpackage the file
	if err := pkg.Unpackage(packagePath, tempDir); err != nil {
		return false, nil, fmt.Errorf("failed to extract package: %w", err)
	}

	// Check integrity
	valid, failures, err := pkg.VerifyIntegrity(tempDir)
	if err != nil {
		return false, nil, fmt.Errorf("integrity check failed: %w", err)
	}

	if !valid {
		return false, failures, fmt.Errorf("package integrity verification failed")
	}

	// Return success
	return true, nil, nil
}
