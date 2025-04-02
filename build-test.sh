#!/bin/bash
echo "Building and testing SentinelStacks..."
./build-clean.sh
chmod +x ./sentinel
echo "Running agent test..."
./sentinel run myname/chatbot:latest
