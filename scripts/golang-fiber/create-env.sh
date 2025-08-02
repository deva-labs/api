#!/bin/bash

set -e

# Print warning if some critical envs are missing
: "${PROJECT_NAME:?PROJECT_NAME is required}"
: "${GO_VERSION:?GO_VERSION is required}"
: "${AIR_VERSION:?AIR_VERSION is required}"
: "${DB_TYPE:?DB_TYPE is required}"
: "${APP_VERSION:?APP_VERSION is required}"
: "${DB_VERSION:?DB_VERSION is required}"
: "${FRAMEWORK:?FRAMEWORK is required}"
: "${WITH_DB:?WITH_DB is required}"
: "${RUN_WITH_DOCKER_COMPOSE:?RUN_WITH_DOCKER_COMPOSE is required}"
: "${ENV:?ENV is required}"
: "${DB_PASS:?DB_PASS is required}"
: "${DB_NAME:?DB_NAME is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PORT:?DB_PORT is required}"
: "${APP_PORT:?APP_PORT is required}"

APP_NAME="${PROJECT_NAME}"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
IMAGE_NAME=${APP_NAME}:${APP_VERSION}
# Check if /app exists
if [ -d "/app" ]; then
  rm -f "/app/public/${PROJECT_NAME}/.env"
  echo "[INFO] Old .env file removed from /app/public/${PROJECT_NAME}/.env"
else
  rm -f "public/${PROJECT_NAME}/.env"
  echo "[INFO] Old .env file removed from public/${PROJECT_NAME}/.env"
fi

# Ensure directory exists
mkdir -p "public/${PROJECT_NAME}"

# Create .env file
cat <<EOF > "public/${PROJECT_NAME}/.env"
GO_VERSION=${GO_VERSION}
AIR_VERSION=${AIR_VERSION}
DB_TYPE=${DB_TYPE}
APP_NAME=${APP_NAME}
APP_VERSION=${APP_VERSION}
DB_VERSION=${DB_VERSION}
FRAMEWORK=${FRAMEWORK}
WITH_DB=${WITH_DB}
RUN_WITH_DOCKER_COMPOSE=${RUN_WITH_DOCKER_COMPOSE}
ENV=${ENV}
DB_PASS=${DB_PASS}
DB_NAME=${DB_NAME}
DB_USER=${DB_USER}
DB_PORT=${DB_PORT}
APP_PORT=${APP_PORT}
EOF

echo "[INFO] .env file created at public/${PROJECT_NAME}/.env"
sleep 1
exit 0
