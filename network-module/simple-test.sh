#!/bin/bash
set -e

# Clean any previous build
rm -f sentinel-network

# Build the binary
echo "Building..."
go mod tidy
go build -o sentinel-network .
echo "Built successfully!"

# Run a simple test
echo -e "\nCreating network..."
./sentinel-network network create simple-network

echo -e "\nListing networks..."
./sentinel-network network ls

echo -e "\nTest completed!"

