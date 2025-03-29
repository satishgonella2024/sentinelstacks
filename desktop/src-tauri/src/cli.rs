use crate::error::AppError;
use serde::{Deserialize, Serialize};
use std::process::Command;
use std::time::SystemTime;

#[derive(Debug, Serialize, Deserialize)]
pub struct Agent {
    pub name: String,
    pub description: String,
    pub model: String,
    pub memory_type: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentDetails {
    pub name: String,
    pub version: String,
    pub description: String,
    pub model: AgentModel,
    pub capabilities: Vec<String>,
    pub memory: AgentMemory,
    pub tools: Vec<AgentTool>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentModel {
    pub provider: String,
    pub name: String,
    pub endpoint: Option<String>,
    pub options: Option<serde_json::Value>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentMemory {
    pub r#type: String,
    pub persistence: bool,
    pub max_items: Option<u32>,
    pub embedding_model: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentTool {
    pub id: String,
    pub version: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RegistryAgent {
    pub name: String,
    pub description: String,
    pub version: String,
    pub author: Option<String>,
    pub tags: Option<Vec<String>>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MemoryEntry {
    pub id: String,
    pub content: String,
    pub timestamp: String,
    pub metadata: Option<serde_json::Value>,
}

pub fn get_agents() -> Result<Vec<Agent>, AppError> {
    let output = Command::new("sentinel")
        .args(&["agent", "list", "--format", "json"])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let agents: Vec<Agent> = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(agents)
}

pub fn create_agent(name: String, description: String, model: String, memory_type: String) -> Result<(), AppError> {
    // Validate inputs
    if name.is_empty() {
        return Err(AppError::InvalidInput("Agent name cannot be empty".to_string()));
    }
    
    let mut args = vec![
        "agent".to_string(),
        "create".to_string(),
        "--name".to_string(),
        name,
        "--description".to_string(),
        description,
    ];
    
    // Add model if provided
    if !model.is_empty() {
        args.push("--model".to_string());
        args.push(model);
    }
    
    // Add memory type if provided
    if !memory_type.is_empty() {
        args.push("--memory".to_string());
        args.push(memory_type);
    }
    
    let output = Command::new("sentinel")
        .args(&args)
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    Ok(())
}

pub fn run_agent(name: String, input: String) -> Result<String, AppError> {
    let output = Command::new("sentinel")
        .args(&["agent", "run", "--name", &name, "--input", &input, "--format", "json"])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let response = String::from_utf8_lossy(&output.stdout).to_string();
    Ok(response)
}

pub fn get_agent_details(name: String) -> Result<AgentDetails, AppError> {
    let output = Command::new("sentinel")
        .args(&["agent", "inspect", "--name", &name, "--format", "json"])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let details: AgentDetails = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(details)
}

pub fn search_registry(query: String) -> Result<Vec<RegistryAgent>, AppError> {
    let output = Command::new("sentinel")
        .args(&["registry", "search", "--query", &query, "--format", "json"])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let agents: Vec<RegistryAgent> = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(agents)
}

pub fn get_registry_agents() -> Result<Vec<RegistryAgent>, AppError> {
    let output = Command::new("sentinel")
        .args(&["registry", "list", "--format", "json"])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let agents: Vec<RegistryAgent> = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(agents)
}

pub fn get_memory(agent_name: String, limit: u32) -> Result<Vec<MemoryEntry>, AppError> {
    let limit_str = limit.to_string();
    let output = Command::new("sentinel")
        .args(&[
            "memory",
            "list",
            "--agent",
            &agent_name,
            "--limit",
            &limit_str,
            "--format",
            "json",
        ])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let entries: Vec<MemoryEntry> = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(entries)
}

pub fn search_memory(agent_name: String, query: String, limit: u32) -> Result<Vec<MemoryEntry>, AppError> {
    let limit_str = limit.to_string();
    let output = Command::new("sentinel")
        .args(&[
            "memory",
            "search",
            "--agent",
            &agent_name,
            "--query",
            &query,
            "--limit",
            &limit_str,
            "--format",
            "json",
        ])
        .output()
        .map_err(|e| AppError::CommandExecution(e.to_string()))?;
    
    if !output.status.success() {
        return Err(AppError::Cli(
            String::from_utf8_lossy(&output.stderr).to_string(),
        ));
    }
    
    let json_str = String::from_utf8_lossy(&output.stdout);
    let entries: Vec<MemoryEntry> = serde_json::from_str(&json_str)
        .map_err(|e| AppError::JsonParse(e.to_string()))?;
    
    Ok(entries)
}