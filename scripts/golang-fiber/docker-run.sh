#!/bin/bash

set -e

# Set default project folder name
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"

# Check if /app directory exists and adjust ENV_PATH accordingly
if [ -d "/app" ]; then
    echo "🔍 Loading environment from $ENV_PATH"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
fi

# Load environment variables
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    exit 1
fi

# Validate required environment variables
: "${APP_NAME:?❌ APP_NAME environment variable is not set}"
: "${FRAMEWORK:?❌ FRAMEWORK environment variable is not set}"
: "${VERSION:?❌ VERSION environment variable is not set}"
: "${APP_PORT:?❌ APP_PORT environment variable is not set}"

# Run Docker container
CONTAINER_NAME="${APP_NAME}-${FRAMEWORK}"
IMAGE_NAME="${APP_NAME}-${FRAMEWORK}:${VERSION}"

docker run -d -p "${APP_PORT}:${APP_PORT}" --name "$CONTAINER_NAME" "$IMAGE_NAME"

sleep 1
exit 0