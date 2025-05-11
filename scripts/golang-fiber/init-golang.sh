#!/bin/bash

set -e

# Use PROJECT_NAME from environment or fallback to default
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
APP_BASE="/app"

if [ -d "$APP_BASE" ]; then
    ENV_PATH="${APP_BASE}/public/${PROJECT_NAME}/.env"
    PROJECT_DIR="${APP_BASE}/public/${PROJECT_NAME}"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    PROJECT_DIR="public/${PROJECT_NAME}"
fi

# Load .env if it exists
if [ -f "$ENV_PATH" ]; then
    echo "📂 Loading environment from $ENV_PATH"
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    echo "⚠️  No .env file found at $ENV_PATH"
fi

# Verify APP_NAME is set
: "${APP_NAME:?❌ APP_NAME environment variable not set}"

# Change into the project directory
cd "$PROJECT_DIR"

# Initialize Go module only if it doesn't exist
if [ ! -f go.mod ]; then
    echo "🧱 Initializing Go module for $APP_NAME..."
    go mod init "$APP_NAME"
else
    echo "📦 go.mod already exists — skipping init."
fi

# Tidy up dependencies (creates go.sum if needed)
echo "📚 Running go mod tidy..."
go mod tidy

echo "✅ Go module initialized and tidy complete."
exit 0
