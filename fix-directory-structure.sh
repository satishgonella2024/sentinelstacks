#!/bin/bash

echo "Fixing directory structure for SentinelStacks..."

# Fix multi_agent_system_repository
echo "Fixing multi-agent system repository..."
rm -rf ./pkg/repository/fs/multi_agent_system_repository.go
cp ./pkg/repository/fs/tmp/mas_repo.go ./pkg/repository/fs/multi_agent_system_repository.go

# Clean up temporary directories
echo "Cleaning up temporary directories..."
rm -rf ./pkg/repository/fs/tmp
rm -rf ./pkg/repository/fs/temp

echo "Directory structure fixed!"
