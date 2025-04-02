package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/satishgonella2024/sentinelstacks/internal/tools"
)

// SearchTool implements a tool for web search
type SearchTool struct {
	tools.BaseTool
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

// SearchResult represents a search result
type SearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

// SearchResponse represents a response from the search API
type SearchResponse struct {
	Results     []SearchResult `json:"results"`
	TotalHits   int            `json:"total_hits"`
	SearchTime  float64        `json:"search_time"`
	Query       string         `json:"query"`
	Error       string         `json:"error,omitempty"`
}

// NewSearchTool creates a new web search tool
func NewSearchTool(apiKey, endpoint string) *SearchTool {
	if endpoint == "" {
		endpoint = "https://api.searchweb.example.com/v1/search"
	}
	
	return &SearchTool{
		BaseTool: tools.BaseTool{
			Name:        "web/search",
			Description: "Search for information on the web",
			Parameters: []tools.Parameter{
				{
					Name:        "query",
					Type:        "string",
					Description: "Search query",
					Required:    true,
				},
				{
					Name:        "num_results",
					Type:        "integer",
					Description: "Number of results to return (max 10)",
					Required:    false,
					Default:     5,
				},
				{
					Name:        "safe_search",
					Type:        "boolean",
					Description: "Whether to enable safe search",
					Required:    false,
					Default:     true,
				},
			},
			Permission: tools.PermissionNetwork,
		},
		apiKey:   apiKey,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute performs a web search
func (t *SearchTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get parameters
	query, _ := params["query"].(string)
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	
	// Get number of results
	numResults := 5 // default
	if numParam, ok := params["num_results"]; ok {
		switch v := numParam.(type) {
		case int:
			numResults = v
		case float64:
			numResults = int(v)
		}
		
		// Enforce limits
		if numResults < 1 {
			numResults = 1
		} else if numResults > 10 {
			numResults = 10
		}
	}
	
	// Get safe search parameter
	safeSearch := true // default
	if safeParam, ok := params["safe_search"].(bool); ok {
		safeSearch = safeParam
	}
	
	// Build search request
	return t.performSearch(ctx, query, numResults, safeSearch)
}

// performSearch executes the search request
func (t *SearchTool) performSearch(ctx context.Context, query string, numResults int, safeSearch bool) (*SearchResponse, error) {
	// If no API key, return mock results
	if t.apiKey == "" {
		return t.getMockSearchResults(query, numResults), nil
	}
	
	// Build request URL
	reqURL, err := url.Parse(t.endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid search endpoint: %w", err)
	}
	
	// Add query parameters
	q := reqURL.Query()
	q.Add("q", query)
	q.Add("num", fmt.Sprintf("%d", numResults))
	q.Add("safe", fmt.Sprintf("%t", safeSearch))
	reqURL.RawQuery = q.Encode()
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add API key header
	req.Header.Add("X-Api-Key", t.apiKey)
	req.Header.Add("Accept", "application/json")
	
	// Execute request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search API returned error: %s", resp.Status)
	}
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse response
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Check for API error
	if searchResp.Error != "" {
		return nil, fmt.Errorf("search API error: %s", searchResp.Error)
	}
	
	return &searchResp, nil
}

// getMockSearchResults returns mock search results for testing without an API key
func (t *SearchTool) getMockSearchResults(query string, numResults int) *SearchResponse {
	results := []SearchResult{
		{
			Title:       "Mock search result 1 for: " + query,
			URL:         "https://example.com/result1",
			Description: "This is a mock search result description for testing purposes. It contains information related to " + query + ".",
			Source:      "example.com",
		},
		{
			Title:       "Mock search result 2 for: " + query,
			URL:         "https://example.org/result2",
			Description: "Another mock search result with different content about " + query + ".",
			Source:      "example.org",
		},
		{
			Title:       "Mock search result 3 for: " + query,
			URL:         "https://example.net/result3",
			Description: "A third mock result discussing various aspects of " + query + ".",
			Source:      "example.net",
		},
		{
			Title:       "Mock search result 4 for: " + query,
			URL:         "https://example.io/result4",
			Description: "Yet another search result containing information about " + query + ".",
			Source:      "example.io",
		},
		{
			Title:       "Mock search result 5 for: " + query,
			URL:         "https://example.edu/result5",
			Description: "A fifth mock result with educational information about " + query + ".",
			Source:      "example.edu",
		},
	}
	
	// Limit results to requested number
	if numResults < len(results) {
		results = results[:numResults]
	}
	
	return &SearchResponse{
		Results:    results,
		TotalHits:  len(results),
		SearchTime: 0.05,
		Query:      query,
	}
}