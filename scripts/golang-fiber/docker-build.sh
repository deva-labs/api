#!/bin/bash

set -e


# Support dynamic project folder via PROJECT_NAME
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"

# Check if /app exists, and adjust ENV_PATH accordingly
if [ -d "/app" ]; then
    echo "🔍 Loading environment variables from $ENV_PATH"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
fi

# Load environment variables
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
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
    exit 1
fi

docker build -t "${APP_NAME}-${FRAMEWORK}:${VERSION}" .

sleep 1
exit 0