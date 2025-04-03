#!/bin/bash

set -e

# ANSI color codes for better readability
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

echo -e "${YELLOW}=== SentinelStacks API Documentation Sync ===${NC}"

# First, run the API documentation generator
echo -e "${BLUE}Generating API documentation...${NC}"
./scripts/generate_api_docs.sh

# Ensure the API directory exists
mkdir -p docs/api

# Copy needed files to the mkdocs structure
echo -e "${BLUE}Syncing files with MkDocs structure...${NC}"

# Copy OpenAPI file if needed
if [ -f "docs/api-reference.yaml" ]; then
  echo -e "${GREEN}✓ Found API reference file${NC}"
else
  echo -e "${RED}Error: API reference file not found${NC}"
  exit 1
fi

# Ensure API directory exists
if [ -f "docs/api/index.md" ]; then
  echo -e "${GREEN}✓ API index page exists${NC}"
else
  echo -e "${YELLOW}Warning: API index page not found. Creating a simple one...${NC}"
  mkdir -p docs/api
  cat > docs/api/index.md << 'EOF'
# API Documentation

This page provides interactive documentation for the SentinelStacks API.

## Interactive API Explorer

<div id="swagger-ui"></div>

<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js"></script>
<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />

<script>
  window.onload = function() {
    window.ui = SwaggerUIBundle({
      url: "../api-reference.yaml",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
      ],
      layout: "StandaloneLayout",
      defaultModelsExpandDepth: -1,
      displayRequestDuration: true
    });
  };
</script>

<style>
  .swagger-ui .topbar { display: none }
</style>

## API Reference

The full OpenAPI specification is available [here](../api-reference.yaml).

See the [API Usage Guide](../api-usage-guide.md) for more information.
EOF
  echo -e "${GREEN}✓ Created API index page${NC}"
fi

# Check if mkdocs.yml contains API section
if grep -q "API:" mkdocs.yml; then
  echo -e "${GREEN}✓ mkdocs.yml contains API section${NC}"
else
  echo -e "${YELLOW}Warning: mkdocs.yml doesn't contain API section. You may need to update it manually.${NC}"
  echo -e "Add the following section to your mkdocs.yml file:"
  echo -e "${BLUE}  - API:${NC}"
  echo -e "${BLUE}    - Overview: api/index.md${NC}"
  echo -e "${BLUE}    - OpenAPI Specification: api-reference.yaml${NC}"
  echo -e "${BLUE}    - Usage Guide: api-usage-guide.md${NC}"
fi

echo -e "\n${GREEN}API Documentation Sync Complete!${NC}"
echo -e "${BLUE}You can preview the documentation with:${NC}"
echo -e "  ./scripts/update_docs.sh serve"
echo -e "${BLUE}Then open:${NC} http://localhost:8000/" 