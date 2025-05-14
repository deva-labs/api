#!/bin/bash
set -e

log() {
  echo "[INFO] $1"
}

if [ -d "$BASE_DIR" ]; then
    ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"
    PROJECT_DIR="${BASE_DIR}/public/${PROJECT_NAME}"
else
    ENV_PATH="public/${PROJECT_NAME}/.env"
    PROJECT_DIR="public/${PROJECT_NAME}"
fi

# Load .env nếu tồn tại
if [ -f "$ENV_PATH" ]; then
    export $(grep -v '^#' "$ENV_PATH" | xargs)
else
    echo "⚠️  No .env file found at $ENV_PATH"
    exit 1
fi

: "${APP_NAME:?❌ APP_NAME environment variable not set}"
: "${GO_VERSION:?❌ GO_VERSION environment variable not set}"

GO_ROOT_DIR="/app/packages/go/versions/${GO_VERSION}"

if ! command -v go &>/dev/null; then
    echo "❌ 'go' command not found in PATH"
    exit 1
fi

export GOROOT="$GO_ROOT_DIR"
export PATH="$GOROOT/bin:$PATH"

CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [ "$CURRENT_GO_VERSION" != "$GO_VERSION" ]; then
    log "❌ Current Go version is $CURRENT_GO_VERSION, expected $GO_VERSION"
    exit 1
fi

cd "$PROJECT_DIR"

if [ ! -f go.mod ]; then
    log "🧱 Initializing Go module for $APP_NAME..."
    go mod init "$APP_NAME"
else
    log "📦 go.mod already exists — skipping init."
fi

go mod tidy

echo "✅ Go module initialized and tidy complete."
sleep 1
exit 0
