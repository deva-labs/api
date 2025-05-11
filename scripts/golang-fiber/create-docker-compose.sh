#!/bin/bash
set -e

# --- Configuration ---
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
WORKDIR="public/${PROJECT_NAME}"
ENV_FILE="${BASE_DIR}/${WORKDIR}/.env"
COMPOSE_FILE="${WORKDIR}/docker-compose.yml"

log() {
  echo "[INFO] $1"
}

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
    log "Loaded environment variables from $ENV_FILE"
else
    ALT_ENV_FILE="./${WORKDIR}/.env"
    if [ -f "$ALT_ENV_FILE" ]; then
        ENV_FILE="$ALT_ENV_FILE"
        export $(grep -v '^#' "$ENV_FILE" | xargs)
        log "Loaded environment variables from fallback $ALT_ENV_FILE"
    else
        echo "❌ .env file not found at $ENV_FILE or fallback $ALT_ENV_FILE"
        exit 1
    fi
fi

# Validate required environment variables
: "${WITH_DB:?WITH_DB environment variable not set}"
: "${APP_NAME:?APP_NAME environment variable not set}"
: "${DB_PASS:?DB_PASS environment variable not set}"
: "${DB_NAME:?DB_NAME environment variable not set}"
: "${DB_USER:?DB_USER environment variable not set}"
: "${DB_PORT:?DB_PORT environment variable not set}"
: "${APP_PORT:?APP_PORT environment variable not set}"
: "${FRAMEWORK:?FRAMEWORK environment variable not set}"
: "${VERSION:?VERSION environment variable not set}"

# Ensure output directory exists
mkdir -p "$WORKDIR"

# --- Generate docker-compose.yml ---
if [ "$WITH_DB" = "yes" ]; then
  cat <<EOF > "$COMPOSE_FILE"
services:
  db:
    image: \${APP_NAME}-mysql:\${VERSION}
    container_name: \${APP_NAME}-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: \${DB_PASS}
      MYSQL_DATABASE: \${DB_NAME}
      MYSQL_USER: \${DB_USER}
      MYSQL_PASSWORD: \${DB_PASS}
    ports:
      - "\${DB_PORT}:3306"
    volumes:
      - db_data:/var/lib/mysql
    networks:
      - backend

  app:
    image: \${APP_NAME}-\${FRAMEWORK}:\${VERSION}
    container_name: \${APP_NAME}-\${FRAMEWORK}
    restart: always
    ports:
      - "\${APP_PORT}:\${APP_PORT}"
    volumes:
      - ./:/app
    environment:
      - DB_HOST=db
      - DB_PORT=3306
      - DB_NAME=\${DB_NAME}
      - DB_USER=\${DB_USER}
      - DB_PASS=\${DB_PASS}
    depends_on:
      - db
    networks:
      - backend

volumes:
  db_data:

networks:
  backend:
    name: \${APP_NAME}-network

x-bake:
  db:
    dockerfile: Dockerfile.mysql
    tags:
      - \${APP_NAME}-mysql:\${VERSION}
    platforms: ["linux/amd64"]
    cache-from: type=registry,ref=\${APP_NAME}-mysql:cache
    cache-to: type=registry,ref=\${APP_NAME}-mysql:cache,mode=max

  app:
    dockerfile: Dockerfile
    tags:
      - \${APP_NAME}-\${FRAMEWORK}:\${VERSION}
    platforms: ["linux/amd64"]
    cache-from: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache
    cache-to: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache,mode=max
EOF

else
  cat <<EOF > "$COMPOSE_FILE"
services:
  app:
    image: \${APP_NAME}-\${FRAMEWORK}:\${VERSION}
    container_name: \${APP_NAME}-\${FRAMEWORK}
    restart: always
    ports:
      - "\${APP_PORT}:\${APP_PORT}"
    volumes:
      - ./:/app
    environment:
      - DB_HOST=localhost
      - DB_PORT=\${DB_PORT}
      - DB_NAME=\${DB_NAME}
      - DB_USER=\${DB_USER}
      - DB_PASS=\${DB_PASS}
    networks:
      - backend

networks:
  backend:
    name: \${APP_NAME}-network

x-bake:
  app:
    dockerfile: Dockerfile
    tags:
      - \${APP_NAME}-\${FRAMEWORK}:\${VERSION}
    platforms: ["linux/amd64"]
    cache-from: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache
    cache-to: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache,mode=max
EOF
fi
