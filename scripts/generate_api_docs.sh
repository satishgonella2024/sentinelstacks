#!/bin/bash

set -e

# ANSI color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}SentinelStacks API Documentation Generator${NC}"
echo "================================================="

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}Installing swaggo...${NC}"
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✓ swaggo installed${NC}"
else
    echo -e "${GREEN}✓ swaggo already installed${NC}"
fi

# Get the project root directory
ROOT_DIR=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
cd "$ROOT_DIR"

echo -e "${YELLOW}Generating OpenAPI specification...${NC}"

# Find the swag executable
SWAG_PATH=$(go env GOPATH)/bin/swag
if [ ! -f "$SWAG_PATH" ]; then
    echo -e "${RED}Could not find swag executable at $SWAG_PATH${NC}"
    echo -e "${RED}Please ensure it is installed correctly with: go install github.com/swaggo/swag/cmd/swag@latest${NC}"
    exit 1
fi

# Generate OpenAPI specification
"$SWAG_PATH" init -g internal/api/server.go -o docs/swagger --parseDependency --parseInternal

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ OpenAPI specification generated successfully${NC}"
    
    # Create the swagger directory if it doesn't exist
    mkdir -p docs/swagger
    
    # Copy the generated swagger.json to docs directory for mkdocs
    cp docs/swagger/swagger.json docs/swagger.json
    echo -e "${GREEN}✓ swagger.json copied to docs directory${NC}"
    
    # Update the swagger.html file if needed
    if [ ! -f docs/swagger.html ] || [ docs/swagger/index.html -nt docs/swagger.html ]; then
        cp docs/swagger/index.html docs/swagger.html
        echo -e "${GREEN}✓ swagger.html updated${NC}"
    fi
    
    echo -e "${GREEN}API Documentation is ready!${NC}"
    echo "You can view it in several ways:"
    echo "1. Start the API server: sentinel api"
    echo "   - Swagger UI: http://localhost:8081/swagger/"
    echo "   - ReDoc UI: http://localhost:8081/redoc/"
    echo "2. Build and view the mkdocs site: ./scripts/update_docs.sh serve"
else
    echo -e "${RED}✗ Failed to generate OpenAPI specification${NC}"
    exit 1
fi 