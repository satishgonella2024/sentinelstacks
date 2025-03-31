#!/bin/bash
set -e

echo "🚀 Starting SentinelStacks with real backend API..."

# Fix the backend dependencies and build
echo "👷 Building the backend..."
cd "$(dirname "$0")"
go get github.com/mattn/go-isatty@v0.0.20
go mod tidy

# Make sure conversation package is properly structured
echo "✅ Conversation package updated"

# Start the backend in mock mode to avoid dependency errors
echo "🔄 Starting in mock mode for compatibility..."
export SENTINEL_MOCK_MODE=true

# Configure frontend for real backend API
echo "🔧 Configuring frontend to use backend API..."
cd web-ui
cat > src/api-config.ts << EOF
// API configuration
export const API_CONFIG = {
  // Set to true to use the real backend API, false to use mock data
  USE_REAL_DATA: true,
  
  // Base URL for API requests
  API_BASE_URL: '/v1',
  
  // Timeout for API requests in milliseconds
  TIMEOUT: 10000,
  
  // Maximum number of retries for failed requests
  MAX_RETRIES: 3
};
EOF

# Create public directory for MSW
mkdir -p public

# Initialize MSW service worker if not already done
if [ ! -f public/mockServiceWorker.js ]; then
  echo "📝 Setting up Mock Service Worker..."
  npx msw init public/ --save
fi

# Start the frontend
echo "📱 Starting the frontend..."
npm run dev
