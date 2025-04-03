#!/bin/bash

set -e

# ANSI color codes for better readability
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

echo -e "${YELLOW}=== SentinelStacks API Documentation Generator ===${NC}"
echo -e "${BLUE}Generating comprehensive API documentation...${NC}"

# Ensure required directories exist
mkdir -p docs/api
mkdir -p public/api

# Check optional dependencies
check_optional_dependency() {
  if ! command -v $1 &> /dev/null; then
    echo -e "${YELLOW}Warning: $1 not found. Some features will be disabled.${NC}"
    if [ "$1" = "redoc-cli" ]; then
      echo -e "${YELLOW}You can install it with: npm i -g redoc-cli${NC}"
    elif [ "$1" = "jq" ]; then
      echo -e "${YELLOW}You can install it with: brew install jq${NC}"
    fi
    return 1
  fi
  return 0
}

# Check optional dependencies
HAVE_REDOC=$(check_optional_dependency "redoc-cli" && echo "yes" || echo "no")
HAVE_JQ=$(check_optional_dependency "jq" && echo "yes" || echo "no")
HAVE_PYTHON=$(command -v python3 &> /dev/null && echo "yes" || echo "no")

# Paths
API_REFERENCE_FILE="docs/api-reference.yaml"
SWAGGER_JSON="public/api/swagger.json"
SWAGGER_YAML="public/api/swagger.yaml"
REDOC_HTML="public/api/redoc.html"

# Process the OpenAPI YAML file
if [ -f "$API_REFERENCE_FILE" ]; then
  echo -e "${BLUE}Found OpenAPI spec file, using as reference...${NC}"
  
  # Create directories if they don't exist
  mkdir -p $(dirname "$SWAGGER_JSON")
  mkdir -p $(dirname "$SWAGGER_YAML")
  
  # Copy the YAML file
  cp "$API_REFERENCE_FILE" "$SWAGGER_YAML"
  echo -e "${GREEN}OpenAPI spec copied to $SWAGGER_YAML${NC}"
  
  # Convert to JSON if possible
  if command -v yq &> /dev/null; then
    yq -o=json "$API_REFERENCE_FILE" > "$SWAGGER_JSON"
    echo -e "${GREEN}OpenAPI spec converted to JSON format${NC}"
  elif [ "$HAVE_PYTHON" = "yes" ]; then
    # Try with Python if yq isn't available
    echo -e "${BLUE}Converting YAML to JSON using Python...${NC}"
    python3 -c "import sys, yaml, json; json.dump(yaml.safe_load(open('$API_REFERENCE_FILE').read()), open('$SWAGGER_JSON', 'w'), indent=2)" 2>/dev/null
    if [ $? -eq 0 ]; then
      echo -e "${GREEN}OpenAPI spec converted to JSON format${NC}"
    else
      echo -e "${YELLOW}Warning: Failed to convert YAML to JSON with Python. Install pyyaml module.${NC}"
    fi
  else
    echo -e "${YELLOW}Warning: Could not convert YAML to JSON. JSON file will not be available.${NC}"
  fi
else
  echo -e "${RED}Error: OpenAPI reference file not found at $API_REFERENCE_FILE${NC}"
  echo -e "${YELLOW}Create an OpenAPI specification file at $API_REFERENCE_FILE${NC}"
  exit 1
fi

# Generate ReDoc HTML documentation if redoc-cli is available
if [ "$HAVE_REDOC" = "yes" ]; then
  echo -e "${BLUE}Generating ReDoc HTML documentation...${NC}"
  redoc-cli bundle "$SWAGGER_YAML" -o "$REDOC_HTML" --options.theme.colors.primary.main="#4a35a9"
  echo -e "${GREEN}ReDoc HTML documentation generated at $REDOC_HTML${NC}"
else
  echo -e "${YELLOW}Skipping ReDoc HTML generation (redoc-cli not available)${NC}"
fi

# Verify the output files exist
echo -e "${BLUE}Validating generated documentation...${NC}"
if [ -f "$SWAGGER_JSON" ]; then
  echo -e "${GREEN}✓ Swagger JSON file generated${NC}"
  
  # Count number of endpoints if jq is available
  if [ "$HAVE_JQ" = "yes" ]; then
    ENDPOINT_COUNT=$(jq '.paths | length' "$SWAGGER_JSON")
    echo -e "${GREEN}✓ API documentation contains $ENDPOINT_COUNT endpoints${NC}"
    
    # Check if memory endpoints are documented
    if jq -e '.paths | keys[] | select(contains("/memory"))' "$SWAGGER_JSON" &>/dev/null; then
      echo -e "${GREEN}✓ Memory API endpoints are documented${NC}"
    else
      echo -e "${YELLOW}Warning: Memory API endpoints may be missing from documentation${NC}"
    fi
  fi
else
  echo -e "${RED}✗ Swagger JSON file generation failed${NC}"
fi

if [ -f "$REDOC_HTML" ]; then
  echo -e "${GREEN}✓ ReDoc HTML documentation generated${NC}"
else
  echo -e "${YELLOW}ReDoc HTML documentation was not generated${NC}"
fi

echo -e "\n${GREEN}API Documentation Generation Complete!${NC}"
echo -e "${BLUE}You can access the documentation at:${NC}"
if [ -f "$REDOC_HTML" ]; then
  echo -e "  - ReDoc UI: file://$PWD/$REDOC_HTML (open this file in your browser)"
fi
echo -e "  - Raw Swagger JSON: $PWD/$SWAGGER_JSON"
echo -e "  - OpenAPI YAML: $PWD/$SWAGGER_YAML" 