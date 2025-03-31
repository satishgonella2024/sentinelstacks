#!/bin/bash
set -e

echo "🚀 Starting SentinelStacks with real backend API..."

# Build the Go backend
echo "👷 Building the backend..."
go mod tidy
go build -o bin/sentinel main.go

# Run backend in background
echo "🔌 Starting the backend API server..."
./bin/sentinel server start --port 8080 &
BACKEND_PID=$!

# Wait a bit for the backend to initialize
sleep 2

# Set up frontend for real data
cd web-ui

# Make sure the frontend uses real API data
echo "🔧 Configuring frontend for real API..."
cat > src/api-config.ts << EOF
// API configuration
export const API_CONFIG = {
  USE_REAL_DATA: true,
  API_BASE_URL: '/v1'
};
EOF

# Start the frontend
echo "🖥️ Starting the frontend..."
npm run dev

# Clean up backend when frontend is stopped
kill $BACKEND_PID
