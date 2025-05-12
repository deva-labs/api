#!/bin/bash

set -e

# Check if /app exists
if [ -d "/app" ]; then
  # Clean up old .env if inside /app
  rm -f "/app/public/${PROJECT_NAME}/.env"
  echo "[INFO] Old .env file removed from /app/public/${PROJECT_NAME}/.env"
else
  # Clean up old .env if outside /app
  rm -f "public/${PROJECT_NAME}/.env"
  echo "[INFO] Old .env file removed from public/${PROJECT_NAME}/.env"
fi

# Ensure directory exists
mkdir -p "public/${PROJECT_NAME}"

# Create new .env file
cat <<EOF > "public/${PROJECT_NAME}/.env"
GO_VERSION=1.24.2
APP_NAME=${PROJECT_NAME}
VERSION=1.0.0
FRAMEWORK=fiber
GO_TAR=go\${GO_VERSION}.linux-amd64.tar.gz
WITH_DB=yes
RUN_WITH_DOCKER_COMPOSE=yes
ENV=dev
IMAGE_NAME=\${APP_NAME}:latest
DB_PASS=${PROJECT_NAME}-admin
DB_NAME=${PROJECT_NAME}
DB_USER=admin
DB_PORT=3307
APP_PORT=2350
EOF

sleep 1
exit 0