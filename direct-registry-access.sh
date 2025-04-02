#!/bin/bash
# Script to directly read and display images from the registry

REGISTRY_DIR="$HOME/.sentinel/images"

echo "Reading images directly from registry: $REGISTRY_DIR"
echo 

if [ ! -d "$REGISTRY_DIR" ]; then
  echo "Registry directory not found!"
  exit 1
fi

# Count the number of JSON files
JSON_COUNT=$(find "$REGISTRY_DIR" -name "*.json" | wc -l)
echo "Found $JSON_COUNT image files"
echo

printf "%-12s %-30s %-10s %-20s %-10s %-20s\n" "IMAGE ID" "NAME" "TAG" "CREATED" "SIZE" "BASE MODEL"
echo "------------------------------------------------------------------------------------------------------"

# Process each JSON file
find "$REGISTRY_DIR" -name "*.json" | while read -r file; do
  # Generate a fake ID
  ID=$(echo "$file" | md5sum | cut -c1-12)
  
  # Extract name and tag from filename
  FILENAME=$(basename "$file")
  NAME=$(echo "$FILENAME" | cut -d'_' -f1 | tr '_' '/')
  TAG=$(echo "$FILENAME" | cut -d'_' -f2 | sed 's/\.json$//')
  
  # Get file size
  SIZE=$(du -h "$file" | cut -f1)
  
  # Get creation time
  CREATED=$(stat -c %y "$file" | cut -d' ' -f1)
  
  # Try to extract baseModel from the JSON
  BASE_MODEL=$(grep -o '"baseModel"[^,}]*' "$file" | cut -d'"' -f4)
  
  printf "%-12s %-30s %-10s %-20s %-10s %-20s\n" "$ID" "$NAME" "$TAG" "$CREATED" "$SIZE" "$BASE_MODEL"
done

echo
echo "End of registry listing"
