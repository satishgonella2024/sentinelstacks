#!/bin/bash
echo "Making all helper scripts executable..."

# Make toggle mode script executable
chmod +x ./toggle-mode.sh

# Make run with backend script executable
chmod +x ./run-with-backend.sh

# Make web-ui scripts executable
chmod +x ./web-ui/toggle-api-mode.sh
chmod +x ./web-ui/fix-redux.sh
chmod +x ./web-ui/fix-styling.sh

# Make other scripts executable
chmod +x ./scripts/fix_dependencies.sh
chmod +x ./scripts/run_fixed.sh

echo "✅ All scripts are now executable!"
