#!/bin/bash
# Script to check the registry directory

echo "Checking the registry directory..."
ls -la ~/.sentinel/images/

echo "Image files in the registry:"
find ~/.sentinel/images/ -type f -name "*.json" | while read file; do
  echo "File: $file"
  cat "$file" | grep -E '"name"|"tag"|"baseModel"' | sed 's/^/  /'
  echo ""
done

echo "End of registry check."
