#!/bin/bash
set -e

echo "🔧 Fixing SentinelStacks UI issues..."

# Generate MSW service worker
echo "📝 Setting up Mock Service Worker..."
npx msw init public/ --save

# Run the application
echo "🚀 Starting the application..."
cd ..
./scripts/dev.sh
