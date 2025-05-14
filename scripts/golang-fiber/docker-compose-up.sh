#!/bin/bash

set -e

BASE_DIR="${BASE_DIR:-/app}"
CONTEXT_NAME="${CONTEXT_NAME:-myremote}"

ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"

# Check BASE_DIR existence
if [ -d "$BASE_DIR" ]; then
  echo "[INFO] $BASE_DIR not found, loading .env from $ENV_PATH"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    echo "[INFO] $BASE_DIR not found, loading .env from $ENV_PATH"
fi

# Load .env variables
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    exit 1
fi

# Ensure essential variables exist
: "${RUN_WITH_DOCKER_COMPOSE:?❌ RUN_WITH_DOCKER_COMPOSE environment variable is not set}"
: "${PROJECT_NAME:?❌ PROJECT_NAME is not set}"
: "${TLSCACERT_PATH:?❌ TLSCACERT_PATH is not set}"
: "${TLSCERT_PATH:?❌ TLSCERT_PATH is not set}"
: "${TLSKEY_PATH:?❌ TLSKEY_PATH is not set}"
: "${CONTEXT_NAME:?❌ CONTEXT_NAME is not set}"

# Change to project directory
cd "${BASE_DIR}/public/${PROJECT_NAME}" >/dev/null 2>&1

# Docker context setup (quiet mode)
if ! docker context inspect "$CONTEXT_NAME" >/dev/null 2>&1; then
  docker context create "$CONTEXT_NAME" \
    --docker "host=tcp://192.168.237.116:2376,ca=${TLSCACERT_PATH},cert=${TLSCERT_PATH},key=${TLSKEY_PATH}" >/dev/null 2>&1
fi

docker context use "$CONTEXT_NAME" >/dev/null 2>&1

# Execute commands based on RUN_WITH_DOCKER_COMPOSE
if [ "$RUN_WITH_DOCKER_COMPOSE" == "true" ]; then
    docker buildx bake --load >/dev/null 2>&1
    sleep 3
    if docker compose version &> /dev/null; then
        docker compose up -d >/dev/null 2>&1
        sleep 1
    else
        docker-compose up -d >/dev/null 2>&1
        sleep 1
    fi
else
    echo "[INFO] Running local Docker build and run using Makefile..."
    cd "${BASE_DIR}"
    make docker-build PROJECT_NAME="$PROJECT_NAME" >/dev/null 2>&1
    make docker-run PROJECT_NAME="$PROJECT_NAME" >/dev/null 2>&1
fi

sleep 1
exit 0