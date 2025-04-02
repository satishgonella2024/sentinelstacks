package registry

import (
	"encoding/json"
	"fmt"
	"time"
)

// ImageSummary is an alias for ImageInfo for backward compatibility
type ImageSummary struct {
	Name        string    `json:"name"`
	Tag         string    `json:"tag"`
	Created     time.Time `json:"created"`
	Size        int64     `json:"size"`
	BaseModel   string    `json:"baseModel"`
	Description string    `json:"description"`
}

// FormatImagesAsJSON formats a list of ImageSummary as JSON
func FormatImagesAsJSON(images []ImageSummary) (string, error) {
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal images to JSON: %w", err)
	}
	
	return string(data), nil
}

// ConvertToImageSummary converts an ImageInfo to an ImageSummary
func ConvertToImageSummary(info ImageInfo) ImageSummary {
	return ImageSummary{
		Name:      info.Name,
		Tag:       info.Tag,
		Created:   info.CreatedAt,
		Size:      info.Size,
		BaseModel: info.BaseModel,
		// Description might not be available in ImageInfo
		Description: "",
	}
}

// List returns all images in the local registry as ImageSummary
func (r *LocalRegistry) List() ([]ImageSummary, error) {
	infoList, err := r.ListImageInfo()
	if err != nil {
		return nil, err
	}
	
	var summaries []ImageSummary
	for _, info := range infoList {
		summaries = append(summaries, ConvertToImageSummary(info))
	}
	
	return summaries, nil
}

