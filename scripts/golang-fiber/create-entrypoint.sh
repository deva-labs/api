#!/bin/bash

set -e

# Load environment variables if .env exists
if [ -f fiber-with-docker/.env ]; then
    export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

: "${APP_NAME:?APP_NAME environment variable not set}"
: "${ENV:?ENV environment variable not set}"

# Clean old file if exists
rm -f fiber-with-docker/entrypoint.sh

if [ "$ENV" = "dev" ]; then
cat <<EOF > fiber-with-docker/entrypoint.sh
#!/bin/bash
set -e

BINARY="./bin/${APP_NAME}"

if [ ! -f "\$BINARY" ]; then
  echo "Binary \$BINARY not found. Building it now..."
  go build -buildvcs=false -o \$BINARY .
fi

echo "Starting the application with Air..."
exec air
EOF

elif [ "$ENV" = "prod" ]; then
cat <<EOF > fiber-with-docker/entrypoint.sh
#!/bin/bash
set -e

BINARY="./${APP_NAME}"

if [ ! -f "\$BINARY" ]; then
  echo "Binary \$BINARY not found. Building it now..."
  go build -buildvcs=false -o \$BINARY .
fi

echo "Starting the application..."
exec \$BINARY
EOF
fi

