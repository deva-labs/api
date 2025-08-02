#!/bin/bash
set -e

# Determine env file location
if [ -d "$BASE_DIR" ]; then
    ENV_FILE="${BASE_DIR}/public/${PROJECT_NAME}/.env"
else
    ENV_FILE="public/${PROJECT_NAME}/.env"
fi

# Load .env if it exists
if [ -f "$ENV_FILE" ]; then
    echo "[INFO] Loading environment from $ENV_FILE"
    set -o allexport
    source "$ENV_FILE"
    set +o allexport
else
    echo "[ERROR] .env file not found: $ENV_FILE"
    exit 1
fi

# Ensure GO_VERSION is set
if [ -z "$GO_VERSION" ]; then
    echo "[ERROR] GO_VERSION not set in .env"
    exit 1
fi

# Target installation path
GO_ROOT_DIR="/app/packages/go/versions/${GO_VERSION}"
GO_BIN="${GO_ROOT_DIR}/bin/go"

if [ -x "$GO_BIN" ]; then
    echo "[INFO] Go $GO_VERSION already installed at $GO_ROOT_DIR"
else
    # Download Go tarball
    GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
    GO_URL="https://go.dev/dl/${GO_TAR}"
    INSTALL_TAR="public/${GO_TAR}"

    echo "[INFO] Downloading $GO_TAR from $GO_URL..."
    curl -sSL -o "$INSTALL_TAR" "$GO_URL"

    mkdir -p "$GO_ROOT_DIR"
    tar -C "$GO_ROOT_DIR" --strip-components=1 -xzf "$INSTALL_TAR"
fi

sleep 1
exit 0
