package registry

import (
	"github.com/satishgonella2024/sentinelstacks/pkg/agent"
)

// Image represents a Sentinel Image
type Image struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Tag          string                 `json:"tag"`
	CreatedAt    int64                  `json:"createdAt"`
	Definition   ImageDefinition        `json:"definition"`
	Dependencies []string               `json:"dependencies,omitempty"`
}

// ImageDefinition represents the definition of a Sentinel Image
type ImageDefinition struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	BaseModel    string                 `json:"baseModel"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Tools        []string               `json:"tools,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// ConvertFromAgentImage converts from an agent.Image to a registry.Image
func ConvertFromAgentImage(image *agent.Image) *Image {
	return &Image{
		ID:           image.ID,
		Name:         image.Name,
		Tag:          image.Tag,
		CreatedAt:    image.CreatedAt,
		Definition: ImageDefinition{
			Name:         image.Definition.Name,
			Description:  image.Definition.Description,
			BaseModel:    image.Definition.BaseModel,
			Capabilities: image.Definition.Capabilities,
			Tools:        image.Definition.Tools,
			Parameters:   image.Definition.Parameters,
		},
		Dependencies: image.Dependencies,
	}
}

// ConvertToAgentImage converts from a registry.Image to an agent.Image
func ConvertToAgentImage(image *Image) *agent.Image {
	return &agent.Image{
		ID:           image.ID,
		Name:         image.Name,
		Tag:          image.Tag,
		CreatedAt:    image.CreatedAt,
		Definition: agent.Definition{
			Name:         image.Definition.Name,
			Description:  image.Definition.Description,
			BaseModel:    image.Definition.BaseModel,
			Capabilities: image.Definition.Capabilities,
			Tools:        image.Definition.Tools,
			Parameters:   image.Definition.Parameters,
		},
		Dependencies: image.Dependencies,
	}
}
