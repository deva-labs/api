#!/bin/bash
set -e

# --- Configuration ---
WORKDIR="public/${PROJECT_NAME}"
ENV_FILE="${WORKDIR}/.env"
COMPOSE_FILE="${WORKDIR}/docker-compose.yml"

# --- Load .env file ---
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

# --- Validate required environment variables ---
: "${WITH_DB:?WITH_DB environment variable not set}"
: "${DB_TYPE:?DB_TYPE environment variable not set (e.g. mysql, postgres, mongodb)}"
: "${APP_NAME:?APP_NAME environment variable not set}"
: "${DB_NAME:?DB_NAME environment variable not set}"
: "${DB_USER:?DB_USER environment variable not set}"
: "${DB_PASS:?DB_PASS environment variable not set}"
: "${DB_PORT:?DB_PORT environment variable not set}"
: "${APP_PORT:?APP_PORT environment variable not set}"
: "${FRAMEWORK:?FRAMEWORK environment variable not set}"
: "${APP_VERSION:?APP_VERSION environment variable not set}"

# --- Ensure output directory exists ---
mkdir -p "$WORKDIR"

# --- Generate DB service ---
generate_db_service() {
  case "$DB_TYPE" in
    mysql)
      cat <<EOF
  db:
    image: \${APP_NAME}-mysql:\${APP_VERSION}
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
EOF
      ;;
    postgres)
      cat <<EOF
  db:
    image: \${APP_NAME}-postgres:\${APP_VERSION}
    container_name: \${APP_NAME}-postgres
    restart: always
    environment:
      POSTGRES_DB: \${DB_NAME}
      POSTGRES_USER: \${DB_USER}
      POSTGRES_PASSWORD: \${DB_PASS}
    ports:
      - "\${DB_PORT}:5432"
    volumes:
      - db_data:/var/lib/postgresql/data
    networks:
      - backend
EOF
      ;;
    mongodb)
      cat <<EOF
  db:
    image: \${APP_NAME}-mongodb:\${APP_VERSION}
    container_name: \${APP_NAME}-mongodb
    restart: always
    environment:
      MONGO_INITDB_DATABASE: \${DB_NAME}
      MONGO_INITDB_ROOT_USERNAME: \${DB_USER}
      MONGO_INITDB_ROOT_PASSWORD: \${DB_PASS}
    ports:
      - "\${DB_PORT}:27017"
    volumes:
      - db_data:/data/db
    networks:
      - backend
EOF
      ;;
    *)
      echo "❌ Unsupported DB_TYPE: $DB_TYPE"
      exit 1
      ;;
  esac
}

# --- Compose file generation ---
{
  echo "services:"

  if [ "$WITH_DB" = "true" ]; then
    generate_db_service
  fi

  # App service
  cat <<EOF
  app:
    image: \${APP_NAME}-\${FRAMEWORK}:\${APP_VERSION}
    container_name: \${APP_NAME}-\${FRAMEWORK}
    restart: always
    ports:
      - "\${APP_PORT}:\${APP_PORT}"
    volumes:
      - ./:/app
EOF

  if [ "$WITH_DB" = "true" ]; then
    cat <<EOF
    environment:
      - DB_HOST=db
      - DB_PORT=\${DB_PORT}
      - DB_NAME=\${DB_NAME}
      - DB_USER=\${DB_USER}
      - DB_PASS=\${DB_PASS}
    depends_on:
      - db
EOF
  fi

  cat <<EOF
    networks:
      - backend

volumes:
  db_data:

networks:
  backend:
    name: \${APP_NAME}-network

x-bake:
EOF

  if [ "$WITH_DB" = "true" ]; then
    cat <<EOF
  db:
    dockerfile: Dockerfile.${DB_TYPE}
    tags:
      - \${APP_NAME}-${DB_TYPE}:\${APP_VERSION}
    platforms: ["linux/amd64"]
    cache-from: type=registry,ref=\${APP_NAME}-${DB_TYPE}:cache
    cache-to: type=registry,ref=\${APP_NAME}-${DB_TYPE}:cache,mode=max
EOF
  fi

  cat <<EOF
  app:
    dockerfile: Dockerfile
    tags:
      - \${APP_NAME}-\${FRAMEWORK}:\${APP_VERSION}
    platforms: ["linux/amd64"]
    cache-from: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache
    cache-to: type=registry,ref=\${APP_NAME}-\${FRAMEWORK}:cache,mode=max
EOF

} > "$COMPOSE_FILE"
