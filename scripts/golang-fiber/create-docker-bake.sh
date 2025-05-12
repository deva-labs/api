#!/bin/bash
set -e

# --- Configuration ---
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
WORKDIR="public/${PROJECT_NAME}"
ENV_FILE="${BASE_DIR}/${WORKDIR}/.env"
OUTPUT_FILE="${WORKDIR}/docker-bake.hcl"

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    ALT_ENV_FILE="./${WORKDIR}/.env"
    if [ -f "$ALT_ENV_FILE" ]; then
        ENV_FILE="$ALT_ENV_FILE"
        export $(grep -v '^#' "$ENV_FILE" | xargs)
    else
        exit 1
    fi
fi

# Check required env variables
: "${APP_NAME:?APP_NAME environment variable not set}"
: "${FRAMEWORK:?FRAMEWORK environment variable not set}"
: "${VERSION:?VERSION environment variable not set}"
: "${DB_PASS:?DB_PASS environment variable not set}"
: "${DB_NAME:?DB_NAME environment variable not set}"
: "${DB_USER:?DB_USER environment variable not set}"

# Ensure output directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Write to docker-bake.hcl
cat <<EOF > "$OUTPUT_FILE"
# Define variables with default values
APP_NAME = "${APP_NAME}"
FRAMEWORK = "${FRAMEWORK}"
VERSION = "${VERSION}"
DB_PASS = "${DB_PASS}"
DB_NAME = "${DB_NAME}"
DB_USER = "${DB_USER}"

group "default" {
  targets = ["app", "db"]
}

target "app" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["\${APP_NAME}-\${FRAMEWORK}:\${VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    APP_NAME = APP_NAME
    VERSION = VERSION
  }
}

target "db" {
  context = "."
  dockerfile = "Dockerfile.mysql"
  tags = ["\${APP_NAME}-mysql:\${VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    MYSQL_ROOT_PASSWORD = DB_PASS
    MYSQL_DATABASE = DB_NAME
    MYSQL_USER = DB_USER
    MYSQL_PASSWORD = DB_PASS
  }
}
EOF

sleep 1
exit 0