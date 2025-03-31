package multimodal

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"
)

func TestNewTextContent(t *testing.T) {
	text := "Hello, world!"
	content := NewTextContent(text)

	if content.Type != MediaTypeText {
		t.Errorf("Expected content type to be %s, got %s", MediaTypeText, content.Type)
	}

	if content.Text != text {
		t.Errorf("Expected text to be %s, got %s", text, content.Text)
	}
}

func TestNewImageContent(t *testing.T) {
	data := []byte("fake-image-data")
	mimeType := "image/jpeg"
	content := NewImageContent(data, mimeType)

	if content.Type != MediaTypeImage {
		t.Errorf("Expected content type to be %s, got %s", MediaTypeImage, content.Type)
	}

	if content.MimeType != mimeType {
		t.Errorf("Expected MIME type to be %s, got %s", mimeType, content.MimeType)
	}

	if !bytes.Equal(content.Data, data) {
		t.Errorf("Expected data to match the input data")
	}
}

func TestNewImageContentFromURL(t *testing.T) {
	url := "https://example.com/image.jpg"
	altText := "Example image"
	content := NewImageContentFromURL(url, altText)

	if content.Type != MediaTypeImage {
		t.Errorf("Expected content type to be %s, got %s", MediaTypeImage, content.Type)
	}

	if content.URL != url {
		t.Errorf("Expected URL to be %s, got %s", url, content.URL)
	}

	if content.Text != altText {
		t.Errorf("Expected alt text to be %s, got %s", altText, content.Text)
	}

	if content.MimeType != "image/jpeg" {
		t.Errorf("Expected MIME type to be image/jpeg, got %s", content.MimeType)
	}
}

func TestGuessMimeTypeFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg", "image/jpeg"},
		{"https://example.com/image.jpeg", "image/jpeg"},
		{"https://example.com/image.png", "image/png"},
		{"https://example.com/image.gif", "image/gif"},
		{"https://example.com/image.webp", "image/webp"},
		{"https://example.com/video.mp4", "video/mp4"},
		{"https://example.com/audio.mp3", "audio/mpeg"},
		{"https://example.com/unknown", "application/octet-stream"},
	}

	for _, test := range tests {
		result := guessMimeTypeFromURL(test.url)
		if result != test.expected {
			t.Errorf("For URL %s, expected MIME type %s, got %s", test.url, test.expected, result)
		}
	}
}

func TestContentToBase64(t *testing.T) {
	data := []byte("test-data")
	content := NewImageContent(data, "image/jpeg")

	expected := base64.StdEncoding.EncodeToString(data)
	result := content.ToBase64()

	if result != expected {
		t.Errorf("Expected base64 %s, got %s", expected, result)
	}

	// Test with empty data
	emptyContent := NewImageContent(nil, "image/jpeg")
	if emptyContent.ToBase64() != "" {
		t.Errorf("Expected empty base64 string for nil data")
	}
}

func TestContentToDataURL(t *testing.T) {
	data := []byte("test-data")
	mimeType := "image/jpeg"
	content := NewImageContent(data, mimeType)

	expected := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
	result := content.ToDataURL()

	if result != expected {
		t.Errorf("Expected data URL %s, got %s", expected, result)
	}

	// Test with empty data
	emptyContent := NewImageContent(nil, mimeType)
	if emptyContent.ToDataURL() != "" {
		t.Errorf("Expected empty data URL string for nil data")
	}
}

func TestMultimodalInput(t *testing.T) {
	input := NewInput()
	text := "Hello, world!"
	imageData := []byte("fake-image-data")
	imageMimeType := "image/png"
	imageURL := "https://example.com/image.jpg"
	altText := "Example image"

	// Test fluent API
	input.AddText(text).
		AddImage(imageData, imageMimeType).
		AddImageURL(imageURL, altText).
		SetMaxTokens(100).
		SetTemperature(0.7).
		SetStream(true).
		SetMetadata("system_prompt", "You are a helpful assistant.")

	// Check contents
	if len(input.Contents) != 3 {
		t.Errorf("Expected 3 contents, got %d", len(input.Contents))
	}

	if input.Contents[0].Type != MediaTypeText || input.Contents[0].Text != text {
		t.Errorf("Text content not added correctly")
	}

	if input.Contents[1].Type != MediaTypeImage || !bytes.Equal(input.Contents[1].Data, imageData) {
		t.Errorf("Image content not added correctly")
	}

	if input.Contents[2].Type != MediaTypeImage || input.Contents[2].URL != imageURL {
		t.Errorf("Image URL content not added correctly")
	}

	// Check other properties
	if input.MaxTokens != 100 {
		t.Errorf("Expected MaxTokens to be 100, got %d", input.MaxTokens)
	}

	if input.Temperature != 0.7 {
		t.Errorf("Expected Temperature to be 0.7, got %f", input.Temperature)
	}

	if !input.Stream {
		t.Errorf("Expected Stream to be true")
	}

	if input.Metadata["system_prompt"] != "You are a helpful assistant." {
		t.Errorf("Metadata not set correctly")
	}
}

func TestMultimodalOutput(t *testing.T) {
	output := NewOutput()
	text := "This is a response."
	imageData := []byte("fake-image-data")
	imageMimeType := "image/jpeg"

	// Test fluent API
	output.AddText(text).
		AddImage(imageData, imageMimeType)

	// Check contents
	if len(output.Contents) != 2 {
		t.Errorf("Expected 2 contents, got %d", len(output.Contents))
	}

	if output.Contents[0].Type != MediaTypeText || output.Contents[0].Text != text {
		t.Errorf("Text content not added correctly")
	}

	if output.Contents[1].Type != MediaTypeImage || !bytes.Equal(output.Contents[1].Data, imageData) {
		t.Errorf("Image content not added correctly")
	}

	// Test GetText
	if output.GetText() != text {
		t.Errorf("Expected GetText to return %s, got %s", text, output.GetText())
	}

	// Test GetFirstImage
	firstImage := output.GetFirstImage()
	if firstImage == nil {
		t.Errorf("Expected GetFirstImage to return an image")
	} else if !bytes.Equal(firstImage.Data, imageData) {
		t.Errorf("GetFirstImage returned incorrect image")
	}

	// Test with text-only output
	textOnlyOutput := NewOutput().AddText(text)
	if textOnlyOutput.GetFirstImage() != nil {
		t.Errorf("Expected GetFirstImage to return nil for text-only output")
	}
}

func TestContentReader(t *testing.T) {
	data := []byte("test-data")
	mimeType := "image/jpeg"
	content := NewImageContent(data, mimeType)

	reader := NewContentReader(content)

	// Check content type
	if reader.ContentType() != mimeType {
		t.Errorf("Expected content type %s, got %s", mimeType, reader.ContentType())
	}

	// Read from reader
	buf := make([]byte, len(data))
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("Unexpected error reading from content reader: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to read %d bytes, got %d", len(data), n)
	}
	if !bytes.Equal(buf, data) {
		t.Errorf("Data read from reader doesn't match original data")
	}
}

func TestMultipleTextContents(t *testing.T) {
	output := NewOutput()
	text1 := "Hello, "
	text2 := "world!"

	output.AddText(text1).AddText(text2)

	combined := output.GetText()
	expected := text1 + text2

	if combined != expected {
		t.Errorf("Expected GetText to return %s, got %s", expected, combined)
	}
}

func TestEmptyContents(t *testing.T) {
	input := NewInput()
	output := NewOutput()

	if len(input.Contents) != 0 {
		t.Errorf("Expected empty input contents")
	}

	if len(output.Contents) != 0 {
		t.Errorf("Expected empty output contents")
	}

	if output.GetText() != "" {
		t.Errorf("Expected empty GetText result")
	}

	if output.GetFirstImage() != nil {
		t.Errorf("Expected nil GetFirstImage result")
	}
}

func TestContentMetadata(t *testing.T) {
	content := NewTextContent("Hello")
	content.Metadata = map[string]interface{}{
		"priority": "high",
		"source":   "user",
	}

	if content.Metadata["priority"] != "high" {
		t.Errorf("Expected metadata priority to be high")
	}

	if content.Metadata["source"] != "user" {
		t.Errorf("Expected metadata source to be user")
	}
}

func TestChunk(t *testing.T) {
	content := NewTextContent("Hello")
	chunk := &Chunk{
		Content: content,
		IsFinal: true,
		Error:   nil,
	}

	if chunk.Content != content {
		t.Errorf("Expected chunk content to match")
	}

	if !chunk.IsFinal {
		t.Errorf("Expected chunk to be final")
	}

	if chunk.Error != nil {
		t.Errorf("Expected chunk error to be nil")
	}

	// Test with error
	errorChunk := &Chunk{
		Content: nil,
		IsFinal: true,
		Error:   io.EOF,
	}

	if errorChunk.Error != io.EOF {
		t.Errorf("Expected chunk error to be EOF")
	}
}
