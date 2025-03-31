// Package multimodal provides support for multimodal content in SentinelStacks
package multimodal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

// MediaType represents the type of media content
type MediaType string

// Media type constants
const (
	MediaTypeText  MediaType = "text"
	MediaTypeImage MediaType = "image"
	MediaTypeAudio MediaType = "audio"
	MediaTypeVideo MediaType = "video"
)

// Content represents a piece of multimodal content
type Content struct {
	Type     MediaType              `json:"type"`
	Data     []byte                 `json:"-"` // Raw binary data (not serialized to JSON)
	MimeType string                 `json:"mime_type,omitempty"`
	Text     string                 `json:"text,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewTextContent creates a new text content
func NewTextContent(text string) *Content {
	return &Content{
		Type: MediaTypeText,
		Text: text,
	}
}

// NewImageContent creates a new image content from raw bytes
func NewImageContent(data []byte, mimeType string) *Content {
	return &Content{
		Type:     MediaTypeImage,
		Data:     data,
		MimeType: mimeType,
	}
}

// NewImageContentFromURL creates a new image content from a URL
func NewImageContentFromURL(url string, altText string) *Content {
	return &Content{
		Type:     MediaTypeImage,
		URL:      url,
		Text:     altText,
		MimeType: guessMimeTypeFromURL(url),
	}
}

// ToBase64 returns the base64 encoded data for the content
func (c *Content) ToBase64() string {
	if c.Data == nil || len(c.Data) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(c.Data)
}

// ToDataURL returns a data URL for the content
func (c *Content) ToDataURL() string {
	if c.Data == nil || len(c.Data) == 0 {
		return ""
	}
	return fmt.Sprintf("data:%s;base64,%s", c.MimeType, c.ToBase64())
}

// ResizeImage resizes the image content to the specified dimensions
// while maintaining aspect ratio
func (c *Content) ResizeImage(maxWidth, maxHeight int) error {
	if c.Type != MediaTypeImage || c.Data == nil {
		return fmt.Errorf("not an image content or no data")
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(c.Data))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Get current dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resize is needed
	if width <= maxWidth && height <= maxHeight {
		return nil // No resize needed
	}

	// Calculate new dimensions maintaining aspect ratio
	ratio := float64(width) / float64(height)
	if width > height {
		width = maxWidth
		height = int(float64(width) / ratio)
		if height > maxHeight {
			height = maxHeight
			width = int(float64(height) * ratio)
		}
	} else {
		height = maxHeight
		width = int(float64(height) * ratio)
		if width > maxWidth {
			width = maxWidth
			height = int(float64(width) / ratio)
		}
	}

	// Resize the image (simplified, in a real implementation we would use a proper resizing library)
	// For now, we'll just re-encode at full size as a placeholder
	var buf bytes.Buffer
	switch c.MimeType {
	case "image/jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return fmt.Errorf("failed to encode jpeg: %w", err)
		}
	case "image/png":
		if err := png.Encode(&buf, img); err != nil {
			return fmt.Errorf("failed to encode png: %w", err)
		}
	default:
		return fmt.Errorf("unsupported image format for resizing: %s", c.MimeType)
	}

	c.Data = buf.Bytes()
	return nil
}

// guessMimeTypeFromURL attempts to guess the MIME type from a URL
func guessMimeTypeFromURL(url string) string {
	lower := strings.ToLower(url)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(lower, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(lower, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(lower, ".webp") {
		return "image/webp"
	}
	if strings.HasSuffix(lower, ".mp4") {
		return "video/mp4"
	}
	if strings.HasSuffix(lower, ".mp3") {
		return "audio/mpeg"
	}
	// Default to octet-stream for unknown types
	return "application/octet-stream"
}

// Input represents a multimodal input for generation
type Input struct {
	Contents    []*Content             `json:"contents"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewInput creates a new multimodal input
func NewInput() *Input {
	return &Input{
		Contents: make([]*Content, 0),
		Metadata: make(map[string]interface{}),
	}
}

// AddText adds text content to the input
func (i *Input) AddText(text string) *Input {
	i.Contents = append(i.Contents, NewTextContent(text))
	return i
}

// AddImage adds image content to the input
func (i *Input) AddImage(data []byte, mimeType string) *Input {
	i.Contents = append(i.Contents, NewImageContent(data, mimeType))
	return i
}

// AddImageURL adds an image URL to the input
func (i *Input) AddImageURL(url string, altText string) *Input {
	i.Contents = append(i.Contents, NewImageContentFromURL(url, altText))
	return i
}

// SetMaxTokens sets the maximum number of tokens to generate
func (i *Input) SetMaxTokens(maxTokens int) *Input {
	i.MaxTokens = maxTokens
	return i
}

// SetTemperature sets the temperature for generation
func (i *Input) SetTemperature(temperature float64) *Input {
	i.Temperature = temperature
	return i
}

// SetStream sets whether to stream the response
func (i *Input) SetStream(stream bool) *Input {
	i.Stream = stream
	return i
}

// SetMetadata sets metadata for the input
func (i *Input) SetMetadata(key string, value interface{}) *Input {
	if i.Metadata == nil {
		i.Metadata = make(map[string]interface{})
	}
	i.Metadata[key] = value
	return i
}

// Output represents a multimodal output from generation
type Output struct {
	Contents   []*Content             `json:"contents"`
	UsedTokens int                    `json:"used_tokens,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NewOutput creates a new multimodal output
func NewOutput() *Output {
	return &Output{
		Contents: make([]*Content, 0),
		Metadata: make(map[string]interface{}),
	}
}

// AddText adds text content to the output
func (o *Output) AddText(text string) *Output {
	o.Contents = append(o.Contents, NewTextContent(text))
	return o
}

// AddImage adds image content to the output
func (o *Output) AddImage(data []byte, mimeType string) *Output {
	o.Contents = append(o.Contents, NewImageContent(data, mimeType))
	return o
}

// GetText returns the concatenated text from all text contents
func (o *Output) GetText() string {
	var result strings.Builder
	for _, content := range o.Contents {
		if content.Type == MediaTypeText {
			result.WriteString(content.Text)
		}
	}
	return result.String()
}

// GetFirstImage returns the first image content, if any
func (o *Output) GetFirstImage() *Content {
	for _, content := range o.Contents {
		if content.Type == MediaTypeImage {
			return content
		}
	}
	return nil
}

// Chunk represents a chunk of multimodal content in a streaming response
type Chunk struct {
	Content *Content `json:"content"`
	IsFinal bool     `json:"is_final"`
	Error   error    `json:"error,omitempty"`
}

// ContentReader provides a reader interface for multimodal content
type ContentReader interface {
	io.Reader
	ContentType() string
}

// contentReader implements ContentReader
type contentReader struct {
	reader io.Reader
	ctype  string
}

// NewContentReader creates a new ContentReader from a Content
func NewContentReader(content *Content) ContentReader {
	return &contentReader{
		reader: bytes.NewReader(content.Data),
		ctype:  content.MimeType,
	}
}

// Read implements io.Reader
func (r *contentReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// ContentType returns the MIME type of the content
func (r *contentReader) ContentType() string {
	return r.ctype
}
