#!/bin/bash

# Script to commit and push the CLI enhancement changes

# Set colors for output
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}Creating a new branch for the CLI enhancements...${NC}"
git checkout -b feature/improved-cli-spinners

echo -e "${CYAN}Adding new files...${NC}"
git add pkg/ui/spinner.go
git add cmd/sentinel/main.go
git add NEXT_STEPS.md

echo -e "${CYAN}Committing changes...${NC}"
git commit -m "Add animated spinners and improve CLI user experience"

echo -e "${GREEN}Changes committed successfully!${NC}"
echo 
echo -e "${CYAN}You can now push your changes with:${NC}"
echo "git push -u origin feature/improved-cli-spinners"
echo
echo -e "${CYAN}And create a pull request to merge into main.${NC}"
