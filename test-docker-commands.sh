#!/bin/bash

# Build the Sentinel application
echo "Building Sentinel..."
go build -o sentinel main.go

echo ""
echo "Testing network commands..."
./sentinel network create test-network
./sentinel network ls
./sentinel network connect test-network agent1
./sentinel network connect test-network agent2
./sentinel network inspect test-network
./sentinel network disconnect test-network agent2
./sentinel network rm test-network

echo ""
echo "Testing volume commands..."
./sentinel volume create test-volume --size 2GB --encrypted
./sentinel volume ls
./sentinel volume mount test-volume agent1 --path /memory
./sentinel volume inspect test-volume
./sentinel volume unmount test-volume agent1
./sentinel volume rm test-volume

echo ""
echo "Test completed!"
