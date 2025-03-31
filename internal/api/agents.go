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
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	Model        string    `json:"model,omitempty"`
	IsMultimodal bool      `json:"is_multimodal,omitempty"`
	Endpoints    struct {
		Chat   string `json:"chat"`
		Events string `json:"events"`
	} `json:"endpoints,omitempty"`
}

// AgentLogsResponse represents the logs for an agent
type AgentLogsResponse struct {
	Logs []AgentLogEntry `json:"logs"`
}

// AgentLogEntry represents a single log entry
type AgentLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// @Summary List all agents
// @Description Get a list of all running agents
// @Tags agents
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]AgentResponse
// @Failure 500 {object} map[string]string
// @Router /agents [get]
func (s *Server) listAgentsHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Printf("Received request to list agents")

	// Get agents from the runtime
	agents, err := s.runtime.GetRunningAgents()
	if err != nil {
		s.log.Printf("Error listing agents: %v", err)
		// Return empty array instead of error
		s.sendJSON(w, http.StatusOK, map[string]interface{}{
			"agents": []AgentResponse{},
		})
		return
	}

	s.log.Printf("Found %d agents", len(agents))

	// Convert runtime.AgentInfo to AgentResponse
	response := make([]AgentResponse, 0, len(agents))
	for _, agent := range agents {
		s.log.Printf("Converting agent: %s (%s)", agent.ID, agent.Name)
		agentResp := convertAgentInfoToResponse(agent, s.config.Host, s.config.Port)
		response = append(response, agentResp)
	}

	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"agents": response,
	})
}

// @Summary Get agent details
// @Description Get details of a specific agent by ID
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 200 {object} AgentResponse
// @Failure 404 {object} map[string]string
// @Router /agents/{id} [get]
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

// @Summary Create a new agent
// @Description Create a new agent from an image
// @Tags agents
// @Accept json
// @Produce json
// @Param agent body AgentRequest true "Agent Request"
// @Success 201 {object} AgentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agents [post]
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

	// Extract name from parameters if available
	name := "Agent"
	model := "default-model"
	if req.Parameters != nil {
		if nameParam, ok := req.Parameters["name"].(string); ok && nameParam != "" {
			name = nameParam
		}
		if modelParam, ok := req.Parameters["model"].(string); ok && modelParam != "" {
			model = modelParam
		}
	}

	// Create agent in runtime
	agent, err := s.runtime.CreateAgent(name, req.Image, model)
	if err != nil {
		s.log.Printf("Error creating agent: %v", err)
		s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create agent: %v", err))
		return
	}

	// Start the agent
	if err := s.runtime.StartAgent(agent.ID); err != nil {
		s.log.Printf("Error starting agent: %v", err)
		// We still return the created agent, but with a warning
		s.log.Printf("Agent created but failed to start: %s", agent.ID)
	}

	// Return agent info
	agentInfo, err := s.runtime.GetAgent(agent.ID)
	if err != nil {
		s.log.Printf("Error getting agent info: %v", err)
		s.sendError(w, http.StatusInternalServerError, "Failed to get agent info")
		return
	}

	response := convertAgentInfoToResponse(agentInfo, s.config.Host, s.config.Port)
	s.sendJSON(w, http.StatusCreated, response)
}

// @Summary Delete an agent
// @Description Delete an agent by ID
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agents/{id} [delete]
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

// @Summary Get agent logs
// @Description Get logs for a specific agent
// @Tags agents
// @Accept json
// @Produce json
// @Param id path string true "Agent ID"
// @Success 200 {object} AgentLogsResponse
// @Failure 404 {object} map[string]string
// @Router /agents/{id}/logs [get]
func (s *Server) getAgentLogsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if agent exists
	_, err := s.runtime.GetAgent(id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Agent not found")
		return
	}

	// TODO: Implement real log retrieval
	// For now, return dummy log data
	logs := []AgentLogEntry{
		{
			Timestamp: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			Level:     "info",
			Message:   "Agent initialized successfully",
		},
		{
			Timestamp: time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
			Level:     "debug",
			Message:   "Connected to LLM provider",
		},
		{
			Timestamp: time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			Level:     "info",
			Message:   "Processed user request",
		},
	}

	s.sendJSON(w, http.StatusOK, AgentLogsResponse{
		Logs: logs,
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
	response.Model = info.Model

	// Determine if agent is multimodal based on model
	response.IsMultimodal = isMultimodalModel(info.Model)

	// Add endpoint URLs
	response.Endpoints.Chat = fmt.Sprintf("ws://%s:%d/v1/agents/%s/chat", host, port, info.ID)
	response.Endpoints.Events = fmt.Sprintf("ws://%s:%d/v1/agents/%s/events", host, port, info.ID)

	return response
}

// isMultimodalModel checks if a model supports multimodal inputs
func isMultimodalModel(model string) bool {
	// List of known multimodal models
	multimodalModels := map[string]bool{
		"claude-3-opus-20240229":   true,
		"claude-3-sonnet-20240229": true,
		"claude-3-haiku-20240307":  true,
		"gpt-4-vision-preview":     true,
		"gpt-4-turbo":              true,
		"gemini-pro-vision":        true,
	}

	return multimodalModels[model]
}
