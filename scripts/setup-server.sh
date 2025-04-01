#!/bin/bash
# SentinelStacks Server Setup Script
# Run this script once on your Proxmox Ubuntu VM to set up the environment

set -e  # Exit immediately if a command exits with a non-zero status

echo "=== Starting SentinelStacks Server Setup ==="

# Update and install dependencies
echo "Updating system packages..."
sudo apt update
sudo apt upgrade -y

echo "Installing required dependencies..."
sudo apt install -y build-essential curl git wget gnupg2 lsb-release software-properties-common nginx

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    
    # Set up Go environment
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    source ~/.bashrc
else
    echo "Go is already installed: $(go version)"
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt install -y nodejs
else
    echo "Node.js is already installed: $(node --version)"
fi

# Create required directories
echo "Setting up application directories..."
mkdir -p ~/sentinelstacks/app ~/sentinelstacks/web-ui ~/sentinelstacks/backup

# Set up Nginx
echo "Setting up Nginx..."
sudo cp docs/deployment/sentinelstacks.nginx /etc/nginx/sites-available/sentinelstacks
sudo ln -sf /etc/nginx/sites-available/sentinelstacks /etc/nginx/sites-enabled/
sudo mkdir -p /var/www/sentinelstacks

# Create web server root directory
sudo mkdir -p /var/www/sentinelstacks
sudo chown -R serveradmin:serveradmin /var/www/sentinelstacks

# Set up systemd service
echo "Setting up systemd service..."
sudo cp docs/deployment/sentinelstacks.service /etc/systemd/system/sentinelstacks.service
sudo systemctl daemon-reload
sudo systemctl enable sentinelstacks

# Test Nginx configuration
echo "Testing Nginx configuration..."
sudo nginx -t
sudo systemctl reload nginx

echo "==== Setup completed successfully ===="
echo "You can now deploy SentinelStacks using GitHub Actions"
echo "Nginx status: $(systemctl is-active nginx)"
echo "Note: The sentinelstacks service will start after the first deployment"
