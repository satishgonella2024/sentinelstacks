#!/bin/bash

# Create SSL directory if it doesn't exist
mkdir -p nginx/ssl

# Generate SSL certificate and key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout nginx/ssl/localhost.key \
    -out nginx/ssl/localhost.crt \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Set proper permissions
chmod 600 nginx/ssl/localhost.key
chmod 644 nginx/ssl/localhost.crt

echo "SSL certificates generated successfully!" 