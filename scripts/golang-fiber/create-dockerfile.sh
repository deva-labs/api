#!/bin/bash
set -e

# --- Configuration ---
BASE_DIR="/app"
WORKDIR="public/${PROJECT_NAME}"
ENV_FILE="${BASE_DIR}/${WORKDIR}/.env"
ALT_ENV_FILE="./${WORKDIR}/.env"
MAIN_DOCKERFILE="${WORKDIR}/Dockerfile"

# DB-specific Dockerfiles
MYSQL_DOCKERFILE="${WORKDIR}/Dockerfile.mysql"
POSTGRES_DOCKERFILE="${WORKDIR}/Dockerfile.postgres"
MONGODB_DOCKERFILE="${WORKDIR}/Dockerfile.mongodb"

# --- Create working directory if it doesn't exist ---
mkdir -p "$WORKDIR"

# --- Load .env ---
if [ -f "$ENV_FILE" ]; then
    SOURCE_ENV="$ENV_FILE"
elif [ -f "$ALT_ENV_FILE" ]; then
    SOURCE_ENV="$ALT_ENV_FILE"
else
    echo "❌ .env file not found"
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
: "${AIR_VERSION:?❌ AIR_VERSION environment variable not set}"
: "${DB_TYPE:?❌ DB_TYPE environment variable not set (mysql, postgres, mongodb)}"
: "${DB_VERSION:?❌ DB_VERSION environment variable not set}"
: "${WITH_DB:?❌ WITH_DB environment variable not set}"

# --- Clean old Dockerfiles ---
rm -f "$MAIN_DOCKERFILE" "$MYSQL_DOCKERFILE" "$POSTGRES_DOCKERFILE" "$MONGODB_DOCKERFILE"

# --- Initialize DB client package variable empty ---
DB_CLIENT_PKG=""

# --- Generate DB Dockerfile and set DB client package only if WITH_DB = "true" ---
if [ "$WITH_DB" = "true" ]; then
  case "$DB_TYPE" in
    mysql)
      cat <<EOF > "$MYSQL_DOCKERFILE"
FROM mysql:${DB_VERSION}
EOF
      DB_CLIENT_PKG="mariadb-client"
      ;;
    postgres)
      cat <<EOF > "$POSTGRES_DOCKERFILE"
FROM postgres:${DB_VERSION}
EOF
      DB_CLIENT_PKG="postgresql-client"
      ;;
    mongodb)
      cat <<EOF > "$MONGODB_DOCKERFILE"
FROM mongo:${DB_VERSION}
EOF
      DB_CLIENT_PKG="mongodb-tools"
      ;;
    *)
      echo "❌ Unsupported DB_TYPE: $DB_TYPE"
      exit 1
      ;;
  esac
fi

# --- Generate main Dockerfile ---
COMMON_PKGS="bash"

if [ "$ENV" = "dev" ]; then
  cat <<EOF > "$MAIN_DOCKERFILE"
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk add --no-cache \\
EOF

  # Nếu có DB client package, thêm vào lệnh apk add
  if [ -n "$DB_CLIENT_PKG" ]; then
    echo "    $DB_CLIENT_PKG \\" >> "$MAIN_DOCKERFILE"
  fi

  # Tiếp tục các gói còn lại
  cat <<EOF >> "$MAIN_DOCKERFILE"
    inotify-tools \\
    $COMMON_PKGS

RUN go install github.com/air-verse/air@${AIR_VERSION} && \\
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
EOF

  if [ -n "$DB_CLIENT_PKG" ]; then
    echo "    $DB_CLIENT_PKG \\" >> "$MAIN_DOCKERFILE"
  fi

  cat <<EOF >> "$MAIN_DOCKERFILE"
    $COMMON_PKGS

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
