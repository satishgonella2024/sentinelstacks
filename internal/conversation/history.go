// Package conversation provides conversation history management
package conversation

import (
	"encoding/json"
	"fmt"
	"io"
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
	AgentID  string     `json:"agent_id,omitempty"`
	Messages []*Message `json:"messages"`
}

// NewHistory creates a new conversation history
func NewHistory() *History {
	return &History{
		ID:       fmt.Sprintf("history_%s", uuid.New().String()),
		Messages: []*Message{},
	}
}

// SetID sets the ID of the history
func (h *History) SetID(id string) {
	h.ID = id
}

// GetID returns the ID of the history
func (h *History) GetID() string {
	return h.ID
}

// SetAgentID sets the agent ID of the history
func (h *History) SetAgentID(agentID string) {
	h.AgentID = agentID
}

// AddMessage adds a message to the history
func (h *History) AddMessage(role string, content string) {
	// Convert role string to MessageType
	var messageType MessageType
	switch role {
	case "user":
		messageType = MessageTypeUser
	case "assistant":
		messageType = MessageTypeAssistant
	case "system":
		messageType = MessageTypeSystem
	default:
		// Default to user if unknown
		messageType = MessageTypeUser
	}

	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      messageType,
		Content:   content,
		Timestamp: time.Now(),
	})
}

// AddMultimodalMessage adds a multimodal message to the history
func (h *History) AddMultimodalMessage(role string, contents []*multimodal.Content) {
	// Convert role string to MessageType
	var messageType MessageType
	switch role {
	case "user":
		messageType = MessageTypeUser
	case "assistant":
		messageType = MessageTypeAssistant
	case "system":
		messageType = MessageTypeSystem
	default:
		// Default to user if unknown
		messageType = MessageTypeUser
	}

	h.Messages = append(h.Messages, &Message{
		ID:        uuid.New().String(),
		Role:      messageType,
		Contents:  contents,
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

// MessageCount returns the number of messages in the history
func (h *History) MessageCount() int {
	return len(h.Messages)
}

// ToMultimodalInput converts the conversation history to a multimodal input
func (h *History) ToMultimodalInput(messageLimit int) (*multimodal.Input, error) {
	// Get the last N messages
	messages := h.GetLastMessages(messageLimit)

	// Create input
	input := multimodal.NewInput()

	// Find system message first
	var systemPrompt string
	for _, msg := range messages {
		if msg.Role == MessageTypeSystem && msg.Content != "" {
			systemPrompt += msg.Content + "\n"
		}
	}
	
	// Add system prompt to metadata if present
	if systemPrompt != "" {
		input.SetMetadata("system", systemPrompt)
	}

	// Add conversation as text context
	var conversationText string
	for _, msg := range messages {
		// Skip system messages in the conversation text
		if msg.Role == MessageTypeSystem {
			continue
		}
		
		// Get the role as string
		roleStr := "User"
		if msg.Role == MessageTypeAssistant {
			roleStr = "Assistant"
		}
		
		// Add message to conversation text
		if msg.Content != "" {
			conversationText += fmt.Sprintf("%s: %s\n\n", roleStr, msg.Content)
		} else if msg.Contents != nil && len(msg.Contents) > 0 {
			// For multimodal contents, just add text parts
			textContent := ""
			for _, content := range msg.Contents {
				if content.Type == multimodal.MediaTypeText {
					textContent += content.Text + " "
				}
			}
			if textContent != "" {
				conversationText += fmt.Sprintf("%s: %s\n\n", roleStr, textContent)
			}
		}
	}
	
	// Add the conversation context
	if conversationText != "" {
		input.AddText("Conversation history:\n" + conversationText)
	}

	// If the last message is from user and has multimodal content,
	// add those contents directly to the input
	if len(messages) > 0 && messages[len(messages)-1].Role == MessageTypeUser {
		lastMsg := messages[len(messages)-1]
		if lastMsg.Contents != nil && len(lastMsg.Contents) > 0 {
			for _, content := range lastMsg.Contents {
				// Only add non-text contents directly (images, etc.)
				if content.Type != multimodal.MediaTypeText {
					input.AddContent(content)
				}
			}
		}
	}

	return input, nil
}

// Save saves the conversation history to a writer
func (h *History) Save(writer io.Writer) error {
	// Marshal to JSON
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(h); err != nil {
		return fmt.Errorf("could not encode history: %w", err)
	}

	return nil
}

// SaveToFile saves the conversation history to a file
func (h *History) SaveToFile(filename string) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	// Save to the file
	return h.Save(file)
}

// Load loads a conversation history from a reader
func Load(reader io.Reader) (*History, error) {
	// Unmarshal from JSON
	var history History
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&history); err != nil {
		return nil, fmt.Errorf("could not decode history: %w", err)
	}

	return &history, nil
}

// LoadFromFile loads a conversation history from a file
func LoadFromFile(filename string) (*History, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Load from the file
	return Load(file)
}
