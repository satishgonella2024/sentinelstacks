#!/bin/bash

# Script to run the simple test
set -e  # Exit on error

echo "Running simple DAG test..."

# Build and run the test
cd $(dirname $0)
go run simple_dag.go

echo "Simple test completed."
