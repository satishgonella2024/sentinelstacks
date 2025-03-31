// Package conversation provides conversation history management
package conversation

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sentinelstacks/sentinel/internal/multimodal"
)

// MessageType represents the type of a message
type MessageType string

const (
	// MessageTypeUser represents a message from the user
	MessageTypeUser MessageType = "user"
	// MessageTypeAssistant represents a message from the assistant
	MessageTypeAssistant MessageType = "assistant"
	// MessageTypeSystem represents a system message
	MessageTypeSystem MessageType = "system"
)

// Message represents a message in the conversation
type Message struct {
	ID        string                `json:"id"`
	Role      MessageType           `json:"role"`
	Content   string                `json:"content,omitempty"`
	Contents  []*multimodal.Content `json:"contents,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}

// History represents a conversation history
type History struct {
	ID       string     `json:"id"`
	AgentID  string     `json:"agent_id"`
	Messages []*Message `json:"messages"`
}

// NewHistory creates a new conversation history
func NewHistory(agentID, sessionID string) *History {
	return &History{
		ID:       sessionID,
		AgentID:  agentID,
		Messages: []*Message{},
	}
}

// AddUserMessage adds a user message to the history
func (h *History) AddUserMessage(content string) {
	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      MessageTypeUser,
		Content:   content,
		Timestamp: time.Now(),
	})
}

// AddUserMultimodalMessage adds a multimodal user message to the history
func (h *History) AddUserMultimodalMessage(contents []*multimodal.Content) {
	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      MessageTypeUser,
		Contents:  contents,
		Timestamp: time.Now(),
	})
}

// AddAssistantMessage adds an assistant message to the history
func (h *History) AddAssistantMessage(content string) {
	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      MessageTypeAssistant,
		Content:   content,
		Timestamp: time.Now(),
	})
}

// AddSystemMessage adds a system message to the history
func (h *History) AddSystemMessage(content string) {
	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      MessageTypeSystem,
		Content:   content,
		Timestamp: time.Now(),
	})
}

// GetLastMessages returns the last n messages in the history
func (h *History) GetLastMessages(n int) []*Message {
	if n <= 0 || n > len(h.Messages) {
		return h.Messages
	}
	return h.Messages[len(h.Messages)-n:]
}

// ToMultimodalInput converts the conversation history to a multimodal input
func (h *History) ToMultimodalInput(messageLimit int) (*multimodal.Input, error) {
	// Get the last N messages
	messages := h.GetLastMessages(messageLimit)

	// Create input
	input := multimodal.NewInput()

	// Group contents by role
	var systemContents []*multimodal.Content
	var userContents []*multimodal.Content
	var assistantContents []*multimodal.Content

	// Convert messages to contents
	for _, msg := range messages {
		if msg.Contents != nil && len(msg.Contents) > 0 {
			// Message has multimodal contents
			for _, content := range msg.Contents {
				// Add to appropriate role group
				switch msg.Role {
				case MessageTypeUser:
					userContents = append(userContents, content)
				case MessageTypeAssistant:
					assistantContents = append(assistantContents, content)
				case MessageTypeSystem:
					systemContents = append(systemContents, content)
				}
			}
		} else if msg.Content != "" {
			// Message has text content
			content := multimodal.NewTextContent(msg.Content)

			// Add to appropriate role group
			switch msg.Role {
			case MessageTypeUser:
				userContents = append(userContents, content)
			case MessageTypeAssistant:
				assistantContents = append(assistantContents, content)
			case MessageTypeSystem:
				systemContents = append(systemContents, content)
			}
		}
	}

	// Add system role information in metadata if present
	if len(systemContents) > 0 {
		// Combine all system messages into one prompt
		systemPrompt := ""
		for _, content := range systemContents {
			if content.Type == multimodal.MediaTypeText {
				systemPrompt += content.Text + "\n"
			}
		}
		if systemPrompt != "" {
			input.SetMetadata("system_prompt", systemPrompt)
		}
	}

	// Add content with appropriate metadata for roles
	// First add all user contents
	for _, content := range userContents {
		content.Type = multimodal.MediaTypeText
		input.AddContent(content)
		input.SetMetadata("role", "user")
	}

	// Then all assistant contents
	for _, content := range assistantContents {
		content.Type = multimodal.MediaTypeText
		input.AddContent(content)
		input.SetMetadata("role", "assistant")
	}

	return input, nil
}

// SaveToFile saves the conversation history to a file
func (h *History) SaveToFile(filename string) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	// Marshal to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(h); err != nil {
		return fmt.Errorf("could not encode history: %w", err)
	}

	return nil
}

// LoadFromFile loads a conversation history from a file
func LoadFromFile(filename string) (*History, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Unmarshal from JSON
	var history History
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&history); err != nil {
		return nil, fmt.Errorf("could not decode history: %w", err)
	}

	return &history, nil
}