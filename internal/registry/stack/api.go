package stack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// RegistryClient handles communication with the stack registry API
type RegistryClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// StackInfo provides metadata about a stack in the registry
type StackInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Tags        []string  `json:"tags"`
	Publisher   string    `json:"publisher"`
	CreatedAt   time.Time `json:"createdAt"`
	Downloads   int       `json:"downloads"`
	AgentCount  int       `json:"agentCount"`
	Size        int64     `json:"size"`
}

// SearchResult represents a stack search result
type SearchResult struct {
	TotalCount int        `json:"totalCount"`
	Items      []StackInfo `json:"items"`
}

// NewRegistryClient creates a new stack registry client
func NewRegistryClient(baseURL string, authToken string) *RegistryClient {
	return &RegistryClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		AuthToken:  authToken,
	}
}

// PushStack pushes a stack to the registry
func (c *RegistryClient) PushStack(ctx context.Context, spec stack.StackSpec, stackFilePath string) error {
	// Parse stack name and tag
	stackName := spec.Name
	stackTag := "latest"
	if spec.Version != "" {
		stackTag = normalizeVersion(spec.Version)
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add stack metadata
	metadataBytes, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal stack spec: %w", err)
	}

	metadataField, err := writer.CreateFormField("metadata")
	if err != nil {
		return fmt.Errorf("failed to create metadata field: %w", err)
	}
	if _, err := metadataField.Write(metadataBytes); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Add stack file
	stackFile, err := os.Open(stackFilePath)
	if err != nil {
		return fmt.Errorf("failed to open stack file: %w", err)
	}
	defer stackFile.Close()

	fileField, err := writer.CreateFormFile("stackfile", filepath.Base(stackFilePath))
	if err != nil {
		return fmt.Errorf("failed to create file field: %w", err)
	}
	if _, err := io.Copy(fileField, stackFile); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/v1/stacks/%s/tags/%s", c.BaseURL, stackName, stackTag)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to push stack: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	return nil
}

// PullStack downloads a stack from the registry
func (c *RegistryClient) PullStack(ctx context.Context, name string, tag string) (*stack.StackSpec, string, error) {
	// If no tag provided, use latest
	if tag == "" {
		tag = "latest"
	}

	// Create request
	url := fmt.Sprintf("%s/v1/stacks/%s/tags/%s", c.BaseURL, name, tag)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("failed to pull stack: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse response as multipart
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/") {
		return nil, "", fmt.Errorf("unexpected content type: %s", contentType)
	}

	// Create temp directory for stack file
	tempDir, err := ioutil.TempDir("", "sentinel-stack-")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Parse multipart response
	reader := multipart.NewReader(resp.Body, extractBoundary(contentType))
	
	var spec stack.StackSpec
	stackFilePath := ""

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("failed to read multipart: %w", err)
		}

		// Handle different parts
		switch part.FormName() {
		case "metadata":
			// Parse metadata
			metadataBytes, err := ioutil.ReadAll(part)
			if err != nil {
				return nil, "", fmt.Errorf("failed to read metadata: %w", err)
			}
			if err := json.Unmarshal(metadataBytes, &spec); err != nil {
				return nil, "", fmt.Errorf("failed to unmarshal metadata: %w", err)
			}

		case "stackfile":
			// Save stackfile
			filename := part.FileName()
			if filename == "" {
				filename = "Stackfile.yaml"
			}
			stackFilePath = filepath.Join(tempDir, filename)
			file, err := os.Create(stackFilePath)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create stack file: %w", err)
			}
			if _, err := io.Copy(file, part); err != nil {
				file.Close()
				return nil, "", fmt.Errorf("failed to write stack file: %w", err)
			}
			file.Close()
		}
	}

	if stackFilePath == "" {
		return nil, "", fmt.Errorf("no stack file found in response")
	}

	return &spec, stackFilePath, nil
}

// SearchStacks searches for stacks in the registry
func (c *RegistryClient) SearchStacks(ctx context.Context, query string, limit int) (*SearchResult, error) {
	// Build URL
	params := url.Values{}
	params.Add("q", query)
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	requestURL := fmt.Sprintf("%s/v1/stacks/search?%s", c.BaseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to search stacks: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetStackTags returns the available tags for a stack
func (c *RegistryClient) GetStackTags(ctx context.Context, name string) ([]string, error) {
	// Create request
	url := fmt.Sprintf("%s/v1/stacks/%s/tags", c.BaseURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get stack tags: %s (status code: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var tags []string
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tags, nil
}

// extractBoundary extracts the boundary from a multipart content type
func extractBoundary(contentType string) string {
	parts := strings.Split(contentType, "boundary=")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// normalizeVersion ensures version follows a consistent format
func normalizeVersion(version string) string {
	// If version doesn't start with 'v', add it
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}
