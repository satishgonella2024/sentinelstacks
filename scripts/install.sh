#!/bin/bash

# Installation directory
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.sentinel"

# Create necessary directories
mkdir -p "$CONFIG_DIR/registry"

# Copy binary
cp sentinel "$INSTALL_DIR/"

# Set permissions
chmod +x "$INSTALL_DIR/sentinel"

# Copy examples if they exist
if [ -d "examples" ]; then
    mkdir -p "$CONFIG_DIR/examples"
    cp -r examples/* "$CONFIG_DIR/examples/"
fi

echo "SentinelStacks installed successfully!"
echo "Binary location: $INSTALL_DIR/sentinel"
echo "Configuration directory: $CONFIG_DIR" 