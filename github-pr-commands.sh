#!/bin/bash
# Commands to create a new branch and push changes for the tool support feature

# Create a new branch
git checkout -b feature/tool-support

# Add new files
git add pkg/tools/
git add examples/tool-agent/
git add .github/PULL_REQUEST_TEMPLATE/tool-support.md

# Commit the new files
git commit -m "Add tool support infrastructure with calculator and weather tools"

# Add modified files
git add pkg/runtime/agent.go 
git add pkg/agentfile/schema.go
git add ROADMAP.md
git add NEXT_STEPS.md

# Commit the modified files
git commit -m "Update agent runtime and documentation to support tools"

# Push to remote
git push -u origin feature/tool-support

echo "Now create a pull request from the feature/tool-support branch"
echo "Use the pull request template at .github/PULL_REQUEST_TEMPLATE/tool-support.md"
