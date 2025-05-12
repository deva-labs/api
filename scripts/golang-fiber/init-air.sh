#!/bin/bash

set -e

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

# Load environment variables if .env exists
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    echo "⚠️  .env file not found at $ENV_PATH"
fi

# Check required environment variables
: "${APP_NAME:?❌ APP_NAME is not set}"
: "${APP_PORT:?❌ APP_PORT is not set}"

# Create .air.toml if missing
if [ ! -f "$AIR_TOML_PATH" ]; then
    mkdir -p "$AIR_TOML_DIR/bin"
    (cd "$AIR_TOML_DIR" && air init)
else
    echo "✅ .air.toml already exists"
fi

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

echo "✅ .air.toml configured successfully"

sleep 1
exit 0