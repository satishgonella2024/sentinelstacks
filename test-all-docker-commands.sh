#!/bin/bash

# Test script for SentinelStacks Docker-inspired commands
# This script tests all the Docker-inspired commands implemented in SentinelStacks

set -e  # Exit on error

COLOR_RED="\033[0;31m"
COLOR_GREEN="\033[0;32m"
COLOR_YELLOW="\033[0;33m"
COLOR_BLUE="\033[0;34m"
COLOR_RESET="\033[0m"

echo -e "${COLOR_BLUE}=== SENTINELSTACKS DOCKER-INSPIRED COMMANDS TEST SCRIPT ===${COLOR_RESET}"

# Function to display section header
section() {
    echo -e "\n${COLOR_YELLOW}=== $1 ===${COLOR_RESET}"
}

# Function to run a command and display its output
run_command() {
    echo -e "${COLOR_BLUE}$ $1${COLOR_RESET}"
    eval "$1"
    if [ $? -eq 0 ]; then
        echo -e "${COLOR_GREEN}Command succeeded${COLOR_RESET}"
    else
        echo -e "${COLOR_RED}Command failed${COLOR_RESET}"
        exit 1
    fi
}

# Fix directory structure
section "Setting up environment"
run_command "sh fix-directory-structure.sh"

# Build the application
section "Building SentinelStacks"
run_command "go build -o sentinel main.go"

# Create data directories if they don't exist
mkdir -p ~/.sentinel/data/networks
mkdir -p ~/.sentinel/data/volumes
mkdir -p ~/.sentinel/data/systems

# Check the version
section "Checking version"
run_command "./sentinel version"

# Test network commands
section "Network Commands"
run_command "./sentinel network create test-network-1"
run_command "./sentinel network create test-network-2 --driver advanced"
run_command "./sentinel network ls"
run_command "./sentinel network connect test-network-1 test-agent-1"
run_command "./sentinel network connect test-network-1 test-agent-2"
run_command "./sentinel network inspect test-network-1"
run_command "./sentinel network disconnect test-network-1 test-agent-2"
run_command "./sentinel network inspect test-network-1"

# Test volume commands
section "Volume Commands"
run_command "./sentinel volume create test-volume-1"
run_command "./sentinel volume create test-volume-2 --size 2GB"
run_command "./sentinel volume create test-volume-3 --size 1GB --encrypted"
run_command "./sentinel volume ls"
run_command "./sentinel volume mount test-volume-1 test-agent-1"
run_command "./sentinel volume mount test-volume-2 test-agent-1 --path /custom/memory"
run_command "./sentinel volume inspect test-volume-1"
run_command "./sentinel volume inspect test-volume-2"
run_command "./sentinel volume unmount test-volume-1 test-agent-1"
run_command "./sentinel volume unmount test-volume-2 test-agent-1"

# Test compose commands
section "Compose Commands"
run_command "cp examples/compose-example.yaml ."
run_command "./sentinel compose up -f compose-example.yaml"
run_command "./sentinel compose ls"
# Get the system ID from the list
SYSTEM_ID=$(./sentinel compose ls | grep -v "ID" | awk '{print $1}' | head -n 1)
if [ -n "$SYSTEM_ID" ]; then
    run_command "./sentinel compose pause $SYSTEM_ID"
    run_command "./sentinel compose resume $SYSTEM_ID"
    run_command "./sentinel compose logs $SYSTEM_ID"
else
    echo -e "${COLOR_RED}No system ID found, skipping system-specific commands${COLOR_RESET}"
fi

# Test system commands
section "System Commands"
run_command "./sentinel system info"
run_command "./sentinel system df"
run_command "./sentinel system events --limit 5"

# Clean up
section "Cleaning up resources"
if [ -n "$SYSTEM_ID" ]; then
    run_command "./sentinel compose down $SYSTEM_ID"
fi
run_command "./sentinel volume rm test-volume-1"
run_command "./sentinel volume rm test-volume-2"
run_command "./sentinel volume rm test-volume-3"
run_command "./sentinel network rm test-network-1"
run_command "./sentinel network rm test-network-2"

# Verify cleanup
section "Verifying cleanup"
run_command "./sentinel network ls"
run_command "./sentinel volume ls"
run_command "./sentinel compose ls"
run_command "./sentinel system info"

echo -e "\n${COLOR_GREEN}All tests completed successfully!${COLOR_RESET}"
