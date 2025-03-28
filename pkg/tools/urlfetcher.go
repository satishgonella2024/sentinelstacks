package tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// URLFetcherTool provides capabilities to fetch content from URLs
type URLFetcherTool struct {
	client *http.Client
}

// NewURLFetcherTool creates a new URL fetcher tool
func NewURLFetcherTool() *URLFetcherTool {
	return &URLFetcherTool{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ID returns the unique identifier for the URL fetcher tool
func (u *URLFetcherTool) ID() string {
	return "urlfetcher"
}

// Name returns a user-friendly name
func (u *URLFetcherTool) Name() string {
	return "URL Fetcher"
}

// Description returns a detailed description
func (u *URLFetcherTool) Description() string {
	return "Fetches content from URLs, allowing agents to access web resources"
}

// Version returns the semantic version
func (u *URLFetcherTool) Version() string {
	return "0.1.0"
}

// ParameterSchema returns the JSON schema for parameters
func (u *URLFetcherTool) ParameterSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type": "string",
				"description": "The URL to fetch content from",
			},
			"method": map[string]interface{}{
				"type": "string",
				"enum": []string{"GET", "HEAD"},
				"default": "GET",
				"description": "HTTP method to use",
			},
			"maxLength": map[string]interface{}{
				"type": "integer",
				"default": 4096,
				"description": "Maximum length of content to return (to avoid excessive responses)",
			},
		},
		"required": []string{"url"},
	}
}

// Execute runs the URL fetcher with the provided parameters
func (u *URLFetcherTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Get URL parameter
	urlStr, ok := params["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url parameter is required and must be a string")
	}

	// Validate URL
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return nil, fmt.Errorf("url must begin with http:// or https://")
	}

	// Get method parameter (default to GET)
	method := "GET"
	if methodParam, ok := params["method"].(string); ok {
		method = strings.ToUpper(methodParam)
		if method != "GET" && method != "HEAD" {
			return nil, fmt.Errorf("method must be GET or HEAD")
		}
	}

	// Get maxLength parameter (default to 4096)
	maxLength := 4096
	if maxLengthParam, ok := params["maxLength"].(float64); ok {
		maxLength = int(maxLengthParam)
		if maxLength <= 0 {
			return nil, fmt.Errorf("maxLength must be positive")
		}
	}

	// Create request
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set a user agent to be polite
	req.Header.Set("User-Agent", "SentinelStacks-Agent/0.1")

	// Execute request
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-success status code: %d %s", 
			resp.StatusCode, resp.Status)
	}

	// For HEAD requests, just return headers
	if method == "HEAD" {
		headers := make(map[string]string)
		for name, values := range resp.Header {
			headers[name] = strings.Join(values, ", ")
		}
		return map[string]interface{}{
			"status": resp.Status,
			"statusCode": resp.StatusCode,
			"headers": headers,
		}, nil
	}

	// For GET requests, read body
	limitedReader := io.LimitReader(resp.Body, int64(maxLength))
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Create a result with meta information and content
	result := map[string]interface{}{
		"status": resp.Status,
		"statusCode": resp.StatusCode,
		"contentType": resp.Header.Get("Content-Type"),
		"content": string(body),
		"truncated": len(body) >= maxLength,
	}

	return result, nil
}
