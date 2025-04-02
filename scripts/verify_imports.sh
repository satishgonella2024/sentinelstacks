#!/bin/bash

# This script verifies that all import paths are consistent

echo "Verifying import paths..."

MODULE_NAME=$(grep -m 1 "module" go.mod | awk '{print $2}')
echo "Module name from go.mod: $MODULE_NAME"

# Find all Go files
echo "Scanning Go files for imports..."
find . -name "*.go" -not -path "./vendor/*" | while read -r file; do
    imports=$(grep -E "^import \(" "$file" -A 20 | grep -E "^\t\"" | grep "$MODULE_NAME" || true)
    if [ -n "$imports" ]; then
        echo "Found references in $file:"
        echo "$imports"
    fi
    
    other_imports=$(grep -E "import \"github.com/" "$file" | grep -v "$MODULE_NAME" || true)
    if [ -n "$other_imports" ]; then
        echo "Found other module imports in $file:"
        echo "$other_imports"
    fi
done

incorrect_imports=$(grep -r "github.com/subrahmanyagonella/the-repo/sentinelstacks" --include="*.go" . || true)
if [ -n "$incorrect_imports" ]; then
    echo -e "\nWARNING: Found incorrect import paths:"
    echo "$incorrect_imports"
    
    echo -e "\nThese imports should be changed to: $MODULE_NAME"
else
    echo "No incorrect imports found."
fi

# Check if the main.go file is correctly importing the cmd package
main_import=$(grep -E "import \(" main.go -A 10 | grep -E "^\t\"$MODULE_NAME/cmd")
if [ -n "$main_import" ]; then
    echo "main.go correctly imports: $main_import"
else
    echo "WARNING: main.go may have incorrect imports!"
fi

echo "Import verification complete."
