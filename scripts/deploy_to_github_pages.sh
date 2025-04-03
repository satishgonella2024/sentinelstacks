#!/bin/bash

set -e

# ANSI color codes for better readability
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

echo -e "${YELLOW}=== SentinelStacks GitHub Pages Deployment ===${NC}"

# Check if the user specified a custom commit message
if [ $# -eq 0 ]; then
  COMMIT_MESSAGE="Update documentation $(date +'%Y-%m-%d %H:%M:%S')"
else
  COMMIT_MESSAGE="$1"
fi

# First, sync the API documentation
echo -e "${BLUE}Syncing API documentation...${NC}"
./scripts/sync_api_docs.sh

# Check if we're on the main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
  echo -e "${YELLOW}Warning: You are not on the main branch.${NC}"
  echo -e "Current branch: ${BLUE}$CURRENT_BRANCH${NC}"
  read -p "Do you want to continue? (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Deployment aborted.${NC}"
    exit 1
  fi
fi

# Build and deploy the documentation
echo -e "${BLUE}Building and deploying documentation...${NC}"
mkdocs gh-deploy --force --message "$COMMIT_MESSAGE"

echo -e "\n${GREEN}Documentation successfully deployed to GitHub Pages!${NC}"
echo -e "${BLUE}Visit your GitHub Pages site to view the documentation.${NC}"
echo -e "${BLUE}The URL is typically: https://YOUR-USERNAME.github.io/YOUR-REPO/${NC}" 