#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

mod cli;
mod error;

use error::AppError;
use std::process::Command;
use std::sync::Mutex;
use tauri::State;

// Agent state to be managed by Tauri
struct AgentState(Mutex<Vec<String>>);

fn main() {
    tauri::Builder::default()
        .manage(AgentState(Mutex::new(Vec::new())))
        .invoke_handler(tauri::generate_handler![
            get_agents,
            create_agent,
            run_agent,
            get_agent_details,
            search_registry,
            get_registry_agents,
            get_memory,
            search_memory
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
fn get_agents() -> Result<Vec<cli::Agent>, AppError> {
    cli::get_agents()
}

#[tauri::command]
fn create_agent(name: String, description: String, model: String, memory_type: String) -> Result<(), AppError> {
    cli::create_agent(name, description, model, memory_type)
}

#[tauri::command]
fn run_agent(name: String, input: String) -> Result<String, AppError> {
    cli::run_agent(name, input)
}

#[tauri::command]
fn get_agent_details(name: String) -> Result<cli::AgentDetails, AppError> {
    cli::get_agent_details(name)
}

#[tauri::command]
fn search_registry(query: String) -> Result<Vec<cli::RegistryAgent>, AppError> {
    cli::search_registry(query)
}

#[tauri::command]
fn get_registry_agents() -> Result<Vec<cli::RegistryAgent>, AppError> {
    cli::get_registry_agents()
}

#[tauri::command]
fn get_memory(agent_name: String, limit: u32) -> Result<Vec<cli::MemoryEntry>, AppError> {
    cli::get_memory(agent_name, limit)
}

#[tauri::command]
fn search_memory(agent_name: String, query: String, limit: u32) -> Result<Vec<cli::MemoryEntry>, AppError> {
    cli::search_memory(agent_name, query, limit)
}