#!/bin/bash

set -e

# Function to log messages
log() {
    echo "[INFO] $1"
}

# Set default project folder if not passed via env
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"

if [ -d "$BASE_DIR" ]; then
    ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"
    AIR_TOML_DIR="${BASE_DIR}/public/${PROJECT_NAME}"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    AIR_TOML_DIR="public/${PROJECT_NAME}"
fi

AIR_TOML_PATH="${AIR_TOML_DIR}/.air.toml"

log "🔍 Loading environment from $ENV_PATH"

# Load environment variables if .env exists
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
    log ".env file loaded successfully"
else
    log "⚠️  .env file not found at $ENV_PATH"
fi

# Check required environment variables
: "${APP_NAME:?❌ APP_NAME is not set}"
: "${APP_PORT:?❌ APP_PORT is not set}"

# Create .air.toml if missing
if [ ! -f "$AIR_TOML_PATH" ]; then
    log "🛠️ Creating .air.toml..."
    mkdir -p "$AIR_TOML_DIR/bin"
    (cd "$AIR_TOML_DIR" && air init)
else
    log "✅ .air.toml already exists"
fi

log "🔧 Updating .air.toml config"

# Ensure tmp_dir is set
if grep -q "^tmp_dir =" "$AIR_TOML_PATH"; then
    sed -i 's|^tmp_dir = .*|tmp_dir = "bin"|' "$AIR_TOML_PATH"
else
    sed -i '1i tmp_dir = "bin"' "$AIR_TOML_PATH"
fi

# Clean existing [build] keys
sed -i '/^\[build\]/,/^\[.*\]/ {
  s|^ *bin = .*||g
  s|^ *cmd = .*||g
  s|^ *pre_cmd = .*||g
  s|^ *poll = .*||g
  s|^ *poll_interval = .*||g
  s|^ *full_bin = .*||g
}' "$AIR_TOML_PATH"

# Inject new build config
sed -i "/^\[build\]/a\
  cmd = \"go build -buildvcs=false -o ./bin/${APP_NAME} .\"\n\
  poll = true\n\
  poll_interval = 500\n\
  full_bin = \"APP_PORT=${APP_PORT} ./bin/${APP_NAME}\"\n\
  pre_cmd = [\"go mod tidy\"]\n\
  bin = \"./bin/${APP_NAME}-api\"" "$AIR_TOML_PATH"

# Set log level if missing
if ! grep -q 'level = "debug"' "$AIR_TOML_PATH"; then
    sed -i '/^\[log\]/ a\  level = "debug"' "$AIR_TOML_PATH"
fi

log "✅ .air.toml configured successfully in $AIR_TOML_DIR"
exit 0
