#!/bin/bash

set -e

# Logging helper
log() {
    echo "[INFO] $1"
}

# Support dynamic project folder via PROJECT_NAME
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"

# Check if /app exists, and adjust ENV_PATH accordingly
if [ -d "/app" ]; then
    log "🔍 Loading environment variables from $ENV_PATH"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    log "🔍 /app directory not found, loading environment variables from $ENV_PATH"
fi

# Load environment variables
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
    log ".env loaded successfully"
else
    echo "❌ .env file not found at $ENV_PATH"
    exit 1
fi

# Verify required environment variables
: "${APP_NAME:?❌ APP_NAME is not set in environment}"
: "${FRAMEWORK:?❌ FRAMEWORK is not set in environment}"
: "${VERSION:?❌ VERSION is not set in environment}"

# Change to project directory
PROJECT_DIR="${BASE_DIR}/public/${PROJECT_NAME}"
if [ -d "$PROJECT_DIR" ]; then
    cd "$PROJECT_DIR"
else
    echo "❌ Project directory $PROJECT_DIR not found"
    exit 1
fi

# Perform Docker build
log "🐳 Building Docker image ${APP_NAME}-${FRAMEWORK}:${VERSION}..."
docker build -t "${APP_NAME}-${FRAMEWORK}:${VERSION}" .

log "✅ Build completed successfully."
