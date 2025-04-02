package runtime

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/satishgonella2024/sentinelstacks/internal/multimodal"
	"github.com/satishgonella2024/sentinelstacks/internal/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultimodalAgent tests the multimodal agent functionality
func TestMultimodalAgent(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := ioutil.TempDir("", "multimodal-agent-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test image file
	imagePath := filepath.Join(tmpDir, "test.jpg")
	require.NoError(t, createTestImageFile(imagePath))

	// Initialize a basic agent
	agent := &Agent{
		ID:       "test-agent",
		Name:     "Test Agent",
		Image:    "mock:latest",
		Status:   StatusRunning,
		Model:    "mock-model",
		StateDir: tmpDir,
	}

	// Create a multimodal agent
	shimConfig := shim.Config{
		Provider: "mock",
		Model:    "mock-model",
		Parameters: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  1024,
		},
	}

	// Create multimodal agent
	mmAgent, err := NewMultimodalAgent(agent, shimConfig)
	require.NoError(t, err)
	assert.NotNil(t, mmAgent)

	// Add system prompt
	mmAgent.AddSystemPrompt("You are a helpful assistant.")

	// Process text input
	ctx := context.Background()
	resp, err := mmAgent.ProcessTextInput(ctx, "Hello, world!")
	require.NoError(t, err)
	assert.Contains(t, resp, "Mock")
	assert.Contains(t, resp, "Hello, world!")

	// Process multimodal input with image
	imageData, err := ioutil.ReadFile(imagePath)
	require.NoError(t, err)

	contents := []*multimodal.Content{
		multimodal.NewTextContent("What's in this image?"),
		multimodal.NewImageContent(imageData, "image/jpeg"),
	}

	output, err := mmAgent.ProcessMultimodalInput(ctx, contents)
	require.NoError(t, err)
	assert.NotNil(t, output)

	// Extract text from output
	var responseText string
	for _, content := range output.Contents {
		if content.Type == multimodal.MediaTypeText {
			responseText = content.Text
			break
		}
	}
	assert.Contains(t, responseText, "Mock multimodal response")

	// Check conversation history
	history := mmAgent.GetConversationHistory()
	assert.NotNil(t, history)
	assert.GreaterOrEqual(t, len(history.Messages), 3) // System prompt + text input + multimodal input
	assert.Equal(t, "system", string(history.Messages[0].Role))

	// Close the agent
	err = mmAgent.Close()
	require.NoError(t, err)

	// Verify conversation file exists
	conversationFile := filepath.Join(tmpDir, "conversations", history.ID+".json")
	_, err = os.Stat(conversationFile)
	assert.NoError(t, err)
}

// Helper function to create a test image file
func createTestImageFile(path string) error {
	// Simple JPG file header (minimal valid JPG)
	data := []byte{
		0xFF, 0xD8, // SOI marker
		0xFF, 0xE0, // APP0 marker
		0x00, 0x10, // Length of APP0
		0x4A, 0x46, 0x49, 0x46, 0x00, // JFIF identifier
		0x01, 0x01, // JFIF version
		0x00,       // Density units
		0x00, 0x01, // X density
		0x00, 0x01, // Y density
		0x00, 0x00, // Thumbnail width and height
		0xFF, 0xD9, // EOI marker
	}
	return ioutil.WriteFile(path, data, 0644)
}
