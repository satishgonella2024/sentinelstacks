#!/bin/bash
set -e

cd "$(dirname "$0")"
echo "Starting fixed app with mock API enabled..."
VITE_USE_MOCK_API=true npm run dev 