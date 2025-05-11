#!/bin/bash

set -e

# Logging helper
log() {
  echo "[INFO] $1"
}

# Support dynamic folder via PROJECT_NAME
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
WORKDIR="public/${PROJECT_NAME}"

# Try absolute path first, then fall back to relative path
ENV_FILE_ABSOLUTE="${BASE_DIR}/${WORKDIR}/.env"
ENV_FILE_RELATIVE="${WORKDIR}/.env"
ENTRYPOINT="${WORKDIR}/entrypoint.sh"

# Check which .env file exists
if [ -f "$ENV_FILE_ABSOLUTE" ]; then
    ENV_FILE="$ENV_FILE_ABSOLUTE"
elif [ -f "$ENV_FILE_RELATIVE" ]; then
    ENV_FILE="$ENV_FILE_RELATIVE"
    log "⚠️  Using relative path for .env file as absolute path not found"
else
    echo "❌ .env not found at either $ENV_FILE_ABSOLUTE or $ENV_FILE_RELATIVE"
    exit 1
fi

log "📂 Working directory: $WORKDIR"
log "📄 Loading env from: $ENV_FILE"

export $(grep -v '^#' "$ENV_FILE" | xargs)
log ".env loaded successfully"

# Check required variables
: "${APP_NAME:?❌ APP_NAME environment variable not set}"
: "${ENV:?❌ ENV environment variable not set}"

# Remove old entrypoint if it exists
rm -f "$ENTRYPOINT"

log "⚙️ Generating entrypoint for ENV=$ENV"

if [ "$ENV" = "dev" ]; then
cat <<EOF > "$ENTRYPOINT"
#!/bin/sh
set -e

BINARY="./bin/${APP_NAME}"

if [ ! -f "\$BINARY" ]; then
  echo "Binary \$BINARY not found. Building it now..."
  go build -buildvcs=false -o "\$BINARY" .
fi

echo "Starting the application with Air..."
exec air
EOF

elif [ "$ENV" = "prod" ]; then
cat <<EOF > "$ENTRYPOINT"
#!/bin/sh
set -e

BINARY="./${APP_NAME}"

if [ ! -f "\$BINARY" ]; then
  echo "Binary \$BINARY not found. Building it now..."
  go build -buildvcs=false -o "\$BINARY" .
fi

echo "Starting the application..."
exec "\$BINARY"
EOF
else
  echo "❌ Unsupported ENV: $ENV. Must be 'dev' or 'prod'."
  exit 1
fi

chmod +x "$ENTRYPOINT"
log "✅ Entrypoint script created at $ENTRYPOINT"