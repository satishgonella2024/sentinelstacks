#!/bin/bash
# Script to override the images command in the sentinel executable

set -e # Exit on error

echo "This script will create a modified 'sentinel' executable that displays real images from the registry."
echo "It will create a new file called 'sentinel-real' to avoid modifying your original executable."
echo

# Locate the hardcoded data in the sentinel binary
echo "Searching for hardcoded image data in the sentinel executable..."
if strings ./sentinel | grep -q "sha256:abcde"; then
  echo "Found hardcoded image data in the sentinel executable."
else
  echo "Could not find hardcoded image data. This approach may not work."
fi

echo
echo "Building standalone images command..."
go build -o real-images real-images.go

echo
echo "Creating sentinel-real script..."
cat > sentinel-real << 'EOF'
#!/bin/bash

if [ "$1" = "images" ]; then
  # Run the real-images command
  ./real-images
else
  # Pass all arguments to the original sentinel
  ./sentinel "$@"
fi
EOF

chmod +x sentinel-real

echo
echo "Done! You can now use './sentinel-real images' to see the real images in your registry."
