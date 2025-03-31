package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ImageInfo contains information about a Sentinel Image
type ImageInfo struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Tag        string                 `json:"tag"`
	CreatedAt  time.Time              `json:"created_at"`
	Size       int64                  `json:"size"`
	LLM        string                 `json:"llm"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// listImagesHandler handles GET /images
func (s *Server) listImagesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual image listing from registry
	// For now, return dummy image data
	images := []ImageInfo{
		{
			ID:        "sha256:abcdef1234567890",
			Name:      "user/chatbot",
			Tag:       "latest",
			CreatedAt: time.Now().Add(-3 * time.Hour),
			Size:      1024 * 1024 * 5, // 5 MB
			LLM:       "claude-3-haiku-20240307",
			Parameters: map[string]interface{}{
				"temperature": 0.7,
				"memoryDepth": 10,
			},
		},
		{
			ID:        "sha256:9876543210abcdef",
			Name:      "user/research-assistant",
			Tag:       "v1.0",
			CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
			Size:      1024 * 1024 * 8, // 8 MB
			LLM:       "claude-3-opus-20240229",
			Parameters: map[string]interface{}{
				"temperature": 0.5,
				"memoryDepth": 20,
			},
		},
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"images": images,
	})
}

// getImageHandler handles GET /images/{id}
func (s *Server) getImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement actual image lookup from registry
	// For now, return dummy image data based on the requested ID

	// Return a 404 if ID doesn't match our dummy data
	if id != "sha256:abcdef1234567890" && id != "sha256:9876543210abcdef" {
		s.sendError(w, http.StatusNotFound, "Image not found")
		return
	}

	var image ImageInfo
	if id == "sha256:abcdef1234567890" {
		image = ImageInfo{
			ID:        "sha256:abcdef1234567890",
			Name:      "user/chatbot",
			Tag:       "latest",
			CreatedAt: time.Now().Add(-3 * time.Hour),
			Size:      1024 * 1024 * 5, // 5 MB
			LLM:       "claude-3-haiku-20240307",
			Parameters: map[string]interface{}{
				"temperature": 0.7,
				"memoryDepth": 10,
			},
		}
	} else {
		image = ImageInfo{
			ID:        "sha256:9876543210abcdef",
			Name:      "user/research-assistant",
			Tag:       "v1.0",
			CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
			Size:      1024 * 1024 * 8, // 8 MB
			LLM:       "claude-3-opus-20240229",
			Parameters: map[string]interface{}{
				"temperature": 0.5,
				"memoryDepth": 20,
			},
		}
	}

	// Add additional details for the detailed view
	details := map[string]interface{}{
		"id":         image.ID,
		"name":       image.Name,
		"tag":        image.Tag,
		"created_at": image.CreatedAt,
		"size":       image.Size,
		"llm":        image.LLM,
		"parameters": image.Parameters,
		"capabilities": []string{
			"web_search",
			"document_analysis",
		},
		"metadata": map[string]string{
			"description": fmt.Sprintf("An AI assistant based on %s", image.LLM),
			"author":      "user",
			"version":     image.Tag,
		},
	}

	s.sendJSON(w, http.StatusOK, details)
}
