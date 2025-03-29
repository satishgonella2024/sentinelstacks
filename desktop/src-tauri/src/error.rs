use serde::Serialize;
use thiserror::Error;

#[derive(Error, Debug, Serialize)]
pub enum AppError {
    #[error("CLI error: {0}")]
    Cli(String),
    
    #[error("Command execution failed: {0}")]
    CommandExecution(String),
    
    #[error("Failed to parse JSON: {0}")]
    JsonParse(String),
    
    #[error("Agent not found: {0}")]
    AgentNotFound(String),
    
    #[error("Registry error: {0}")]
    Registry(String),
    
    #[error("IO error: {0}")]
    Io(String),
    
    #[error("Invalid input: {0}")]
    InvalidInput(String),
}

impl From<std::io::Error> for AppError {
    fn from(err: std::io::Error) -> Self {
        AppError::Io(err.to_string())
    }
}

impl From<serde_json::Error> for AppError {
    fn from(err: serde_json::Error) -> Self {
        AppError::JsonParse(err.to_string())
    }
}