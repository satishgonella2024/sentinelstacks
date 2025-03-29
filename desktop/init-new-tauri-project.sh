#!/bin/bash

# Create a temporary directory
tmp_dir=$(mktemp -d)
cd "$tmp_dir" || exit 1

# Initialize a new Tauri project to get the correct configuration format
cargo init --bin temp-tauri-app
cd temp-tauri-app || exit 1
cargo add tauri@2.0.0
cargo add tauri-build@2.0.0 --build
mkdir -p src-tauri/src

# Create a minimal main.rs
cat > src-tauri/src/main.rs << EOF
fn main() {
    println!("Hello, world!");
}
EOF

# Initialize Tauri
echo "Initializing a fresh Tauri project to get the correct configuration..."
cargo tauri init

# Copy the generated tauri.conf.json to the real project
echo "Copying the generated tauri.conf.json to your project..."
cp src-tauri/tauri.conf.json "/Users/subrahmanyagonella/SentinelStacks/desktop/src-tauri/"

# Clean up
cd ../../
rm -rf "$tmp_dir"

echo "Done! A fresh tauri.conf.json has been copied to your project."
echo "You may need to modify it to match your project's specific needs."
