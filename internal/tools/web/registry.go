package web

import (
	"github.com/spf13/viper"
	"github.com/satishgonella2024/sentinelstacks/internal/tools"
)

// RegisterWebTools registers all web tools
func RegisterWebTools() error {
	// Get search API key from config
	apiKey := viper.GetString("web.search_api_key")
	endpoint := viper.GetString("web.search_endpoint")
	
	// Create tools
	searchTool := NewSearchTool(apiKey, endpoint)
	
	// Register tools
	registry := tools.GetRegistry()
	if err := registry.RegisterTool(searchTool); err != nil {
		return err
	}
	
	return nil
}

// UpdateSearchAPIKey updates the API key for the search tool
func UpdateSearchAPIKey(apiKey string) error {
	// Get registry
	registry := tools.GetRegistry()
	
	// Update search tool
	searchTool, err := registry.GetTool("web/search")
	if err == nil {
		if webTool, ok := searchTool.(*SearchTool); ok {
			webTool.apiKey = apiKey
		}
	}
	
	return nil
}

// UpdateSearchEndpoint updates the endpoint for the search tool
func UpdateSearchEndpoint(endpoint string) error {
	// Get registry
	registry := tools.GetRegistry()
	
	// Update search tool
	searchTool, err := registry.GetTool("web/search")
	if err == nil {
		if webTool, ok := searchTool.(*SearchTool); ok {
			webTool.endpoint = endpoint
		}
	}
	
	return nil
}