#!/bin/bash

set -e

# ANSI color codes for better readability
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

PORT=8082
DOCS_DIR="public/api"

# Check if the documentation exists
if [ ! -f "$DOCS_DIR/swagger.yaml" ]; then
  echo -e "${RED}Error: API documentation not found at $DOCS_DIR/swagger.yaml${NC}"
  echo -e "${YELLOW}Run ./scripts/generate_api_docs.sh first to generate the documentation${NC}"
  exit 1
fi

# Check for Python to run a simple HTTP server
if command -v python3 &> /dev/null; then
  echo -e "${GREEN}Starting API documentation server on http://localhost:$PORT ${NC}"
  echo -e "${BLUE}Press Ctrl+C to stop the server${NC}"
  cd "$(dirname "$DOCS_DIR")"
  python3 -m http.server $PORT
elif command -v python &> /dev/null; then
  echo -e "${GREEN}Starting API documentation server on http://localhost:$PORT ${NC}"
  echo -e "${BLUE}Press Ctrl+C to stop the server${NC}"
  cd "$(dirname "$DOCS_DIR")"
  python -m SimpleHTTPServer $PORT
else
  echo -e "${RED}Error: Python not found. Unable to start an HTTP server.${NC}"
  echo -e "${YELLOW}Install Python or use another method to serve the files from $DOCS_DIR${NC}"
  exit 1
fi 