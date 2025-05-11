#!/bin/bash

set -euo pipefail

# --- Configuration ---
PROJECT_NAME="${PROJECT_NAME:-fiber}"
BASE_PATH="/app"
WORKDIR="${BASE_PATH}/public/${PROJECT_NAME}"
ENV_FILE="$WORKDIR/.env"

log() {
  echo "[INFO] $1"
}

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
    log "Loaded environment variables from $ENV_FILE"
else
    # Try without BASE_PATH
    ALT_WORKDIR="./public/${PROJECT_NAME}"
    ALT_ENV_FILE="${ALT_WORKDIR}/.env"
    if [ -f "$ALT_ENV_FILE" ]; then
        WORKDIR="$ALT_WORKDIR"
        ZIP_PATH="./public/${PROJECT_NAME}.zip"
        export $(grep -v '^#' "$ALT_ENV_FILE" | xargs)
        log "Loaded environment variables from fallback path $ALT_ENV_FILE"
    else
        echo "❌ .env file not found at $ENV_FILE or fallback $ALT_ENV_FILE"
        exit 1
    fi
fi

# Validate required environment variables
: "${APP_NAME:?APP_NAME environment variable is not set}"

# Ensure ZIP_PATH is set after fallback
ZIP_PATH="${ZIP_PATH:-${WORKDIR}.zip}"

# Step 1: Stop and remove Docker containers
if [ -d "$WORKDIR" ]; then
  echo -e "🔴 Stopping and removing Docker containers..."
  cd "$WORKDIR"
  docker compose down --rmi all --volumes --remove-orphans || true
  cd -
else
  echo -e "⚠️ Working directory $WORKDIR not found. Skipping Docker cleanup."
fi

# Step 2: Zip working directory
if [ -d "$WORKDIR" ]; then
  echo -e "📦 Zipping $WORKDIR to $ZIP_PATH..."
  mkdir -p "$(dirname "$ZIP_PATH")"
  zip -rq "$ZIP_PATH" "$WORKDIR"
else
  echo -e "⚠️ Directory $WORKDIR not found. Skipping zipping."
fi

# Step 3: Delete working directory
if [ -d "$WORKDIR" ]; then
  echo -e "🧹 Deleting working directory $WORKDIR..."
  rm -rf "$WORKDIR"
fi
