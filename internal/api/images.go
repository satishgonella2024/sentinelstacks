package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ImageInfo represents information about an agent image
type ImageInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
	BaseModel string    `json:"base_model"`
	Features  []string  `json:"features"`
}

// ImageResponse represents the response for image endpoints
type ImageResponse struct {
	Images []ImageInfo `json:"images"`
}

// @Summary List all images
// @Description Get a list of all available agent images
// @Tags images
// @Accept json
// @Produce json
// @Success 200 {object} ImageResponse
// @Failure 500 {object} map[string]string
// @Router /images [get]
func (s *Server) listImagesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement real image listing
	// For now, return mock image data
	images := []ImageInfo{
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

	s.sendJSON(w, http.StatusOK, ImageResponse{
		Images: images,
	})
}

// @Summary Get image details
// @Description Get details of a specific image by ID
// @Tags images
// @Accept json
// @Produce json
// @Param id path string true "Image ID"
// @Success 200 {object} ImageInfo
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /images/{id} [get]
func (s *Server) getImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement real image retrieval
	// For now, return mock data based on the requested ID
	if id == "sha256:abcdef1234567890" {
		image := ImageInfo{
			ID:        "sha256:abcdef1234567890",
			Name:      "user/chatbot",
			Tag:       "latest",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			Size:      1024 * 1024 * 5, // 5MB
			BaseModel: "claude-3-haiku-20240307",
			Features:  []string{"chat", "reasoning"},
		}
		s.sendJSON(w, http.StatusOK, image)
		return
	} else if id == "sha256:9876543210abcdef" {
		image := ImageInfo{
			ID:        "sha256:9876543210abcdef",
			Name:      "user/research-assistant",
			Tag:       "v1.0",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			Size:      1024 * 1024 * 8, // 8MB
			BaseModel: "claude-3-opus-20240229",
			Features:  []string{"research", "summarization", "multimodal"},
		}
		s.sendJSON(w, http.StatusOK, image)
		return
	} else if id == "sha256:fedcba0987654321" {
		image := ImageInfo{
			ID:        "sha256:fedcba0987654321",
			Name:      "user/translator",
			Tag:       "v2.1",
			CreatedAt: time.Now().Add(-120 * time.Hour),
			Size:      1024 * 1024 * 3, // 3MB
			BaseModel: "claude-3-sonnet-20240229",
			Features:  []string{"translation", "language-detection"},
		}
		s.sendJSON(w, http.StatusOK, image)
		return
	}

	s.sendError(w, http.StatusNotFound, "Image not found")
}
