package network

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MessageFormat represents the format of a message
type MessageFormat string

const (
	// TextFormat represents plain text messages
	TextFormat MessageFormat = "text"
	// ImageFormat represents image messages
	ImageFormat MessageFormat = "image"
	// AudioFormat represents audio messages
	AudioFormat MessageFormat = "audio"
	// VideoFormat represents video messages
	VideoFormat MessageFormat = "video"
	// BinaryFormat represents binary data
	BinaryFormat MessageFormat = "binary"
	// JSONFormat represents structured JSON data
	JSONFormat MessageFormat = "json"
)

// Message represents a communication message between agents
type Message struct {
	ID          string                 `json:"id"`
	NetworkID   string                 `json:"network_id"`
	SenderID    string                 `json:"sender_id"`
	Format      MessageFormat          `json:"format"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Attachments []Attachment           `json:"attachments,omitempty"`
}

// Attachment represents a file or binary data attached to a message
type Attachment struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Format   MessageFormat `json:"format"`
	Path     string       `json:"path"`
	Size     int64        `json:"size"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MessageHandler is a function that processes a message
type MessageHandler func(message Message) error

// MessagingSystem manages message delivery between agents
type MessagingSystem struct {
	networkManager *NetworkManager
	dataDir        string
	mutex          sync.RWMutex
	handlers       map[MessageFormat][]MessageHandler
	subscribers    map[string][]MessageHandler // subscriber handlers by network ID
}

// NewMessagingSystem creates a new messaging system
func NewMessagingSystem(networkManager *NetworkManager) *MessagingSystem {
	return &MessagingSystem{
		networkManager: networkManager,
		dataDir:        networkManager.DataDir,
		handlers:       make(map[MessageFormat][]MessageHandler),
		subscribers:    make(map[string][]MessageHandler),
	}
}

// RegisterFormatHandler registers a handler for a specific message format
func (ms *MessagingSystem) RegisterFormatHandler(format MessageFormat, handler MessageHandler) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	
	if _, exists := ms.handlers[format]; !exists {
		ms.handlers[format] = []MessageHandler{}
	}
	
	ms.handlers[format] = append(ms.handlers[format], handler)
}

// SubscribeToNetwork registers a handler for messages in a specific network
func (ms *MessagingSystem) SubscribeToNetwork(networkID string, handler MessageHandler) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	
	if _, exists := ms.subscribers[networkID]; !exists {
		ms.subscribers[networkID] = []MessageHandler{}
	}
	
	ms.subscribers[networkID] = append(ms.subscribers[networkID], handler)
}

// SendMessage sends a message to a network
func (ms *MessagingSystem) SendMessage(networkName string, message Message) (Message, error) {
	// Use a mutex to avoid concurrent map writes
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Assign ID and timestamp if not present
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}
	
	// Get network
	network, err := ms.networkManager.GetNetworkByName(networkName)
	if err != nil {
		return Message{}, err
	}
	
	// Verify format is supported
	formatSupported := false
	for _, format := range network.SupportedFormats {
		if string(message.Format) == format {
			formatSupported = true
			break
		}
	}
	
	if !formatSupported {
		return Message{}, fmt.Errorf("format '%s' is not supported by network '%s'", message.Format, networkName)
	}
	
	// Set network ID
	message.NetworkID = network.ID
	
	// Verify sender is connected to the network
	senderConnected := false
	for _, agentID := range network.Agents {
		if agentID == message.SenderID {
			senderConnected = true
			break
		}
	}
	
	if !senderConnected {
		return Message{}, fmt.Errorf("sender '%s' is not connected to network '%s'", message.SenderID, networkName)
	}
	
	// Process attachments if any
	if len(message.Attachments) > 0 {
		attachmentDir := filepath.Join(ms.dataDir, "attachments", message.ID)
		if err := os.MkdirAll(attachmentDir, 0755); err != nil {
			return Message{}, fmt.Errorf("failed to create attachment directory: %w", err)
		}
		
		for i, attachment := range message.Attachments {
			if attachment.ID == "" {
				attachment.ID = uuid.New().String()
				message.Attachments[i] = attachment
			}
		}
	}
	
	// Save message
	if err := ms.saveMessage(message); err != nil {
		return Message{}, err
	}
	
	// Process message with format handlers
	ms.mutex.RLock()
	formatHandlers, exists := ms.handlers[message.Format]
	networkHandlers := ms.subscribers[network.ID]
	ms.mutex.RUnlock()
	
	// Execute format handlers
	if exists {
		for _, handler := range formatHandlers {
			if err := handler(message); err != nil {
				// Log error but continue processing
				fmt.Printf("Error in format handler: %v\n", err)
			}
		}
	}
	
	// Execute network-specific handlers
	for _, handler := range networkHandlers {
		if err := handler(message); err != nil {
			// Log error but continue processing
			fmt.Printf("Error in network handler: %v\n", err)
		}
	}
	
	return message, nil
}

// GetMessages retrieves messages for a specific network
func (ms *MessagingSystem) GetMessages(networkName string, limit int, offset int) ([]Message, error) {
	// Get network
	network, err := ms.networkManager.GetNetworkByName(networkName)
	if err != nil {
		return nil, err
	}
	
	return ms.getMessagesForNetwork(network.ID, limit, offset)
}

// getMessagesForNetwork is an internal helper to get messages for a network ID
func (ms *MessagingSystem) getMessagesForNetwork(networkID string, limit int, offset int) ([]Message, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	
	messagesDir := filepath.Join(ms.dataDir, "messages", networkID)
	if _, err := os.Stat(messagesDir); os.IsNotExist(err) {
		return []Message{}, nil
	}
	
	files, err := os.ReadDir(messagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages directory: %w", err)
	}
	
	var messages []Message
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		filePath := filepath.Join(messagesDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}
		
		var message Message
		if err := json.Unmarshal(data, &message); err != nil {
			continue // Skip files that can't be parsed
		}
		
		messages = append(messages, message)
	}
	
	// Apply pagination
	if limit > 0 && offset >= 0 {
		if offset >= len(messages) {
			return []Message{}, nil
		}
		
		end := offset + limit
		if end > len(messages) {
			end = len(messages)
		}
		
		return messages[offset:end], nil
	}
	
	return messages, nil
}

// saveMessage saves a message to disk
func (ms *MessagingSystem) saveMessage(message Message) error {
	// Create messages directory if it doesn't exist
	messagesDir := filepath.Join(ms.dataDir, "messages", message.NetworkID)
	if err := os.MkdirAll(messagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create messages directory: %w", err)
	}
	
	// Convert to JSON
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %w", err)
	}
	
	// Write to file
	filePath := filepath.Join(messagesDir, message.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write message data: %w", err)
	}
	
	return nil
}

// GetMessageByID retrieves a specific message by ID
func (ms *MessagingSystem) GetMessageByID(networkID, messageID string) (Message, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	
	filePath := filepath.Join(ms.dataDir, "messages", networkID, messageID+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Message{}, fmt.Errorf("message not found: %s", messageID)
		}
		return Message{}, fmt.Errorf("failed to read message data: %w", err)
	}
	
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return Message{}, fmt.Errorf("failed to parse message data: %w", err)
	}
	
	return message, nil
}
