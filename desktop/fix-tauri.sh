#!/bin/bash

# Make sure we're in the correct directory
cd "$(dirname "$0")"

echo "Fixing Tauri configuration for SentinelStacks Desktop..."

# Make sure the Tauri CLI is installed
if ! command -v cargo &> /dev/null || ! command -v rustc &> /dev/null; then
  echo "Rust and/or Cargo not found. Please install Rust first:"
  echo "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
  exit 1
fi

# Check if Tauri CLI is installed
if ! command -v cargo-tauri &> /dev/null; then
  echo "Installing Tauri CLI version 2.0.0..."
  cargo install tauri-cli --version "^2.0.0"
fi

# Check for Tauri plugins
echo "Installing required Tauri plugins..."
cd src-tauri || { echo "src-tauri directory not found!"; exit 1; }
cargo add tauri@2.0.0
cargo add tauri-build@2.0.0 --build
cargo add tauri-plugin-shell@2.0.0
cargo add tauri-plugin-dialog@2.0.0
cargo add tauri-plugin-fs@2.0.0
cargo add tauri-plugin-process@2.0.0

echo "Tauri setup fixed! You can now run 'npm run tauri dev' to start the desktop application."
