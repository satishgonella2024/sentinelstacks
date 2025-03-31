#!/bin/bash
set -e

echo "🚀 Starting SentinelStacks Enhanced Application..."

# Navigate to project root
cd "$(dirname "$0")"

# Fix Go dependencies
echo "👷 Fixing dependencies..."
go get github.com/mattn/go-isatty@v0.0.20
go mod tidy

# Set up mock mode for easier development
export SENTINEL_MOCK_MODE=true

# Create directories for MSW
echo "📝 Setting up Mock Service Worker..."
mkdir -p web-ui/public

# Initialize MSW service worker if needed
if [ ! -f web-ui/public/mockServiceWorker.js ]; then
  cd web-ui
  npx msw init public/ --save
  cd ..
fi

# Configure application for mock mode
echo "🔧 Configuring for mock mode..."
cat > web-ui/src/api-config.ts << EOF
// API configuration
export const API_CONFIG = {
  // Set to false to use mock data
  USE_REAL_DATA: false,
  
  // Base URL for API requests
  API_BASE_URL: '/v1',
  
  // Timeout for API requests in milliseconds
  TIMEOUT: 10000,
  
  // Maximum number of retries for failed requests
  MAX_RETRIES: 3
};
EOF

# Create a style fix for the background
echo "🎨 Applying styling fixes..."
cat > web-ui/public/additional-styles.css << 'EOF'
body {
  background-color: #0a1017;
  color: white;
}

/* Additional styles for think bubbles */
.think-bubble {
  background-color: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(71, 85, 105, 0.3);
  border-radius: 0.5rem;
  backdrop-filter: blur(8px);
}

.think-bubble.insight {
  background-color: rgba(20, 83, 153, 0.2);
  border-color: rgba(37, 99, 235, 0.4);
}

.think-bubble.guidance {
  background-color: rgba(21, 128, 61, 0.2);
  border-color: rgba(34, 197, 94, 0.4);
}

.think-bubble.suggestion {
  background-color: rgba(124, 58, 237, 0.2);
  border-color: rgba(139, 92, 246, 0.4);
}

.think-bubble.achievement {
  background-color: rgba(180, 83, 9, 0.2);
  border-color: rgba(245, 158, 11, 0.4);
}
EOF

# Start the frontend
echo "🖥️ Starting the enhanced application..."
cd web-ui
npm run dev
