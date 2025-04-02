#!/bin/bash

echo "=== TESTING SENTINELSTACKS DOCKER-INSPIRED COMMANDS ==="

# First, build the application
echo -e "\n\n[STEP 1] Building SentinelStacks..."
go build -o sentinel main.go

# Create the home directory structure
echo -e "\n\n[STEP 2] Setting up test environment..."
mkdir -p ./testdata/networks
mkdir -p ./testdata/volumes
mkdir -p ./testdata/systems

# Test Network Commands
echo -e "\n\n[STEP 3] Testing Network Commands..."
echo "Creating network 'agent-net'..."
./sentinel network create agent-net
echo "Creating network 'data-net'..."
./sentinel network create data-net
echo "Listing networks..."
./sentinel network ls
echo "Connecting an agent to the network..."
./sentinel network connect agent-net agent1
echo "Inspecting network..."
./sentinel network inspect agent-net

# Test Volume Commands
echo -e "\n\n[STEP 4] Testing Volume Commands..."
echo "Creating volume 'memory-vol'..."
./sentinel volume create memory-vol --size 2GB
echo "Creating encrypted volume 'secure-vol'..."
./sentinel volume create secure-vol --size 1GB --encrypted
echo "Listing volumes..."
./sentinel volume ls
echo "Mounting volume to an agent..."
./sentinel volume mount memory-vol agent1 --path /memory
echo "Inspecting volume..."
./sentinel volume inspect memory-vol

# Test Multi-Agent System (Compose) Commands
echo -e "\n\n[STEP 5] Testing Compose Commands..."
echo "Starting a multi-agent system from a compose file..."
./sentinel compose up -f examples/compose-example.yaml
echo "Listing multi-agent systems..."
./sentinel compose ls

# Test System Commands
echo -e "\n\n[STEP 6] Testing System Commands..."
echo "Displaying system information..."
./sentinel system info
echo "Showing disk usage..."
./sentinel system df
echo "Viewing system events..."
./sentinel system events --limit 5

# Clean up
echo -e "\n\n[STEP 7] Cleaning up resources..."
echo "Unmounting volumes..."
./sentinel volume unmount memory-vol agent1
echo "Removing volumes..."
./sentinel volume rm memory-vol
./sentinel volume rm secure-vol
echo "Disconnecting agents from networks..."
./sentinel network disconnect agent-net agent1
echo "Removing networks..."
./sentinel network rm agent-net
./sentinel network rm data-net

echo -e "\n\n[STEP 8] Verifying cleanup..."
echo "Listing networks..."
./sentinel network ls
echo "Listing volumes..."
./sentinel volume ls
echo "Listing multi-agent systems..."
./sentinel compose ls

echo -e "\n\n=== TEST COMPLETE ==="
