#!/bin/bash

set -e

# --- Configuration ---
PROJECT_NAME="${PROJECT_NAME:-fiber}"
WORKDIR="${BASE_DIR}/public/${PROJECT_NAME}"
ENV_FILE="$WORKDIR/.env"
RUN_WITH_DOCKER_COMPOSE="${RUN_WITH_DOCKER_COMPOSE:-false}"

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    ALT_WORKDIR="./public/${PROJECT_NAME}"
    ALT_ENV_FILE="${ALT_WORKDIR}/.env"
    if [ -f "$ALT_ENV_FILE" ]; then
        WORKDIR="$ALT_WORKDIR"
        ZIP_PATH="./public/${PROJECT_NAME}.zip"
        export $(grep -v '^#' "$ALT_ENV_FILE" | xargs)
    else
        echo "❌ .env file not found!" >&2
        exit 1
    fi
fi

: "${APP_NAME:?APP_NAME environment variable is not set}"
: "${FRAMEWORK:?FRAMEWORK environment variable is not set}"
: "${APP_VERSION:?APP_VERSION environment variable is not set}"
ZIP_PATH="${ZIP_PATH:-${WORKDIR}.zip}"

CONTAINER_NAME="${APP_NAME}-${FRAMEWORK}"
IMAGE_NAME="${CONTAINER_NAME}:${APP_VERSION}"

# --- Step 1: Cleanup containers/images ---
if [ "$RUN_WITH_DOCKER_COMPOSE" = "true" ]; then
  if [ -d "$WORKDIR" ]; then
    echo "🔴 Stopping and removing Docker Compose containers..." >&2
    cd "$WORKDIR"
    docker compose down --rmi all --volumes --remove-orphans > /dev/null 2>&1 || true
    cd - > /dev/null 2>&1 || exit
  fi
else
  echo "🔴 Stopping and removing standalone Docker container..." >&2
  docker stop "$CONTAINER_NAME" > /dev/null 2>&1 || true
  docker rm "$CONTAINER_NAME" > /dev/null 2>&1 || true

  echo "🗑️ Removing Docker image: $IMAGE_NAME" >&2
  docker rmi "$IMAGE_NAME" > /dev/null 2>&1 || true
fi

# --- Step 2: Zip working directory ---
if [ -d "$WORKDIR" ]; then
  echo "📦 Zipping project folder..." >&2
  mkdir -p "$(dirname "$ZIP_PATH")"
  cd "$WORKDIR"
  zip -rq "$ZIP_PATH" . -x "*/node_modules/*"
  cd - > /dev/null 2>&1
else
  echo "⚠️ Directory not found. Skipping zipping." >&2
fi

# --- Step 3: Delete working directory ---
if [ -d "$WORKDIR" ]; then
  echo "🧹 Removing project folder..." >&2
  rm -rf "$WORKDIR"
fi

sleep 1
exit 0
