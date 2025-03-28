package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/satishgonella2024/sentinelstacks/internal/registry"
	"gopkg.in/yaml.v3"
)

type Agent struct {
	Name         string    `json:"name" yaml:"name"`
	Version      string    `json:"version" yaml:"version"`
	Description  string    `json:"description" yaml:"description"`
	Capabilities []string  `json:"capabilities" yaml:"capabilities"`
	Commands     []Command `json:"commands" yaml:"commands"`
}

type Command struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Args        []Arg  `json:"args" yaml:"args"`
}

type Arg struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Required    bool   `json:"required" yaml:"required"`
	Description string `json:"description" yaml:"description"`
}

func StartServer(port string) error {
	http.HandleFunc("/api/agents", handleAgents)
	http.HandleFunc("/api/agents/", handleAgentDetails)

	fmt.Printf("Starting API server on port %s...\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agents, err := registry.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing agents: %v", err), http.StatusInternalServerError)
		return
	}

	// Get detailed information for each agent
	var detailedAgents []Agent
	for _, a := range agents {
		agent, err := getAgentDetails(a.Name, a.Version)
		if err != nil {
			fmt.Printf("Warning: Error getting details for agent %s:%s: %v\n", a.Name, a.Version, err)
			continue
		}
		detailedAgents = append(detailedAgents, agent)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(detailedAgents) == 0 {
		// Return empty array instead of null
		w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(detailedAgents); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleAgentDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract agent name and version from URL
	path := r.URL.Path[len("/api/agents/"):]
	agent, err := getAgentDetails(path, "latest")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}

func getAgentDetails(name, version string) (Agent, error) {
	agentDir := filepath.Join(os.Getenv("HOME"), ".sentinel", "agents", name, version)
	configFile := filepath.Join(agentDir, "agent.yaml")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Agent{}, fmt.Errorf("error reading agent config: %v", err)
	}

	var agent Agent
	if err := yaml.Unmarshal(data, &agent); err != nil {
		return Agent{}, fmt.Errorf("error parsing agent config: %v", err)
	}

	// Ensure required fields are set
	if agent.Name == "" {
		agent.Name = name
	}
	if agent.Version == "" {
		agent.Version = version
	}

	return agent, nil
}
