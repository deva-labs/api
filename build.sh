#!/bin/bash

set -e

PROJECT_NAME="$1"

if [ -z "$PROJECT_NAME" ]; then
  echo "❌ Please provide a project name as the first argument."
  exit 1
fi

BASE_DIR="/app"
TLSCACERT_PATH="${BASE_DIR}/store/secrets/ca.pem"
TLSCERT_PATH="${BASE_DIR}/store/secrets/cert.pem"
TLSKEY_PATH="${BASE_DIR}/store/secrets/key.pem"

export PROJECT_NAME="$PROJECT_NAME"

CONTEXT_NAME="myremote"

if [ ! -f Makefile ]; then
  echo "❌ Makefile not found in the current directory."
  exit 1
fi

if ! docker context inspect "$CONTEXT_NAME" >/dev/null 2>&1; then
  docker context create myremote \
    --docker "host=tcp://192.168.237.116:2376,ca=/app/store/secrets/ca.pem,cert=/app/store/secrets/cert.pem,key=/app/store/secrets/key.pem"
fi

docker context use "$CONTEXT_NAME"

if ! docker info >/dev/null 2>&1; then
  echo "❌ Failed to connect to Docker daemon using context '$CONTEXT_NAME'."
  echo "🔧 Please check your certificate files and DOCKER_HOST."
  exit 1
fi

chmod +x ./scripts/*.sh
make build