// Package multimodal provides types and utilities for handling multimodal content
package multimodal

// MediaType represents the type of media content
type MediaType string

const (
	// MediaTypeText represents text content
	MediaTypeText MediaType = "text"
	// MediaTypeImage represents image content
	MediaTypeImage MediaType = "image"
	// MediaTypeAudio represents audio content
	MediaTypeAudio MediaType = "audio"
	// MediaTypeVideo represents video content
	MediaTypeVideo MediaType = "video"
)

// Content represents a piece of multimodal content
type Content struct {
	Type     MediaType `json:"type"`
	Text     string    `json:"text,omitempty"`
	Data     []byte    `json:"data,omitempty"`
	MimeType string    `json:"mime_type,omitempty"`
	URI      string    `json:"uri,omitempty"`
}

// NewTextContent creates a new text content
func NewTextContent(text string) *Content {
	return &Content{
		Type: MediaTypeText,
		Text: text,
	}
}

// NewImageContent creates a new image content
func NewImageContent(data []byte, mimeType string) *Content {
	return &Content{
		Type:     MediaTypeImage,
		Data:     data,
		MimeType: mimeType,
	}
}

// NewImageURIContent creates a new image content with a URI
func NewImageURIContent(uri string, mimeType string) *Content {
	return &Content{
		Type:     MediaTypeImage,
		URI:      uri,
		MimeType: mimeType,
	}
}

// Input represents multimodal input for LLM providers
type Input struct {
	Contents    []*Content             `json:"contents"`
	MaxTokens   int                    `json:"max_tokens"`
	Temperature float64                `json:"temperature"`
	Stream      bool                   `json:"stream"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewInput creates a new multimodal input
func NewInput() *Input {
	return &Input{
		Contents:    []*Content{},
		MaxTokens:   4096,
		Temperature: 0.7,
		Stream:      false,
		Metadata:    make(map[string]interface{}),
	}
}

// AddContent adds content to the input
func (i *Input) AddContent(content *Content) {
	i.Contents = append(i.Contents, content)
}

// AddText adds text content to the input
func (i *Input) AddText(text string) {
	i.AddContent(NewTextContent(text))
}

// AddImage adds image content to the input
func (i *Input) AddImage(data []byte, mimeType string) {
	i.AddContent(NewImageContent(data, mimeType))
}

// AddImageURI adds image content with a URI to the input
func (i *Input) AddImageURI(uri string, mimeType string) {
	i.AddContent(NewImageURIContent(uri, mimeType))
}

// SetMetadata sets a metadata value
func (i *Input) SetMetadata(key string, value interface{}) {
	i.Metadata[key] = value
}

// GetMetadata gets a metadata value
func (i *Input) GetMetadata(key string) (interface{}, bool) {
	value, ok := i.Metadata[key]
	return value, ok
}

// Output represents multimodal output from LLM providers
type Output struct {
	Contents []*Content             `json:"contents"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewOutput creates a new multimodal output
func NewOutput() *Output {
	return &Output{
		Contents: []*Content{},
		Metadata: make(map[string]interface{}),
	}
}

// AddContent adds content to the output
func (o *Output) AddContent(content *Content) {
	o.Contents = append(o.Contents, content)
}

// AddText adds text content to the output
func (o *Output) AddText(text string) {
	o.AddContent(NewTextContent(text))
}

// Chunk represents a chunk of a multimodal streaming response
type Chunk struct {
	Content  *Content               `json:"content"`
	IsFinal  bool                   `json:"is_final"`
	Error    error                  `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewChunk creates a new chunk
func NewChunk(content *Content, isFinal bool) *Chunk {
	return &Chunk{
		Content:  content,
		IsFinal:  isFinal,
		Metadata: make(map[string]interface{}),
	}
}
