# SentinelStacks Proxmox Deployment Guide

This guide explains how to deploy SentinelStacks on an Ubuntu VM running on Proxmox.

## Prerequisites

- Ubuntu 22.04 LTS or newer VM running on Proxmox
- SSH access to the VM with sudo privileges
- Git installed on the VM
- Internet access for downloading dependencies
- API keys for LLM providers (Claude, OpenAI, or local Ollama setup)

## Server Setup

### System Updates

```bash
# Update package lists and upgrade system
sudo apt update
sudo apt upgrade -y

# Install required system dependencies
sudo apt install -y build-essential curl git wget gnupg2 lsb-release software-properties-common
```

### Install Go

```bash
# Install Go 1.21 or newer
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
rm go1.21.0.linux-amd64.tar.gz

# Set up environment
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### Install Node.js and npm

```bash
# Install Node.js and npm
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Verify installation
node --version
npm --version
```

### Install Nginx

```bash
# Install Nginx
sudo apt install -y nginx

# Enable and start Nginx
sudo systemctl enable nginx
sudo systemctl start nginx

# Allow Nginx through firewall
sudo ufw allow 'Nginx Full'
```

## Create Systemd Service

Create a systemd service to manage the SentinelStacks backend:

```bash
# Create systemd service file
sudo nano /etc/systemd/system/sentinelstacks.service
```

Add the following content:

```ini
[Unit]
Description=SentinelStacks AI Agent Management System
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME
ExecStart=/usr/local/bin/sentinel api serve
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sentinelstacks
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

Replace `YOUR_USERNAME` with your actual username.

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable sentinelstacks
sudo systemctl start sentinelstacks
```

## Configure Nginx

```bash
# Create Nginx configuration file
sudo nano /etc/nginx/sites-available/sentinelstacks
```

Add the following content:

```nginx
server {
    listen 80;
    server_name your-domain-or-ip;

    location / {
        root /var/www/sentinelstacks;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable the site and reload Nginx:

```bash
sudo ln -s /etc/nginx/sites-available/sentinelstacks /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## GitHub Secrets Setup

To enable GitHub Actions deployment, you need to add the following secrets to your GitHub repository:

1. **SSH_PRIVATE_KEY**: Your private SSH key for connecting to the server
2. **HOST_IP**: The IP address or hostname of your Proxmox VM
3. **SSH_USER**: The SSH username to use for deployment

## Manual Deployment

If you prefer to deploy manually:

```bash
# Clone the repository
git clone https://github.com/satishgonella2024/sentinelstacks.git
cd sentinelstacks

# Build the application
go build -o sentinel ./cmd/sentinel

# Install the binary
sudo mv sentinel /usr/local/bin/

# Build the web UI
cd web-ui
npm install
npm run build

# Deploy the web UI
sudo mkdir -p /var/www/sentinelstacks
sudo cp -r dist/* /var/www/sentinelstacks/

# Restart the service
sudo systemctl restart sentinelstacks
```

## Troubleshooting

### Check Service Status

```bash
sudo systemctl status sentinelstacks
```

### View Logs

```bash
sudo journalctl -u sentinelstacks -f
```

### Check Nginx Configuration

```bash
sudo nginx -t
```

### Check Web UI Files

```bash
ls -la /var/www/sentinelstacks
```