#!/bin/sh
set -e

BINARY="./bin/dockerwizard-api"

if [ ! -f "$BINARY" ]; then
  echo "Binary $BINARY not found. Building it now..."
  go build -buildvcs=false -o "$BINARY" .
fi

echo "Starting the application with Air..."
exec air
