#!/bin/bash

# Toggle between real API and mock data

# Check current status
if grep -q "USE_REAL_DATA: true" src/api-config.ts; then
  echo "🔄 Switching to mock data mode..."
  # Update to mock mode
  sed -i '' 's/USE_REAL_DATA: true/USE_REAL_DATA: false/' src/api-config.ts
  echo "✅ Now using mock data. No backend server needed."
else
  echo "🔄 Switching to real API mode..."
  # Update to real mode
  sed -i '' 's/USE_REAL_DATA: false/USE_REAL_DATA: true/' src/api-config.ts
  echo "✅ Now using real API. Make sure the backend server is running on port 8080."
fi

# Show current status
if grep -q "USE_REAL_DATA: true" src/api-config.ts; then
  echo "📊 Current mode: Real API (backend required)"
else
  echo "📊 Current mode: Mock Data (standalone)"
fi
