#!/bin/bash
set -e

# --- Configuration ---
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
        echo "❌ .env file not found"
        exit 1
    fi
fi

# Check required env variables
: "${APP_NAME:?APP_NAME environment variable not set}"
: "${FRAMEWORK:?FRAMEWORK environment variable not set}"
: "${APP_VERSION:?APP_VERSION environment variable not set}"
: "${WITH_DB:?WITH_DB environment variable not set}"
: "${DB_PASS:?DB_PASS environment variable not set}"
: "${DB_NAME:?DB_NAME environment variable not set}"
: "${DB_USER:?DB_USER environment variable not set}"

# Ensure output directory exists
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Write header lines of docker-bake.hcl
cat <<EOF > "$OUTPUT_FILE"
APP_NAME = "${APP_NAME}"
FRAMEWORK = "${FRAMEWORK}"
APP_VERSION = "${APP_VERSION}"
DB_PASS = "${DB_PASS}"
DB_NAME = "${DB_NAME}"
DB_USER = "${DB_USER}"

group "default" {
  targets = ["app"$( [ "$WITH_DB" = "true" ] && echo ', "db"' )]
}

target "app" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["\${APP_NAME}-\${FRAMEWORK}:\${APP_VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    APP_NAME = APP_NAME
    APP_VERSION = APP_VERSION
  }
}
EOF

# Nếu WITH_DB = "true" thì thêm target db
if [ "$WITH_DB" = "true" ]; then
cat <<EOF >> "$OUTPUT_FILE"

target "db" {
  context = "."
  dockerfile = "Dockerfile.mysql"
  tags = ["\${APP_NAME}-mysql:\${APP_VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    MYSQL_ROOT_PASSWORD = DB_PASS
    MYSQL_DATABASE = DB_NAME
    MYSQL_USER = DB_USER
    MYSQL_PASSWORD = DB_PASS
  }
}
EOF
fi

sleep 1
exit 0
