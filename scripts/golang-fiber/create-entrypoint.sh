#!/bin/bash

set -e

# Support dynamic folder via PROJECT_NAME
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
else
    exit 1
fi

export $(grep -v '^#' "$ENV_FILE" | xargs)

# Check required variables
: "${APP_NAME:?❌ APP_NAME environment variable not set}"
: "${ENV:?❌ ENV environment variable not set}"

# Remove old entrypoint if it exists
rm -f "$ENTRYPOINT"

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

sleep 1
exit 0