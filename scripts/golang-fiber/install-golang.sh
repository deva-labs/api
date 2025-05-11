#!/bin/bash

set -e

# Use PROJECT_NAME from environment or default
PROJECT_NAME="${PROJECT_NAME:-fiber}"
APP_BASE="/app"

if [ -d "$APP_BASE" ]; then
    WORK_DIR="public/$PROJECT_NAME"
    ENV_FILE="${APP_BASE}/${WORK_DIR}/.env"
else
    WORK_DIR="public/$PROJECT_NAME"
    ENV_FILE="${WORK_DIR}/.env"
fi

# Load .env if it exists
if [ -f "$ENV_FILE" ]; then
    echo "📂 Loading environment from $ENV_FILE"
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    echo "⚠️  No .env file found at $ENV_FILE"
    exit 1
fi

# Ensure GO_VERSION is set
if [ -z "$GO_VERSION" ]; then
    echo "❌ GO_VERSION is not set. Please define it in .env or environment."
    exit 1
fi

# Check if Go is already installed
if command -v go &> /dev/null; then
    INSTALLED_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$INSTALLED_GO_VERSION" == "$GO_VERSION" ]; then
        echo "✅ Go $GO_VERSION is already installed!"
        exit 0
    else
        echo "⚠️ Go is already installed, but version mismatch. Installing Go $GO_VERSION..."
    fi
else
    echo "❌ Go is not installed. Installing Go $GO_VERSION..."
fi

# Download and install Go
GO_TAR="go$GO_VERSION.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/$GO_TAR"
INSTALL_DIR="/usr/local"

echo "📏 Getting size of $GO_URL..."
CONTENT_LENGTH=$(curl -sI "$GO_URL" | grep -i Content-Length | awk '{print $2}' | tr -d '\r')

if [ -z "$CONTENT_LENGTH" ]; then
    echo "❌ Cannot determine content length of $GO_URL"
    exit 1
fi

echo "⬇️ Downloading Go $GO_VERSION..."
curl -s "$GO_URL" | pv -s "$CONTENT_LENGTH" > "$GO_TAR"

echo "🧹 Removing any existing Go installation at $INSTALL_DIR/go"
sudo rm -rf "$INSTALL_DIR/go"

echo "📦 Extracting $GO_TAR to $INSTALL_DIR"
sudo tar -C "$INSTALL_DIR" -xzf "$GO_TAR"

echo "🧼 Cleaning up"
rm "$GO_TAR"

# Add Go to PATH for current session
export PATH=$PATH:/usr/local/go/bin

echo "✅ Go $GO_VERSION installed successfully!"
echo "👉 Add this to your ~/.bashrc or ~/.zshrc if not already set:"
echo 'export PATH=$PATH:/usr/local/go/bin'

# Verify installation
echo "🔍 Checking if Go is properly installed..."
if ! command -v go &> /dev/null; then
    echo "❌ Go installation failed. Please ensure that Go is correctly installed."
    exit 1
else
    echo "✅ Go is correctly installed!"
fi
