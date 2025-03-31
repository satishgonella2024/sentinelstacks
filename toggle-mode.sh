#!/bin/bash

echo "🔄 SentinelStacks Mode Toggler"
echo "=============================="

# Check current status
if grep -q "USE_REAL_DATA: true" web-ui/src/api-config.ts; then
  echo "Currently using: Real Backend API"
  echo "Switching to: Mock Data mode"
  
  # Update to mock mode
  sed -i '' 's/USE_REAL_DATA: true/USE_REAL_DATA: false/' web-ui/src/api-config.ts
  echo "✅ Configuration updated to use mock data."
  echo ""
  echo "To run the application:"
  echo "  cd web-ui"
  echo "  npm run dev"
else
  echo "Currently using: Mock Data mode"
  echo "Switching to: Real Backend API"
  
  # Update to real mode
  sed -i '' 's/USE_REAL_DATA: false/USE_REAL_DATA: true/' web-ui/src/api-config.ts
  echo "✅ Configuration updated to use real backend API."
  echo ""
  echo "To run with the real backend:"
  echo "  ./run-with-backend.sh"
fi

echo ""
echo "Mode switch complete! 🎉"
