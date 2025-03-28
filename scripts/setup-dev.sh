#!/bin/bash

# Exit on error
set -e

# Create necessary directories
mkdir -p nginx/conf.d nginx/ssl auth

# Generate SSL certificate for development
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout nginx/ssl/key.pem \
    -out nginx/ssl/cert.pem \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Create initial user for registry
docker run --entrypoint htpasswd \
    httpd:2 -Bbn admin admin > auth/htpasswd

# Build and start services
docker-compose up -d

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 10

# Check health
curl -k https://localhost/health || echo "Health check failed"

echo "Development environment is ready!"
echo "Landing page: https://localhost"
echo "Registry: https://localhost/v2/"
echo "API: https://localhost/api/"
echo ""
echo "Default credentials:"
echo "Username: admin"
echo "Password: admin" 