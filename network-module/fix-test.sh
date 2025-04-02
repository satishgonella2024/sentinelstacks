#!/bin/bash
set -e

# Make scripts executable
chmod +x test-network.sh
chmod +x run-test.sh
chmod +x make-test-executable.sh

# Create the most simplified test
echo '#!/bin/bash
set -e

# Clean any previous build
rm -f sentinel-network

# Build the binary
echo "Building..."
go mod tidy
go build -o sentinel-network .
echo "Built successfully!"

# Run a simple test
echo -e "\nCreating network..."
./sentinel-network network create simple-network

echo -e "\nListing networks..."
./sentinel-network network ls

echo -e "\nTest completed!"
' > simple-test.sh

chmod +x simple-test.sh

echo "Fixed! Run ./simple-test.sh for the simplest test"
