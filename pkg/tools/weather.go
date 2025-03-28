package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// WeatherTool fetches current weather information for a location
type WeatherTool struct {
	ApiKey string
}

// NewWeatherTool creates a new weather tool with the given API key
func NewWeatherTool(apiKey string) *WeatherTool {
	return &WeatherTool{
		ApiKey: apiKey,
	}
}

// ID returns the unique identifier for the weather tool
func (w *WeatherTool) ID() string {
	return "weather"
}

// Name returns a user-friendly name
func (w *WeatherTool) Name() string {
	return "Weather"
}

// Description returns a detailed description
func (w *WeatherTool) Description() string {
	return "Fetches current weather information for a city or location"
}

// Version returns the semantic version
func (w *WeatherTool) Version() string {
	return "0.1.0"
}

// ParameterSchema returns the JSON schema for parameters
func (w *WeatherTool) ParameterSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type": "string",
				"description": "The city or location to get weather for (e.g., 'London', 'New York')",
			},
			"units": map[string]interface{}{
				"type": "string",
				"enum": []string{"metric", "imperial"},
				"default": "metric",
				"description": "Units of measurement (metric: Celsius, imperial: Fahrenheit)",
			},
		},
		"required": []string{"location"},
	}
}

// Execute runs the weather tool with the provided parameters
func (w *WeatherTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Check if we have an API key
	if w.ApiKey == "" {
		return nil, fmt.Errorf("weather API key is not set")
	}

	// Extract parameters
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("location parameter is required")
	}

	units := "metric"
	if unitParam, ok := params["units"].(string); ok && unitParam != "" {
		units = unitParam
	}

	// Build the API URL (using OpenWeatherMap API as an example)
	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	queryParams := url.Values{}
	queryParams.Add("q", location)
	queryParams.Add("units", units)
	queryParams.Add("appid", w.ApiKey)

	// Make the request
	resp, err := http.Get(baseURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API error: %s", string(body))
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

	// Return formatted output
	return formatWeatherResult(result, units), nil
}

// formatWeatherResult creates a simplified weather summary
func formatWeatherResult(data map[string]interface{}, units string) map[string]interface{} {
	result := map[string]interface{}{
		"location": getNestedString(data, "name"),
		"country":  getNestedString(data, "sys", "country"),
		"weather":  getNestedString(data, "weather", 0, "main"),
		"description": getNestedString(data, "weather", 0, "description"),
	}

	// Extract temperature
	if main, ok := data["main"].(map[string]interface{}); ok {
		if temp, ok := main["temp"].(float64); ok {
			result["temperature"] = temp
			if units == "metric" {
				result["temperature_unit"] = "°C"
			} else {
				result["temperature_unit"] = "°F"
			}
		}
		if humidity, ok := main["humidity"].(float64); ok {
			result["humidity"] = humidity
		}
	}

	// Extract wind
	if wind, ok := data["wind"].(map[string]interface{}); ok {
		if speed, ok := wind["speed"].(float64); ok {
			result["wind_speed"] = speed
			if units == "metric" {
				result["wind_unit"] = "m/s"
			} else {
				result["wind_unit"] = "mph"
			}
		}
	}

	return result
}

// getNestedString safely navigates a nested structure to extract a string value
func getNestedString(data map[string]interface{}, keys ...interface{}) string {
	var current interface{} = data
	
	for _, key := range keys {
		switch k := key.(type) {
		case string:
			if m, ok := current.(map[string]interface{}); ok {
				current = m[k]
			} else {
				return ""
			}
		case int:
			if arr, ok := current.([]interface{}); ok && k < len(arr) {
				current = arr[k]
			} else {
				return ""
			}
		}
	}
	
	if str, ok := current.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", current)
}
