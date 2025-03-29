#!/bin/bash

# Make sure we're in the correct directory
cd "$(dirname "$0")"

# Install UUID for mock data
npm install --save uuid
npm install --save-dev @types/uuid

# Success message
echo "Dependencies installed. The application should now run with the mock data service."
echo "Run 'npm run dev' to start the development server."