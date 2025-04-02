#!/bin/bash

# Comprehensive validation script for SentinelStacks project
echo "Validating SentinelStacks project..."

# 1. Check module path in go.mod
MODULE_PATH=$(grep -m 1 "module" go.mod | awk '{print $2}')
echo "Current module path: $MODULE_PATH"
if [ "$MODULE_PATH" != "github.com/satishgonella2024/sentinelstacks" ]; then
    echo "WARNING: Module path should be 'github.com/satishgonella2024/sentinelstacks'"
fi

# 2. Check for import cycle issues
echo "Checking for import cycles..."
CYCLES=$(go build -o /dev/null ./... 2>&1 | grep "import cycle not allowed" || true)
if [ -n "$CYCLES" ]; then
    echo "WARNING: Import cycles detected:"
    echo "$CYCLES"
    
    # Extract problematic packages
    for pkg in $(echo "$CYCLES" | grep -o "github.com/[^ :]*" | sort | uniq); do
        echo "  * $pkg"
    done
fi

# 3. Check for missing dependencies
echo "Checking for missing dependencies..."
MISSING=$(go build -o /dev/null ./... 2>&1 | grep "no required module provides package" || true)
if [ -n "$MISSING" ]; then
    echo "WARNING: Missing dependencies detected:"
    echo "$MISSING"
    
    # Extract missing packages
    for pkg in $(echo "$MISSING" | grep -o "package [^ ]*" | awk '{print $2}' | sort | uniq); do
        echo "  * $pkg"
    done
fi

# 4. Suggest next steps
echo
echo "Validation complete!"
echo
echo "Suggested next steps:"

if [ -n "$CYCLES" ]; then
    echo "1. Fix import cycle issues:"
    echo "   - Extract common interfaces into separate packages"
    echo "   - Break circular dependencies between packages"
    echo "   - Temporary disable problematic imports for testing"
fi

if [ -n "$MISSING" ]; then
    echo "2. Install missing dependencies:"
    echo "   - Run 'go get' for each missing package"
    echo "   - Update go.mod with correct requirements"
fi

echo "3. Run the simple test to verify DAG functionality:"
echo "   - chmod +x scripts/run_simple_test.sh"
echo "   - ./scripts/run_simple_test.sh"

echo "4. When ready, create a reduced version of the core functionality:"
echo "   - Focus only on stack engine and DAG implementation"
echo "   - Remove dependencies on complex subsystems for testing"

echo
echo "Once these issues are resolved, we can proceed with implementing the registry system."
