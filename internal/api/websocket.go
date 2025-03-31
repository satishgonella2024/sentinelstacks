package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections and message routing
type WebSocketManager struct {
	// Connection upgrader with configuration
	upgrader websocket.Upgrader

	// Connections maps agent IDs to active client connections
	connections map[string]map[*websocket.Conn]bool

	// mutex for protecting the connections map
	mu sync.RWMutex

	// Logger for WebSocket events
	log *log.Logger
}

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Message types for client -> server
	TypeMessage     MessageType = "message"      // Regular message from user
	TypeToolRequest MessageType = "tool_request" // Request to use a tool

	// Message types for server -> client
	TypeResponse    MessageType = "response"     // Response from agent
	TypeEvent       MessageType = "event"        // Event notification
	TypeToolResult  MessageType = "tool_result"  // Result from tool execution
	TypeError       MessageType = "error"        // Error message
	TypeStreamStart MessageType = "stream_start" // Start of stream
	TypeStreamChunk MessageType = "stream_chunk" // Stream chunk
	TypeStreamEnd   MessageType = "stream_end"   // End of stream
)

// ClientMessage represents a message sent from the client
type ClientMessage struct {
	Type       MessageType            `json:"type"`
	Content    string                 `json:"content,omitempty"`
	MessageID  string                 `json:"message_id,omitempty"`
	Tool       string                 `json:"tool,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
}

// ServerMessage represents a message sent to the client
type ServerMessage struct {
	Type        MessageType            `json:"type"`
	Content     string                 `json:"content,omitempty"`
	MessageID   string                 `json:"message_id,omitempty"`
	ResponseID  string                 `json:"response_id,omitempty"`
	ToolsUsed   []string               `json:"tools_used,omitempty"`
	EventType   string                 `json:"event_type,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	IsComplete  bool                   `json:"is_complete,omitempty"`
	ChunkIndex  int                    `json:"chunk_index,omitempty"`
	TotalChunks int                    `json:"total_chunks,omitempty"`
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(logger *log.Logger) *WebSocketManager {
	if logger == nil {
		logger = log.New(os.Stdout, "[WebSocket] ", log.LstdFlags)
	}

	return &WebSocketManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now, this should be more restrictive in production
				return true
			},
		},
		connections: make(map[string]map[*websocket.Conn]bool),
		log:         logger,
	}
}

// HandleAgentChat handles WebSocket connections for agent chat
func (s *Server) HandleAgentChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	// Check if agent exists
	agent, err := s.runtime.GetAgent(agentID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := s.wsManager.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	// Register the connection
	s.wsManager.registerConnection(agentID, conn)
	defer s.wsManager.unregisterConnection(agentID, conn)

	// Send welcome message
	welcomeMsg := ServerMessage{
		Type:      TypeEvent,
		EventType: "connection_established",
		Timestamp: time.Now().Format(time.RFC3339),
		Content:   fmt.Sprintf("Connected to agent %s (%s)", agent.Name, agent.ID),
		Data: map[string]interface{}{
			"agent_id":   agent.ID,
			"agent_name": agent.Name,
			"status":     agent.Status,
		},
	}

	if err := conn.WriteJSON(welcomeMsg); err != nil {
		s.log.Printf("Failed to send welcome message: %v", err)
		return
	}

	// Handle incoming messages
	for {
		var msg ClientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Printf("WebSocket closed unexpectedly: %v", err)
			}
			break
		}

		// Handle the message based on its type
		switch msg.Type {
		case TypeMessage:
			s.handleChatMessage(agentID, conn, msg)
		case TypeToolRequest:
			s.handleToolRequest(agentID, conn, msg)
		default:
			// Send error for unknown message type
			errMsg := ServerMessage{
				Type:      TypeError,
				MessageID: msg.MessageID,
				Error:     "Unknown message type",
			}
			conn.WriteJSON(errMsg)
		}
	}
}

// HandleAgentEvents handles WebSocket connections for agent events
func (s *Server) HandleAgentEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	// Check if agent exists
	agent, err := s.runtime.GetAgent(agentID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := s.wsManager.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	// Register the connection for events
	s.wsManager.registerConnection(agentID+"-events", conn)
	defer s.wsManager.unregisterConnection(agentID+"-events", conn)

	// Send initial event with agent status
	statusMsg := ServerMessage{
		Type:      TypeEvent,
		EventType: "agent_status",
		Timestamp: time.Now().Format(time.RFC3339),
		Content:   fmt.Sprintf("Agent %s is %s", agent.Name, agent.Status),
		Data: map[string]interface{}{
			"agent_id":   agent.ID,
			"agent_name": agent.Name,
			"status":     agent.Status,
			"memory":     agent.Memory,
			"api_usage":  agent.APIUsage,
		},
	}

	if err := conn.WriteJSON(statusMsg); err != nil {
		s.log.Printf("Failed to send status message: %v", err)
		return
	}

	// Keep the connection alive until it's closed
	for {
		// Read messages to detect disconnection
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Printf("Events WebSocket closed unexpectedly: %v", err)
			}
			break
		}
	}
}

// registerConnection registers a WebSocket connection for an agent
func (wm *WebSocketManager) registerConnection(agentID string, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Create map for agent if it doesn't exist
	if _, ok := wm.connections[agentID]; !ok {
		wm.connections[agentID] = make(map[*websocket.Conn]bool)
	}

	// Register the connection
	wm.connections[agentID][conn] = true
	wm.log.Printf("Client connected to agent %s, total connections: %d", agentID, len(wm.connections[agentID]))
}

// unregisterConnection removes a WebSocket connection for an agent
func (wm *WebSocketManager) unregisterConnection(agentID string, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Remove the connection
	if _, ok := wm.connections[agentID]; ok {
		delete(wm.connections[agentID], conn)

		// Clean up empty maps
		if len(wm.connections[agentID]) == 0 {
			delete(wm.connections, agentID)
		}
	}

	// Close the connection
	conn.Close()
	wm.log.Printf("Client disconnected from agent %s", agentID)
}

// broadcastToAgent sends a message to all connections for an agent
func (wm *WebSocketManager) broadcastToAgent(agentID string, msg interface{}) {
	wm.mu.RLock()
	conns := wm.connections[agentID]
	wm.mu.RUnlock()

	// Send message to all connections
	for conn := range conns {
		err := conn.WriteJSON(msg)
		if err != nil {
			wm.log.Printf("Failed to broadcast message: %v", err)
			// Don't unregister here, let the connection handler do it
		}
	}
}

// handleChatMessage processes a chat message from a client
func (s *Server) handleChatMessage(agentID string, conn *websocket.Conn, msg ClientMessage) {
	// Generate a response ID if not provided
	if msg.MessageID == "" {
		msg.MessageID = uuid.New().String()
	}

	// In a real implementation, this would send the message to the agent runtime
	// For now, we'll simulate a response after a short delay

	// Send a "thinking" event
	thinkingMsg := ServerMessage{
		Type:      TypeEvent,
		EventType: "thinking",
		Timestamp: time.Now().Format(time.RFC3339),
		MessageID: msg.MessageID,
		Content:   "Agent is processing your message...",
	}
	conn.WriteJSON(thinkingMsg)

	// Simulate processing time
	go func() {
		// Generate a response ID
		responseID := uuid.New().String()

		// Start streaming (for demo)
		time.Sleep(500 * time.Millisecond)
		conn.WriteJSON(ServerMessage{
			Type:       TypeStreamStart,
			MessageID:  msg.MessageID,
			ResponseID: responseID,
			Timestamp:  time.Now().Format(time.RFC3339),
		})

		// Stream chunks (for demo)
		chunks := []string{
			"I'll help you with that request. ",
			"Based on the information provided, ",
			"here's what I can tell you: ",
			"The answer to your question requires some analysis. ",
			"Let me break it down for you step by step.",
		}

		for i, chunk := range chunks {
			time.Sleep(300 * time.Millisecond)
			conn.WriteJSON(ServerMessage{
				Type:        TypeStreamChunk,
				MessageID:   msg.MessageID,
				ResponseID:  responseID,
				Content:     chunk,
				ChunkIndex:  i,
				TotalChunks: len(chunks),
				Timestamp:   time.Now().Format(time.RFC3339),
			})
		}

		// Finish with complete response
		time.Sleep(500 * time.Millisecond)
		conn.WriteJSON(ServerMessage{
			Type:       TypeStreamEnd,
			MessageID:  msg.MessageID,
			ResponseID: responseID,
			Content:    "I've processed your message: \"" + msg.Content + "\". This is a simulated response since this is a prototype implementation. In a full implementation, I would use the agent runtime to generate a proper response.",
			ToolsUsed:  []string{},
			Timestamp:  time.Now().Format(time.RFC3339),
		})

		// Also send agent status update to event channel
		s.wsManager.broadcastToAgent(agentID+"-events", ServerMessage{
			Type:      TypeEvent,
			EventType: "message_processed",
			Timestamp: time.Now().Format(time.RFC3339),
			Content:   "Agent processed a message",
			Data: map[string]interface{}{
				"agent_id":    agentID,
				"message_id":  msg.MessageID,
				"response_id": responseID,
			},
		})
	}()
}

// handleToolRequest processes a tool request from a client
func (s *Server) handleToolRequest(agentID string, conn *websocket.Conn, msg ClientMessage) {
	// In a real implementation, this would forward the tool request to the agent runtime
	// For now, we'll simulate a response

	// Simulate processing time
	go func() {
		time.Sleep(1 * time.Second)

		// Send tool result
		toolResult := ServerMessage{
			Type:      TypeToolResult,
			RequestID: msg.RequestID,
			Status:    "success",
			Data: map[string]interface{}{
				"result": fmt.Sprintf("Simulated result for tool '%s'", msg.Tool),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		conn.WriteJSON(toolResult)

		// Also send agent status update to event channel
		s.wsManager.broadcastToAgent(agentID+"-events", ServerMessage{
			Type:      TypeEvent,
			EventType: "tool_used",
			Timestamp: time.Now().Format(time.RFC3339),
			Content:   fmt.Sprintf("Agent used tool: %s", msg.Tool),
			Data: map[string]interface{}{
				"agent_id":   agentID,
				"tool":       msg.Tool,
				"request_id": msg.RequestID,
			},
		})
	}()
}
