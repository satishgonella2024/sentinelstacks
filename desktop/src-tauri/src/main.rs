// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use serde::{Deserialize, Serialize};
use tauri::command;
use std::process::Command;
use anyhow::{Result, Context};
use std::fmt;

#[derive(Debug, Serialize)]
struct CommandError {
    message: String,
    details: Option<String>,
}

impl fmt::Display for CommandError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        if let Some(details) = &self.details {
            write!(f, "{}: {}", self.message, details)
        } else {
            write!(f, "{}", self.message)
        }
    }
}

impl From<anyhow::Error> for CommandError {
    fn from(err: anyhow::Error) -> Self {
        CommandError {
            message: err.to_string(),
            details: err.chain().skip(1).map(|e| e.to_string()).collect::<Vec<_>>().join(": "),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
struct Agent {
    id: String,
    name: String,
    description: String,
    status: String,
    model: Model,
    tools: Vec<Tool>,
    capabilities: Vec<String>,
    memory: Memory,
    created_at: String,
    updated_at: String,
    last_active_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Model {
    provider: String,
    name: String,
    options: ModelOptions,
}

#[derive(Debug, Serialize, Deserialize)]
struct ModelOptions {
    temperature: f32,
    max_tokens: i32,
}

#[derive(Debug, Serialize, Deserialize)]
struct Tool {
    id: String,
    name: String,
    description: String,
    version: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Memory {
    persistence: bool,
    vector_storage: bool,
    message_count: i32,
    last_updated: String,
}

fn run_sentinel_command(args: &[&str]) -> Result<String> {
    let command_str = format!("sentinel {}", args.join(" "));
    println!("Executing command: {}", command_str);

    let output = Command::new("sentinel")
        .args(args)
        .output()
        .with_context(|| format!("Failed to execute command: {}", command_str))?;

    if output.status.success() {
        String::from_utf8(output.stdout)
            .with_context(|| "Failed to parse command output as UTF-8")
    } else {
        let error_msg = String::from_utf8_lossy(&output.stderr);
        Err(anyhow::anyhow!("Command failed: {}", error_msg))
    }
}

#[command]
async fn get_agents() -> Result<Vec<Agent>, CommandError> {
    let output = run_sentinel_command(&["agent", "list", "--json"])
        .with_context(|| "Failed to list agents")?;
    
    serde_json::from_str(&output)
        .with_context(|| "Failed to parse agent list response")
        .map_err(Into::into)
}

#[command]
async fn get_agent(id: String) -> Result<Agent, CommandError> {
    let output = run_sentinel_command(&["agent", "get", &id, "--json"])
        .with_context(|| format!("Failed to get agent with ID: {}", id))?;
    
    serde_json::from_str(&output)
        .with_context(|| format!("Failed to parse agent response for ID: {}", id))
        .map_err(Into::into)
}

#[command]
async fn create_agent(agent: Agent) -> Result<Agent, CommandError> {
    let agent_json = serde_json::to_string(&agent)
        .with_context(|| "Failed to serialize agent data")?;
    
    let output = run_sentinel_command(&["agent", "create", "--json", &agent_json])
        .with_context(|| "Failed to create agent")?;
    
    serde_json::from_str(&output)
        .with_context(|| "Failed to parse created agent response")
        .map_err(Into::into)
}

#[command]
async fn update_agent(id: String, agent: Agent) -> Result<Agent, CommandError> {
    let agent_json = serde_json::to_string(&agent)
        .with_context(|| format!("Failed to serialize agent data for ID: {}", id))?;
    
    let output = run_sentinel_command(&["agent", "update", &id, "--json", &agent_json])
        .with_context(|| format!("Failed to update agent with ID: {}", id))?;
    
    serde_json::from_str(&output)
        .with_context(|| format!("Failed to parse updated agent response for ID: {}", id))
        .map_err(Into::into)
}

#[command]
async fn delete_agent(id: String) -> Result<(), CommandError> {
    run_sentinel_command(&["agent", "delete", &id])
        .with_context(|| format!("Failed to delete agent with ID: {}", id))
        .map(|_| ())
        .map_err(Into::into)
}

#[command]
async fn start_agent(id: String) -> Result<(), CommandError> {
    run_sentinel_command(&["agent", "start", &id])
        .with_context(|| format!("Failed to start agent with ID: {}", id))
        .map(|_| ())
        .map_err(Into::into)
}

#[command]
async fn stop_agent(id: String) -> Result<(), CommandError> {
    run_sentinel_command(&["agent", "stop", &id])
        .with_context(|| format!("Failed to stop agent with ID: {}", id))
        .map(|_| ())
        .map_err(Into::into)
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            get_agents,
            get_agent,
            create_agent,
            update_agent,
            delete_agent,
            start_agent,
            stop_agent,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
