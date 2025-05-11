#!/bin/bash

set -e

# Set default project folder name
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"

# Check if /app exists and adjust ENV_PATH accordingly
if [ -d "/app" ]; then
    echo "[INFO] /app directory found, loading .env from $ENV_PATH"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    echo "[INFO] /app directory not found, loading .env from $ENV_PATH"
fi

# Load environment variables
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    echo "❌ .env file not found at $ENV_PATH"
    exit 1
fi

# Ensure required variables are defined
: "${RUN_WITH_DOCKER_COMPOSE:?❌ RUN_WITH_DOCKER_COMPOSE environment variable is not set}"
: "${PROJECT_NAME:?❌ PROJECT_NAME is not set}"

# Change into project directory
cd "${BASE_DIR}/public/${PROJECT_NAME}"

# Remote build when RUN_WITH_DOCKER_COMPOSE=yes
if [ "$RUN_WITH_DOCKER_COMPOSE" == "yes" ]; then
    docker buildx bake --load

    if docker compose version &> /dev/null; then
        docker compose up -d
    else
        docker-compose up -d
    fi
else
    make docker-build
    make docker-run
fi
