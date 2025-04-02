#!/bin/bash

# This script prepares and performs a check-in

echo "Preparing for check-in..."

# Make scripts executable
chmod +x scripts/*.sh

# Fix imports
echo "Fixing imports..."
./scripts/fix_imports.sh

# Verify the project builds
echo "Verifying build..."
go build -o sentinel-test ./main.go
if [ $? -ne 0 ]; then
    echo "Build failed! Aborting check-in."
    exit 1
fi
rm -f sentinel-test

# Show status
echo "Current git status:"
git status --porcelain

# Suggested commit message
echo "
Ready for check-in!

Suggested commit message:
-----------------------
feat: Implement stack engine with memory management

- Renamed 'compose' to 'stack' with full functionality
- Added complete memory management system with multiple backends
- Implemented flexible agent runtime system
- Added CLI commands for managing memory
- Updated documentation with implementation progress
"

# Remind about git commands
echo "
To commit and push, use:
-----------------------
git add .
git commit -m \"feat: Implement stack engine with memory management\"
git push origin main
"

echo "Check-in preparation complete!"
