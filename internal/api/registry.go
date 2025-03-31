package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
}

// RegistrySearchRequest represents a search request for the registry
type RegistrySearchRequest struct {
	Query  string   `json:"query"`
	Tags   []string `json:"tags,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
}

// RegistrySearchResponse represents a search response from the registry
type RegistrySearchResponse struct {
	Results []ImageInfo `json:"results"`
	Total   int         `json:"total"`
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
}

// PushImageRequest represents a request to push an image to the registry
type PushImageRequest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Tag         string            `json:"tag"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	File        string            `json:"file"` // Base64 encoded file
}

// PullImageRequest represents a request to pull an image from the registry
type PullImageRequest struct {
	Name string `json:"name"`
	Tag  string `json:"tag,omitempty"` // Optional, defaults to "latest"
}

// @Summary Search registry
// @Description Search for images in the registry
// @Tags registry
// @Accept json
// @Produce json
// @Param search body RegistrySearchRequest true "Search parameters"
// @Success 200 {object} RegistrySearchResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/search [get]
func (s *Server) searchRegistry(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limit := 10 // Default limit
	offset := 0 // Default offset

	// TODO: Implement real registry search
	// For now, return mock image data based on the query
	mockImages := []ImageInfo{
		{
			ID:        "sha256:abcdef1234567890",
			Name:      "user/chatbot",
			Tag:       "latest",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			Size:      1024 * 1024 * 5, // 5MB
			BaseModel: "claude-3-haiku-20240307",
			Features:  []string{"chat", "reasoning"},
		},
		{
			ID:        "sha256:9876543210abcdef",
			Name:      "user/research-assistant",
			Tag:       "v1.0",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			Size:      1024 * 1024 * 8, // 8MB
			BaseModel: "claude-3-opus-20240229",
			Features:  []string{"research", "summarization", "multimodal"},
		},
		{
			ID:        "sha256:fedcba0987654321",
			Name:      "user/translator",
			Tag:       "v2.1",
			CreatedAt: time.Now().Add(-120 * time.Hour),
			Size:      1024 * 1024 * 3, // 3MB
			BaseModel: "claude-3-sonnet-20240229",
			Features:  []string{"translation", "language-detection"},
		},
	}

	// Filter based on query if provided
	results := []ImageInfo{}
	if query != "" {
		for _, img := range mockImages {
			if strings.Contains(img.Name, query) || strings.Contains(img.BaseModel, query) {
				results = append(results, img)
			}
		}
	} else {
		results = mockImages
	}

	// Create response
	response := RegistrySearchResponse{
		Results: results,
		Total:   len(results),
		Offset:  offset,
		Limit:   limit,
	}

	s.sendJSON(w, http.StatusOK, response)
}

// @Summary Push image
// @Description Push an image to the registry
// @Tags registry
// @Accept json
// @Produce json
// @Param image body PushImageRequest true "Image to push"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/push [post]
func (s *Server) pushImage(w http.ResponseWriter, r *http.Request) {
	var req PushImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Name == "" || req.Tag == "" {
		s.sendError(w, http.StatusBadRequest, "Name and tag are required")
		return
	}

	// TODO: Implement real image push
	// For now, just return success

	s.sendJSON(w, http.StatusOK, map[string]string{
		"id":      "sha256:" + generateMockSHA256(),
		"name":    req.Name,
		"tag":     req.Tag,
		"status":  "pushed",
		"message": "Image pushed successfully",
	})
}

// @Summary Pull image
// @Description Pull an image from the registry
// @Tags registry
// @Accept json
// @Produce json
// @Param image body PullImageRequest true "Image to pull"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/pull [post]
func (s *Server) pullImage(w http.ResponseWriter, r *http.Request) {
	var req PullImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Name == "" {
		s.sendError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Use "latest" as default tag if not specified
	if req.Tag == "" {
		req.Tag = "latest"
	}

	// TODO: Implement real image pull
	// For now, just return success

	s.sendJSON(w, http.StatusOK, map[string]string{
		"id":      "sha256:" + generateMockSHA256(),
		"name":    req.Name,
		"tag":     req.Tag,
		"status":  "pulled",
		"message": "Image pulled successfully",
	})
}

// Helper function to generate a mock SHA256 hash
func generateMockSHA256() string {
	const chars = "0123456789abcdef"
	result := make([]byte, 64)
	for i := 0; i < 64; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
