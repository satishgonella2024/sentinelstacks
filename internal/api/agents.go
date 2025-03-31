package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sentinelstacks/sentinel/internal/runtime"
)

// AgentRequest represents a request to create a new agent
type AgentRequest struct {
	Image       string                 `json:"image"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
}

// AgentResponse represents an agent response
type AgentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Endpoints struct {
		Chat   string `json:"chat"`
		Events string `json:"events"`
	} `json:"endpoints,omitempty"`
}

// listAgentsHandler handles GET /agents
func (s *Server) listAgentsHandler(w http.ResponseWriter, r *http.Request) {
	// Get agents from the runtime
	agents, err := s.runtime.GetRunningAgents()
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to list agents")
		s.log.Printf("Error listing agents: %v", err)
		return
	}

	// Convert runtime.AgentInfo to AgentResponse
	response := make([]AgentResponse, 0, len(agents))
	for _, agent := range agents {
		agentResp := convertAgentInfoToResponse(agent, s.config.Host, s.config.Port)
		response = append(response, agentResp)
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"agents": response,
	})
}

// getAgentHandler handles GET /agents/{id}
func (s *Server) getAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	agent, err := s.runtime.GetAgent(id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	response := convertAgentInfoToResponse(agent, s.config.Host, s.config.Port)
	s.sendJSON(w, http.StatusOK, response)
}

// createAgentHandler handles POST /agents
func (s *Server) createAgentHandler(w http.ResponseWriter, r *http.Request) {
	var req AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Image == "" {
		s.sendError(w, http.StatusBadRequest, "Image is required")
		return
	}

	// TODO: Implement agent creation
	// For now, return a dummy response with a "not implemented" message
	s.sendError(w, http.StatusNotImplemented, "Agent creation not implemented")
}

// deleteAgentHandler handles DELETE /agents/{id}
func (s *Server) deleteAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete agent from runtime
	err := s.runtime.DeleteAgent(id)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete agent: %v", err))
		return
	}

	s.sendJSON(w, http.StatusOK, map[string]string{
		"id":     id,
		"status": "stopping",
	})
}

// getAgentLogsHandler handles GET /agents/{id}/logs
func (s *Server) getAgentLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if agent exists
	_, err := s.runtime.GetAgent(id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// TODO: Implement log retrieval
	// For now, return dummy log data
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			"level":     "info",
			"message":   "Agent initialized successfully",
		},
		{
			"timestamp": time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
			"level":     "debug",
			"message":   "Connected to LLM provider",
		},
		{
			"timestamp": time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			"level":     "info",
			"message":   "Processed user request",
		},
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"logs": logs,
	})
}

// Helper function to convert AgentInfo to AgentResponse
func convertAgentInfoToResponse(info runtime.AgentInfo, host string, port int) AgentResponse {
	var response AgentResponse
	response.ID = info.ID
	response.Name = info.Name
	response.Image = info.Image
	response.Status = info.Status
	response.CreatedAt = info.CreatedAt

	// Add endpoint URLs
	response.Endpoints.Chat = fmt.Sprintf("ws://%s:%d/v1/agents/%s/chat", host, port, info.ID)
	response.Endpoints.Events = fmt.Sprintf("ws://%s:%d/v1/agents/%s/events", host, port, info.ID)

	return response
}
