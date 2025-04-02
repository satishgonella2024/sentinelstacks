#!/bin/bash

# This script fixes all incorrect import paths in the codebase

# First old import path
OLD_IMPORT1="github.com/sentinelstacks/sentinel"
# Second old import path
OLD_IMPORT2="github.com/subrahmanyagonella/the-repo/sentinelstacks"
# New correct import path
NEW_IMPORT="github.com/satishgonella2024/sentinelstacks"

echo "Fixing import paths..."

# Find all Go files with the first old import path and fix them
find . -name "*.go" -not -path "./vendor/*" -exec grep -l "$OLD_IMPORT1" {} \; | while read -r file; do
    echo "Fixing imports in $file (sentinelstacks/sentinel -> satishgonella2024/sentinelstacks)"
    sed -i '' "s|$OLD_IMPORT1|$NEW_IMPORT|g" "$file" 
done

# Find all Go files with the second old import path and fix them
find . -name "*.go" -not -path "./vendor/*" -exec grep -l "$OLD_IMPORT2" {} \; | while read -r file; do
    echo "Fixing imports in $file (subrahmanyagonella/the-repo -> satishgonella2024)"
    sed -i '' "s|$OLD_IMPORT2|$NEW_IMPORT|g" "$file" 
done

echo "Import path fixing complete."
echo "Now run 'go mod tidy' to update dependencies."
