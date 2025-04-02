#!/bin/bash

# This script fixes all incorrect import paths in the codebase

OLD_IMPORT="github.com/sentinelstacks/sentinel"
NEW_IMPORT="github.com/subrahmanyagonella/the-repo/sentinelstacks"

echo "Fixing import paths from $OLD_IMPORT to $NEW_IMPORT..."

# Find all Go files with the old import path and fix them
find . -name "*.go" -not -path "./vendor/*" -exec grep -l "$OLD_IMPORT" {} \; | while read -r file; do
    echo "Fixing imports in $file"
    sed -i '' "s|$OLD_IMPORT|$NEW_IMPORT|g" "$file" 
done

echo "Import path fixing complete."
echo "Now run 'go mod tidy' to update dependencies."
