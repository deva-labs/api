#!/bin/sh
set -e

# Load .env if exists
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Verify required variables
: "${APP_NAME:?APP_NAME environment variable not set}"

BINARY="./bin/${APP_NAME}"

# Build if binary doesn't exist
if [ ! -f "$BINARY" ]; then
  echo "Binary $BINARY not found. Building it now..."
  go build -buildvcs=false -o "$BINARY" .
fi

echo "Starting the application with Air..."
exec air