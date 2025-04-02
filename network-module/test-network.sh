#!/bin/bash

# Robust test script for the network module
# This script tests all network-related commands including multimodal support

# Exit on error
set -e

# Set colors for better output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== NETWORK MODULE TEST SCRIPT ===${NC}"

# Find the correct directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
cd "$SCRIPT_DIR"

# Clean up any previous build
rm -f sentinel-network

# Build the module
echo -e "\n${YELLOW}Building network module...${NC}"

# Download dependencies if needed
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to download dependencies!${NC}"
    exit 1
fi

# Build the binary
echo -e "${YELLOW}Building binary...${NC}"
go build -o sentinel-network
if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
echo -e "${GREEN}Build successful!${NC}"

# Set test data directory
TEST_DATA_DIR="$SCRIPT_DIR/test-data"
mkdir -p "$TEST_DATA_DIR/networks"
mkdir -p "$TEST_DATA_DIR/messages"
mkdir -p "$TEST_DATA_DIR/attachments"

# Create sample test files for attachments
echo "This is a sample text file" > "$TEST_DATA_DIR/sample.txt"
echo '{"key": "value", "array": [1, 2, 3]}' > "$TEST_DATA_DIR/sample.json"
echo "<svg><rect width='100' height='100' /></svg>" > "$TEST_DATA_DIR/sample.svg"

# Run network commands
echo -e "\n${YELLOW}=== Testing Network Management ===${NC}"

echo -e "\n${YELLOW}1. Creating networks with multimodal support...${NC}"
./sentinel-network network create test-network-1 --formats=text,image,audio
./sentinel-network network create test-network-2 --driver advanced --formats=text,json --config='{"max_message_size": 10485760, "encryption": true}'

echo -e "\n${YELLOW}2. Listing networks...${NC}"
./sentinel-network network ls

echo -e "\n${YELLOW}3. Inspecting network with multimodal support...${NC}"
./sentinel-network network inspect test-network-2

echo -e "\n${YELLOW}4. Connecting agents to networks...${NC}"
./sentinel-network network connect test-network-1 agent1
./sentinel-network network connect test-network-1 agent2
./sentinel-network network connect test-network-2 agent3

echo -e "\n${YELLOW}5. Updating network configuration...${NC}"
./sentinel-network network config test-network-1 --config='{"supported_formats": ["text", "image", "audio", "video"], "max_message_size": 20971520}'

echo -e "\n${YELLOW}6. Inspecting updated network...${NC}"
./sentinel-network network inspect test-network-1

echo -e "\n${YELLOW}=== Testing Multimodal Messaging ===${NC}"

echo -e "\n${YELLOW}1. Sending text message...${NC}"
./sentinel-network network message send test-network-1 agent1 --content="Hello, this is a test message" --format=text

echo -e "\n${YELLOW}2. Sending message with attachment...${NC}"
./sentinel-network network message send test-network-1 agent2 --content="Here's a file" --format=text --attach="text:$TEST_DATA_DIR/sample.txt" --metadata='{"priority": "high"}'

echo -e "\n${YELLOW}3. Sending JSON message...${NC}"
./sentinel-network network message send test-network-2 agent3 --content="$TEST_DATA_DIR/sample.json" --format=json

echo -e "\n${YELLOW}4. Listing messages...${NC}"
./sentinel-network network message ls test-network-1

echo -e "\n${YELLOW}5. Getting message details...${NC}"
# Get the first message ID (this is a simplification - in real use, you'd use a specific message ID)
MSG_ID=$(./sentinel-network network message ls test-network-1 | grep -v "^ID" | grep -v "^--" | head -1 | awk '{print $1}')
if [ -n "$MSG_ID" ]; then
    ./sentinel-network network message get test-network-1 "$MSG_ID"
fi

echo -e "\n${YELLOW}=== Testing Cleanup ===${NC}"

echo -e "\n${YELLOW}1. Disconnecting agents...${NC}"
./sentinel-network network disconnect test-network-1 agent2

echo -e "\n${YELLOW}2. Removing networks...${NC}"
# Should fail as there are still connected agents
./sentinel-network network rm test-network-1 || echo -e "${GREEN}Expected failure with connected agent${NC}"

# Force remove
./sentinel-network network rm test-network-1 --force
./sentinel-network network rm test-network-2

echo -e "\n${YELLOW}3. Verifying networks are removed...${NC}"
./sentinel-network network ls

# Clean up test data
rm -rf "$TEST_DATA_DIR"
rm -f sentinel-network

echo -e "\n${GREEN}Test completed successfully!${NC}"
