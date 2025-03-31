#!/bin/bash
set -e

echo "🛠️ Fixing Redux and API configuration..."

# Create public directory for MSW
mkdir -p public

# Initialize MSW service worker
echo "📝 Setting up Mock Service Worker..."
npx msw init public/ --save

# Run the application
echo "🚀 Starting the application..."
cd ..
./scripts/dev.sh
