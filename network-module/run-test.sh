#!/bin/bash

# Simple script to build and test the network module
set -e

# Clean up any previous build
rm -f sentinel-network

# Build the binary
echo "Building network module..."
go mod tidy
go build -o sentinel-network .

# Run simple tests
echo -e "\nTesting network create..."
./sentinel-network network create test-network-1 --formats=text,image,audio

echo -e "\nTesting network list..."
./sentinel-network network ls

echo -e "\nTesting network inspect..."
./sentinel-network network inspect test-network-1

echo -e "\nTesting agent connection..."
./sentinel-network network connect test-network-1 agent1
./sentinel-network network connect test-network-1 agent2

echo -e "\nTesting network remove..."
./sentinel-network network rm test-network-1 --force

echo -e "\nTest completed successfully!"

# Clean up
rm -f sentinel-network
rm -rf ~/.sentinel/data
