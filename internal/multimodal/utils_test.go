package multimodal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultImageOptions(t *testing.T) {
	options := DefaultImageOptions()

	if options.MaxWidth != 1024 {
		t.Errorf("Expected MaxWidth to be 1024, got %d", options.MaxWidth)
	}

	if options.MaxHeight != 1024 {
		t.Errorf("Expected MaxHeight to be 1024, got %d", options.MaxHeight)
	}

	if !options.Compress {
		t.Errorf("Expected Compress to be true")
	}

	if options.Quality != 85 {
		t.Errorf("Expected Quality to be 85, got %d", options.Quality)
	}
}

func TestLoadImageFromFile(t *testing.T) {
	// Create a temporary image file
	tempDir := t.TempDir()
	imagePath := filepath.Join(tempDir, "test.png")

	// Create a simple test image
	img := createTestImage(100, 100)

	if err := createImageFile(img, imagePath); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test with default options
	content, err := LoadImageFromFile(imagePath, nil)
	if err != nil {
		t.Fatalf("LoadImageFromFile failed: %v", err)
	}

	if content.Type != MediaTypeImage {
		t.Errorf("Expected Type to be MediaTypeImage, got %v", content.Type)
	}

	if content.MimeType != "image/png" {
		t.Errorf("Expected MimeType to be image/png, got %s", content.MimeType)
	}

	if content.Text != "test.png" {
		t.Errorf("Expected Text to be test.png, got %s", content.Text)
	}

	// Test with custom options
	options := ImageProcessingOptions{
		MaxWidth:  50,
		MaxHeight: 50,
		Compress:  true,
		Quality:   75,
	}

	contentResized, err := LoadImageFromFile(imagePath, &options)
	if err != nil {
		t.Fatalf("LoadImageFromFile with options failed: %v", err)
	}

	// Decode the image to check its dimensions
	img2, _, err := image.Decode(bytes.NewReader(contentResized.Data))
	if err != nil {
		t.Fatalf("Failed to decode processed image: %v", err)
	}

	bounds := img2.Bounds()
	if bounds.Dx() > options.MaxWidth {
		t.Errorf("Expected width to be <= %d, got %d", options.MaxWidth, bounds.Dx())
	}

	if bounds.Dy() > options.MaxHeight {
		t.Errorf("Expected height to be <= %d, got %d", options.MaxHeight, bounds.Dy())
	}

	// Test with non-existent file
	_, err = LoadImageFromFile(filepath.Join(tempDir, "nonexistent.png"), nil)
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestLoadImageFromURL(t *testing.T) {
	// Create a test server that serves a test image
	testImage := createTestImage(100, 100)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, testImage)
	}))
	defer testServer.Close()

	// Test with default options
	content, err := LoadImageFromURL(testServer.URL, "Test Image", nil)
	if err != nil {
		t.Fatalf("LoadImageFromURL failed: %v", err)
	}

	if content.Type != MediaTypeImage {
		t.Errorf("Expected Type to be MediaTypeImage, got %v", content.Type)
	}

	if content.MimeType != "image/png" {
		t.Errorf("Expected MimeType to be image/png, got %s", content.MimeType)
	}

	if content.Text != "Test Image" {
		t.Errorf("Expected Text to be 'Test Image', got %s", content.Text)
	}

	if content.URL != testServer.URL {
		t.Errorf("Expected URL to be %s, got %s", testServer.URL, content.URL)
	}

	// Test with custom options
	options := ImageProcessingOptions{
		MaxWidth:  50,
		MaxHeight: 50,
		Compress:  true,
		Quality:   75,
	}

	contentResized, err := LoadImageFromURL(testServer.URL, "Test Image", &options)
	if err != nil {
		t.Fatalf("LoadImageFromURL with options failed: %v", err)
	}

	// Decode the image to check its dimensions
	img2, _, err := image.Decode(bytes.NewReader(contentResized.Data))
	if err != nil {
		t.Fatalf("Failed to decode processed image: %v", err)
	}

	bounds := img2.Bounds()
	if bounds.Dx() > options.MaxWidth {
		t.Errorf("Expected width to be <= %d, got %d", options.MaxWidth, bounds.Dx())
	}

	if bounds.Dy() > options.MaxHeight {
		t.Errorf("Expected height to be <= %d, got %d", options.MaxHeight, bounds.Dy())
	}

	// Test with invalid URL
	_, err = LoadImageFromURL("http://nonexistent.example.com", "", nil)
	if err == nil {
		t.Errorf("Expected error for invalid URL, got nil")
	}
}

func TestSaveImageToFile(t *testing.T) {
	// Create a test image content
	img := createTestImage(100, 100)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	content := &Content{
		Type:     MediaTypeImage,
		Data:     buf.Bytes(),
		MimeType: "image/png",
	}

	// Save to temporary file
	tempDir := t.TempDir()
	savePath := filepath.Join(tempDir, "saved.png")

	err := SaveImageToFile(content, savePath)
	if err != nil {
		t.Fatalf("SaveImageToFile failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", savePath)
	}

	// Test save with non-image content
	textContent := &Content{
		Type: MediaTypeText,
		Text: "This is text, not an image",
	}

	err = SaveImageToFile(textContent, savePath)
	if err == nil {
		t.Errorf("Expected error when saving non-image content, got nil")
	}

	// Test with empty data
	emptyContent := &Content{
		Type:     MediaTypeImage,
		Data:     nil,
		MimeType: "image/png",
	}

	err = SaveImageToFile(emptyContent, savePath)
	if err == nil {
		t.Errorf("Expected error when saving with empty data, got nil")
	}
}

func TestProcessImage(t *testing.T) {
	// Create a test image
	img := createTestImage(100, 100)
	var buf bytes.Buffer
	png.Encode(&buf, img)
	imageData := buf.Bytes()

	// Test resizing
	options := ImageProcessingOptions{
		MaxWidth:  50,
		MaxHeight: 50,
		Compress:  false,
	}

	content, err := processImage(imageData, "image/png", options)
	if err != nil {
		t.Fatalf("processImage failed: %v", err)
	}

	// Decode the image to check its dimensions
	img2, _, err := image.Decode(bytes.NewReader(content.Data))
	if err != nil {
		t.Fatalf("Failed to decode processed image: %v", err)
	}

	bounds := img2.Bounds()
	if bounds.Dx() > options.MaxWidth {
		t.Errorf("Expected width to be <= %d, got %d", options.MaxWidth, bounds.Dx())
	}

	if bounds.Dy() > options.MaxHeight {
		t.Errorf("Expected height to be <= %d, got %d", options.MaxHeight, bounds.Dy())
	}

	// Test compression
	options = ImageProcessingOptions{
		MaxWidth:  0,
		MaxHeight: 0,
		Compress:  true,
		Quality:   50,
	}

	compressedContent, err := processImage(imageData, "image/jpeg", options)
	if err != nil {
		t.Fatalf("processImage for compression failed: %v", err)
	}

	// Check if compression resulted in smaller file
	if compressedContent.MimeType != "image/jpeg" {
		t.Errorf("Expected compressed image MIME type to be image/jpeg, got %s", compressedContent.MimeType)
	}
}

func TestBuildMultimodalInput(t *testing.T) {
	// Create a test image content
	img := createTestImage(100, 100)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	imageContent := &Content{
		Type:     MediaTypeImage,
		Data:     buf.Bytes(),
		MimeType: "image/png",
	}

	// Test with text and image
	opts := map[string]interface{}{
		"max_tokens":    100,
		"temperature":   0.7,
		"stream":        true,
		"system_prompt": "You are a helpful assistant",
	}

	input := BuildMultimodalInput("Analyze this image", imageContent, opts)

	if len(input.Contents) != 2 {
		t.Errorf("Expected 2 contents, got %d", len(input.Contents))
	}

	if input.Contents[0].Type != MediaTypeText || input.Contents[0].Text != "Analyze this image" {
		t.Errorf("Expected first content to be text 'Analyze this image'")
	}

	if input.Contents[1].Type != MediaTypeImage || !bytes.Equal(input.Contents[1].Data, imageContent.Data) {
		t.Errorf("Expected second content to be the test image")
	}

	if input.MaxTokens != 100 {
		t.Errorf("Expected MaxTokens to be 100, got %d", input.MaxTokens)
	}

	if input.Temperature != 0.7 {
		t.Errorf("Expected Temperature to be 0.7, got %f", input.Temperature)
	}

	if !input.Stream {
		t.Errorf("Expected Stream to be true")
	}

	if v, ok := input.Metadata["system_prompt"]; !ok || v != "You are a helpful assistant" {
		t.Errorf("Expected system_prompt metadata to be set correctly")
	}

	// Test with text only
	textInput := BuildMultimodalInput("Text only input", nil, opts)

	if len(textInput.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(textInput.Contents))
	}

	if textInput.Contents[0].Type != MediaTypeText || textInput.Contents[0].Text != "Text only input" {
		t.Errorf("Expected content to be text 'Text only input'")
	}
}

func TestExtractTextFromOutput(t *testing.T) {
	output := NewOutput()
	output.AddText("Hello, ")
	output.AddText("world!")

	text := ExtractTextFromOutput(output)
	if text != "Hello, world!" {
		t.Errorf("Expected 'Hello, world!', got '%s'", text)
	}

	// Test with nil output
	nilText := ExtractTextFromOutput(nil)
	if nilText != "" {
		t.Errorf("Expected empty string for nil output, got '%s'", nilText)
	}
}

func TestExtractImagesFromOutput(t *testing.T) {
	// Create a test image content
	img := createTestImage(100, 100)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	imageData := buf.Bytes()

	output := NewOutput()
	output.AddText("This is text")
	output.AddImage(imageData, "image/png")
	output.AddText("More text")

	images := ExtractImagesFromOutput(output)

	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}

	if images[0].Type != MediaTypeImage || !bytes.Equal(images[0].Data, imageData) {
		t.Errorf("Extracted image doesn't match original")
	}

	// Test with nil output
	nilImages := ExtractImagesFromOutput(nil)
	if nilImages != nil {
		t.Errorf("Expected nil for nil output, got %v", nilImages)
	}

	// Test with no images
	textOutput := NewOutput().AddText("Text only")
	noImages := ExtractImagesFromOutput(textOutput)
	if len(noImages) != 0 {
		t.Errorf("Expected empty slice for output with no images, got %d items", len(noImages))
	}
}

// Helper functions

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a simple pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}

	return img
}

func createImageFile(img image.Image, filePath string) error {
	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode image
	if strings.HasSuffix(filePath, ".png") {
		return png.Encode(f, img)
	}

	return nil
}
