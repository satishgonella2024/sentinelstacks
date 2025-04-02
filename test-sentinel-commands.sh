#!/bin/bash

# Build the Sentinel application
echo "Building SentinelStacks..."
go build -o sentinel main.go

echo -e "\n\n===== TESTING NETWORK COMMANDS ====="
echo "Creating network 'test-network'..."
./sentinel network create test-network

echo -e "\nListing networks..."
./sentinel network ls

echo -e "\nConnecting agent 'agent1' to 'test-network'..."
./sentinel network connect test-network agent1

echo -e "\nConnecting agent 'agent2' to 'test-network'..."
./sentinel network connect test-network agent2

echo -e "\nInspecting network 'test-network'..."
./sentinel network inspect test-network

echo -e "\nDisconnecting agent 'agent2' from 'test-network'..."
./sentinel network disconnect test-network agent2

echo -e "\nInspecting network 'test-network' after disconnect..."
./sentinel network inspect test-network

echo -e "\n\n===== TESTING VOLUME COMMANDS ====="
echo "Creating volume 'test-volume'..."
./sentinel volume create test-volume --size 2GB --encrypted

echo -e "\nListing volumes..."
./sentinel volume ls

echo -e "\nMounting volume 'test-volume' to 'agent1'..."
./sentinel volume mount test-volume agent1 --path /memory

echo -e "\nInspecting volume 'test-volume'..."
./sentinel volume inspect test-volume

echo -e "\nUnmounting volume 'test-volume' from 'agent1'..."
./sentinel volume unmount test-volume agent1

echo -e "\nInspecting volume 'test-volume' after unmount..."
./sentinel volume inspect test-volume

echo -e "\n\n===== CLEANUP ====="
echo "Removing network 'test-network'..."
./sentinel network rm test-network

echo -e "\nRemoving volume 'test-volume'..."
./sentinel volume rm test-volume

echo -e "\nListing networks and volumes after cleanup..."
./sentinel network ls
./sentinel volume ls

echo -e "\n\nTest completed successfully!"
