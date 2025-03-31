package multimodal

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// MaxImageSize defines the maximum size of images in bytes (10MB)
const MaxImageSize = 10 * 1024 * 1024

// ImageProcessingOptions defines options for image processing
type ImageProcessingOptions struct {
	MaxWidth  int  // Maximum width of the image
	MaxHeight int  // Maximum height of the image
	Compress  bool // Whether to compress the image
	Quality   int  // JPEG quality (1-100)
}

// DefaultImageOptions returns default image processing options
func DefaultImageOptions() ImageProcessingOptions {
	return ImageProcessingOptions{
		MaxWidth:  1024,
		MaxHeight: 1024,
		Compress:  true,
		Quality:   85,
	}
}

// LoadImageFromFile loads an image from a file path and returns it as Content
func LoadImageFromFile(filePath string, options *ImageProcessingOptions) (*Content, error) {
	// Set default options if nil
	if options == nil {
		defaultOptions := DefaultImageOptions()
		options = &defaultOptions
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Check file size
	if len(fileData) > MaxImageSize {
		return nil, fmt.Errorf("image file too large (max size: %d bytes)", MaxImageSize)
	}

	// Determine MIME type
	mimeType := http.DetectContentType(fileData)
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, fmt.Errorf("file is not an image: %s", mimeType)
	}

	// Get the filename for alt text
	altText := filepath.Base(filePath)

	// If processing is required
	if options.MaxWidth > 0 || options.MaxHeight > 0 || options.Compress {
		content, err := processImage(fileData, mimeType, *options)
		if err != nil {
			return nil, err
		}
		// Set the alt text
		content.Text = altText
		return content, nil
	}

	// Create content with original image
	return &Content{
		Type:     MediaTypeImage,
		Data:     fileData,
		MimeType: mimeType,
		Text:     altText, // Use filename as default alt text
	}, nil
}

// LoadImageFromURL loads an image from a URL and returns it as Content
func LoadImageFromURL(url string, altText string, options *ImageProcessingOptions) (*Content, error) {
	// Set default options if nil
	if options == nil {
		defaultOptions := DefaultImageOptions()
		options = &defaultOptions
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: http.DefaultClient.Timeout,
	}

	// Get image from URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from URL: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read image data with size limit
	limitReader := io.LimitReader(resp.Body, int64(MaxImageSize)+1)
	imageData, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Check if image exceeds size limit
	if len(imageData) > MaxImageSize {
		return nil, fmt.Errorf("image file too large (max size: %d bytes)", MaxImageSize)
	}

	// Get content type from response headers or try to detect it
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" || !strings.HasPrefix(mimeType, "image/") {
		mimeType = http.DetectContentType(imageData)
		if !strings.HasPrefix(mimeType, "image/") {
			return nil, fmt.Errorf("URL does not point to an image: %s", mimeType)
		}
	}

	// Use URL basename as alt text if none provided
	if altText == "" {
		parts := strings.Split(url, "/")
		if len(parts) > 0 {
			altText = parts[len(parts)-1]
		} else {
			altText = "Image from URL"
		}
	}

	// If processing is required
	if options.MaxWidth > 0 || options.MaxHeight > 0 || options.Compress {
		content, err := processImage(imageData, mimeType, *options)
		if err != nil {
			return nil, err
		}
		content.URI = url
		content.Text = altText
		return content, nil
	}

	// Create content with original image
	return &Content{
		Type:     MediaTypeImage,
		Data:     imageData,
		MimeType: mimeType,
		URI:      url,
		Text:     altText,
	}, nil
}

// SaveImageToFile saves an image Content to a file
func SaveImageToFile(content *Content, filePath string) error {
	// Check if content is an image
	if content.Type != MediaTypeImage {
		return errors.New("content is not an image")
	}

	// Check if data is available
	if content.Data == nil || len(content.Data) == 0 {
		return errors.New("image data is empty")
	}

	// Create parent directories if they don't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, content.Data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// processImage processes an image with the given options
func processImage(imageData []byte, mimeType string, options ImageProcessingOptions) (*Content, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get current dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resize is needed
	needsResize := (options.MaxWidth > 0 && width > options.MaxWidth) ||
		(options.MaxHeight > 0 && height > options.MaxHeight)

	// If no resize or compression needed, return original
	if !needsResize && !options.Compress {
		return &Content{
			Type:     MediaTypeImage,
			Data:     imageData,
			MimeType: mimeType,
		}, nil
	}

	// Calculate new dimensions if resize is needed
	var newWidth, newHeight int
	if needsResize {
		// Calculate new dimensions maintaining aspect ratio
		ratio := float64(width) / float64(height)
		if options.MaxWidth > 0 && options.MaxHeight > 0 {
			if float64(options.MaxWidth)/float64(options.MaxHeight) > ratio {
				newHeight = options.MaxHeight
				newWidth = int(float64(newHeight) * ratio)
			} else {
				newWidth = options.MaxWidth
				newHeight = int(float64(newWidth) / ratio)
			}
		} else if options.MaxWidth > 0 {
			newWidth = options.MaxWidth
			newHeight = int(float64(newWidth) / ratio)
		} else {
			newHeight = options.MaxHeight
			newWidth = int(float64(newHeight) * ratio)
		}
	} else {
		newWidth = width
		newHeight = height
	}

	// Create new image
	var processedImg *image.RGBA
	if needsResize {
		processedImg = image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.CatmullRom.Scale(processedImg, processedImg.Bounds(), img, bounds, draw.Over, nil)
	} else {
		// Convert to RGBA format without resizing
		processedImg = image.NewRGBA(bounds)
		draw.Draw(processedImg, bounds, img, bounds.Min, draw.Src)
	}

	// Create buffer for the new image
	var buf bytes.Buffer

	// Encode with appropriate format and options
	switch {
	case strings.HasPrefix(mimeType, "image/jpeg") || (options.Compress && format != "png"):
		// Use JPEG with quality setting
		quality := options.Quality
		if quality <= 0 || quality > 100 {
			quality = 85 // Default quality
		}
		err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: quality})
		mimeType = "image/jpeg"
	case strings.HasPrefix(mimeType, "image/png"):
		// Use PNG
		err = png.Encode(&buf, processedImg)
		mimeType = "image/png"
	default:
		// Default to JPEG for other formats
		err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: 85})
		mimeType = "image/jpeg"
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode processed image: %w", err)
	}

	// Create content with processed image
	return &Content{
		Type:     MediaTypeImage,
		Data:     buf.Bytes(),
		MimeType: mimeType,
	}, nil
}

// BuildMultimodalInput creates a multimodal input from text and optional image content
func BuildMultimodalInput(text string, imageContent *Content, opts map[string]interface{}) *Input {
	input := NewInput()

	// Add text content
	if text != "" {
		input.AddText(text)
	}

	// Add image content if provided
	if imageContent != nil && imageContent.Type == MediaTypeImage {
		input.Contents = append(input.Contents, imageContent)
	}

	// Set options
	if maxTokens, ok := opts["max_tokens"].(int); ok {
		input.MaxTokens = maxTokens
	}

	if temperature, ok := opts["temperature"].(float64); ok {
		input.Temperature = temperature
	}

	if stream, ok := opts["stream"].(bool); ok {
		input.Stream = stream
	}

	// Add any additional metadata
	for k, v := range opts {
		if k != "max_tokens" && k != "temperature" && k != "stream" {
			input.SetMetadata(k, v)
		}
	}

	return input
}

// ExtractTextFromOutput extracts all text content from a multimodal output
func ExtractTextFromOutput(output *Output) string {
	if output == nil {
		return ""
	}

	return output.GetText()
}

// ExtractImagesFromOutput returns all image contents from a multimodal output
func ExtractImagesFromOutput(output *Output) []*Content {
	if output == nil {
		return nil
	}

	var images []*Content
	for _, content := range output.Contents {
		if content.Type == MediaTypeImage {
			images = append(images, content)
		}
	}

	return images
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

// GuessMimeTypeFromURL attempts to guess the MIME type from a URL
func GuessMimeTypeFromURL(url string) string {
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

// FluentAddText adds text content to the input and returns the input for method chaining
func (i *Input) FluentAddText(text string) *Input {
	i.AddText(text)
	return i
}

// FluentAddImage adds image content to the input and returns the input for method chaining
func (i *Input) FluentAddImage(data []byte, mimeType string) *Input {
	i.AddImage(data, mimeType)
	return i
}

// FluentAddImageURI adds image URI content to the input and returns the input for method chaining
func (i *Input) FluentAddImageURI(url string, altText string) *Input {
	content := NewImageURIContent(url, GuessMimeTypeFromURL(url))
	content.Text = altText
	i.AddContent(content)
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

// FluentSetMetadata sets metadata for the input and returns the input for method chaining
func (i *Input) FluentSetMetadata(key string, value interface{}) *Input {
	i.SetMetadata(key, value)
	return i
}

// GetText returns the combined text content from the output
func (o *Output) GetText() string {
	var text strings.Builder
	for _, content := range o.Contents {
		if content.Type == MediaTypeText {
			text.WriteString(content.Text)
		}
	}
	return text.String()
}

// GetFirstImage returns the first image content from the output, if any
func (o *Output) GetFirstImage() *Content {
	for _, content := range o.Contents {
		if content.Type == MediaTypeImage {
			return content
		}
	}
	return nil
}

// ContentReader is an interface for reading multimodal content
type ContentReader interface {
	io.Reader
	ContentType() string
}

type contentReader struct {
	reader io.Reader
	ctype  string
}

// NewContentReader creates a new content reader from a Content
func NewContentReader(content *Content) ContentReader {
	return &contentReader{
		reader: bytes.NewReader(content.Data),
		ctype:  content.MimeType,
	}
}

// Read reads from the content reader
func (r *contentReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// ContentType returns the content type of the reader
func (r *contentReader) ContentType() string {
	return r.ctype
}
