#!/bin/bash
# SentinelStacks Deployment Script
# This script should be placed on your Proxmox Ubuntu VM

set -e  # Exit immediately if a command exits with a non-zero status

# Configuration
APP_NAME="sentinelstacks"
APP_DIR="/home/$(whoami)/sentinelstacks"
WEB_DIR="/var/www/sentinelstacks"
LOG_FILE="$APP_DIR/deploy.log"

# Log function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Create necessary directories
mkdir -p "$APP_DIR/backup" "$APP_DIR/app" "$APP_DIR/web-ui"

# Backup current version
backup() {
    log "Creating backup..."
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    
    # Backup executable if it exists
    if [ -f "/usr/local/bin/sentinel" ]; then
        cp /usr/local/bin/sentinel "$APP_DIR/backup/sentinel_$TIMESTAMP"
    fi
    
    # Backup web files if they exist
    if [ -d "$WEB_DIR" ] && [ "$(ls -A $WEB_DIR)" ]; then
        mkdir -p "$APP_DIR/backup/web_$TIMESTAMP"
        cp -r "$WEB_DIR/"* "$APP_DIR/backup/web_$TIMESTAMP/"
    fi
    
    log "Backup created at $APP_DIR/backup with timestamp $TIMESTAMP"
}

# Deploy new version
deploy() {
    log "Starting deployment process..."
    
    # Check if new binary exists
    if [ -f "$APP_DIR/app/sentinel" ]; then
        log "Installing new sentinel binary..."
        chmod +x "$APP_DIR/app/sentinel"
        sudo mv "$APP_DIR/app/sentinel" /usr/local/bin/
        log "Binary installed to /usr/local/bin/sentinel"
    else
        log "ERROR: No sentinel binary found in $APP_DIR/app/"
        return 1
    fi
    
    # Check if web files exist
    if [ -d "$APP_DIR/web-ui" ] && [ "$(ls -A $APP_DIR/web-ui)" ]; then
        log "Deploying web UI files..."
        sudo mkdir -p "$WEB_DIR"
        sudo cp -r "$APP_DIR/web-ui/"* "$WEB_DIR/"
        sudo chown -R www-data:www-data "$WEB_DIR"
        log "Web UI deployed to $WEB_DIR"
    else
        log "WARNING: No web UI files found in $APP_DIR/web-ui/"
    fi
    
    # Restart service
    log "Restarting service..."
    sudo systemctl restart sentinelstacks
    
    # Check service status
    if sudo systemctl is-active --quiet sentinelstacks; then
        log "Service restarted successfully"
    else
        log "ERROR: Service failed to start"
        log "Service status:"
        sudo systemctl status sentinelstacks --no-pager
        return 1
    fi
    
    log "Deployment completed successfully"
}

# Main function
main() {
    log "==== Deployment started ===="
    
    # Create backup
    backup
    
    # Deploy new version
    deploy
    
    log "==== Deployment finished ===="
}

# Run main function
main