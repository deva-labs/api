#!/bin/bash
set -e

# --- Configuration ---
PROJECT_NAME="${PROJECT_NAME:-fiber-with-docker}"
BASE_DIR="/app"
WORKDIR="public/${PROJECT_NAME}"
ENV_FILE="${BASE_DIR}/${WORKDIR}/.env"
ALT_ENV_FILE="./${WORKDIR}/.env"
MAIN_DOCKERFILE="${WORKDIR}/Dockerfile"
MYSQL_DOCKERFILE="${WORKDIR}/Dockerfile.mysql"

# --- Create working directory if it doesn't exist ---
mkdir -p "$WORKDIR"

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    SOURCE_ENV="$ENV_FILE"
elif [ -f "$ALT_ENV_FILE" ]; then
    SOURCE_ENV="$ALT_ENV_FILE"
else
    echo "❌ .env file not found at $ENV_FILE or fallback $ALT_ENV_FILE"
    exit 1
fi

while IFS='=' read -r key value; do
    [[ $key =~ ^#.*$ || -z $key ]] && continue
    value=$(echo "$value" | sed "s/^['\"]//;s/['\"]$//")
    export "$key"="$value"
done < "$SOURCE_ENV"

# --- Verify required variables ---
: "${ENV:?❌ ENV environment variable not set}"
: "${GO_VERSION:?❌ GO_VERSION environment variable not set}"

# --- Clean old Dockerfiles ---
rm -f "$MAIN_DOCKERFILE" "$MYSQL_DOCKERFILE"

# --- Generate MySQL Dockerfile ---
cat <<EOF > "$MYSQL_DOCKERFILE"
FROM mysql:latest
EOF

# --- Generate main Dockerfile ---

if [ "$ENV" = "dev" ]; then
  cat <<EOF > "$MAIN_DOCKERFILE"
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk add --no-cache \\
    mariadb-client \\
    inotify-tools \\
    bash

RUN go install github.com/air-verse/air@latest && \\
    mv "\$(go env GOPATH)/bin/air" /usr/local/bin/

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

COPY .air.toml ./
COPY . .

ENTRYPOINT ["/bin/sh", "entrypoint.sh"]
EOF

else
  cat <<EOF > "$MAIN_DOCKERFILE"
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk add --no-cache \\
    mariadb-client \\
    bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

COPY . .

ENTRYPOINT ["/bin/sh", "entrypoint.sh"]
EOF
fi

sleep 1
exit 0