package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections
type WebSocketManager struct {
	connections map[string][]*websocket.Conn
	mu          sync.RWMutex
	log         *log.Logger
}

// ChatMessage represents a message in a chat session
type ChatMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// EventMessage represents an event message
type EventMessage struct {
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(logger *log.Logger) *WebSocketManager {
	if logger == nil {
		logger = log.New(os.Stdout, "[WebSocket] ", log.LstdFlags)
	}

	return &WebSocketManager{
		connections: make(map[string][]*websocket.Conn),
		log:         logger,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// @Summary Connect to agent chat
// @Description Start a WebSocket connection for chatting with an agent
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 101 {string} string "Switching to WebSocket protocol"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agents/{id}/chat [get]
func (s *Server) HandleAgentChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	// Check if agent exists
	agent, err := s.runtime.GetAgent(agentID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Register connection
	s.wsManager.mu.Lock()
	connectionKey := fmt.Sprintf("chat:%s", agentID)
	s.wsManager.connections[connectionKey] = append(s.wsManager.connections[connectionKey], conn)
	s.wsManager.mu.Unlock()

	// Send welcome message
	welcomeMsg := ChatMessage{
		Role:      "system",
		Content:   fmt.Sprintf("Connected to %s. Send a message to start chatting.", agent.Name),
		Timestamp: time.Now().Format(time.RFC3339),
	}
	conn.WriteJSON(welcomeMsg)

	// WebSocket message handling loop
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			s.log.Printf("Error reading message: %v", err)
			break
		}

		// Handle message based on type
		if messageType == websocket.TextMessage {
			var chatMsg ChatMessage
			if err := json.Unmarshal(p, &chatMsg); err != nil {
				s.log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Set timestamp if not provided
			if chatMsg.Timestamp == "" {
				chatMsg.Timestamp = time.Now().Format(time.RFC3339)
			}

			// Echo message back for now
			conn.WriteJSON(chatMsg)

			// Simulate agent response after a short delay
			go func() {
				time.Sleep(1 * time.Second)
				responseMsg := ChatMessage{
					Role:      "assistant",
					Content:   fmt.Sprintf("This is a simulated response to: \"%s\"", chatMsg.Content),
					Timestamp: time.Now().Format(time.RFC3339),
				}
				conn.WriteJSON(responseMsg)
			}()
		}
	}

	// Unregister connection when done
	s.wsManager.mu.Lock()
	connections := s.wsManager.connections[connectionKey]
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	s.wsManager.connections[connectionKey] = connections
	s.wsManager.mu.Unlock()
}

// @Summary Connect to agent events
// @Description Start a WebSocket connection for receiving agent events
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 101 {string} string "Switching to WebSocket protocol"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /agents/{id}/events [get]
func (s *Server) HandleAgentEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	// Check if agent exists
	_, err := s.runtime.GetAgent(agentID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Register connection
	s.wsManager.mu.Lock()
	connectionKey := fmt.Sprintf("events:%s", agentID)
	s.wsManager.connections[connectionKey] = append(s.wsManager.connections[connectionKey], conn)
	s.wsManager.mu.Unlock()

	// Send initial connection event
	connectedEvent := EventMessage{
		Type:      "connected",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]string{
			"message": "Connected to event stream",
			"agentId": agentID,
		},
	}
	conn.WriteJSON(connectedEvent)

	// Send periodic status updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				statusEvent := EventMessage{
					Type:      "status",
					Timestamp: t.Format(time.RFC3339),
					Data: map[string]string{
						"status":  "running",
						"message": "Agent is processing requests",
					},
				}
				if err := conn.WriteJSON(statusEvent); err != nil {
					s.log.Printf("Error writing status: %v", err)
					return
				}
			}
		}
	}()

	// Keep connection open until closed by client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			s.log.Printf("WebSocket read error: %v", err)
			break
		}
	}

	// Unregister connection when done
	s.wsManager.mu.Lock()
	connections := s.wsManager.connections[connectionKey]
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	s.wsManager.connections[connectionKey] = connections
	s.wsManager.mu.Unlock()
}

// BroadcastEvent sends an event to all clients subscribed to a specific agent
func (s *WebSocketManager) BroadcastEvent(agentID, eventType string, data interface{}) {
	s.mu.RLock()
	connections := s.connections[fmt.Sprintf("events:%s", agentID)]
	s.mu.RUnlock()

	event := EventMessage{
		Type:      eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}

	for _, conn := range connections {
		if err := conn.WriteJSON(event); err != nil {
			s.log.Printf("Error broadcasting event: %v", err)
		}
	}
}
