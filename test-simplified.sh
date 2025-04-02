#!/bin/bash

# Test script for SentinelStacks Docker-inspired commands with simplified implementation
# This script tests the basic functionality of the commands

echo "=== SENTINELSTACKS SIMPLIFIED TEST SCRIPT ==="

# Build the application
echo "Building SentinelStacks..."
go build -o sentinel main.go

# Test network commands
echo -e "\n\n=== Testing Network Commands ==="
echo "Creating network..."
./sentinel network create test-network

echo "Listing networks..."
./sentinel network ls

echo "Connecting agent to network..."
./sentinel network connect test-network agent1

echo "Inspecting network..."
./sentinel network inspect test-network

# Test volume commands
echo -e "\n\n=== Testing Volume Commands ==="
echo "Creating volume..."
./sentinel volume create test-volume --size 2GB

echo "Listing volumes..."
./sentinel volume ls

echo "Mounting volume..."
./sentinel volume mount test-volume agent1

echo "Inspecting volume..."
./sentinel volume inspect test-volume

# Test compose commands
echo -e "\n\n=== Testing Compose Commands ==="
echo "Creating example compose file..."
cat > example-compose.json << EOF
{
  "name": "test-system",
  "networks": {
    "default": {
      "driver": "default"
    }
  },
  "volumes": {
    "data": {
      "size": "1GB",
      "encrypted": false
    }
  },
  "agents": {
    "agent1": {
      "image": "agent:latest",
      "networks": ["default"],
      "volumes": ["data:/memory"],
      "environment": {
        "ROLE": "default"
      },
      "resources": {
        "memory": "1GB"
      }
    }
  }
}
EOF

echo "Starting multi-agent system..."
./sentinel compose up -f example-compose.json

echo "Listing multi-agent systems..."
./sentinel compose ls

# Test system commands
echo -e "\n\n=== Testing System Commands ==="
echo "Displaying system information..."
./sentinel system info

echo "All tests completed!"
