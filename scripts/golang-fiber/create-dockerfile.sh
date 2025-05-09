#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
    export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

# Remove old Dockerfile
rm -f fiber-with-docker/Dockerfile

: "${ENV:?ENV environment variable not set}"
# Build Dockerfile content
if [ "$ENV" = "dev" ]; then
    cat <<EOF > fiber-with-docker/Dockerfile
FROM golang:${GO_VERSION}
RUN apt update && apt install -y default-mysql-client
RUN go install github.com/air-verse/air@latest && mv \$(go env GOPATH)/bin/air /usr/local/bin/
WORKDIR /app
COPY .air.toml ./
COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . ./
RUN air init || true
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/bin/sh", "-c", "/app/entrypoint.sh"]
EOF
else
    cat <<EOF > fiber-with-docker/Dockerfile
FROM golang:${GO_VERSION}
RUN apt update && apt install -y default-mysql-client
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . ./
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/bin/sh", "-c", "/app/entrypoint.sh"]
EOF
fi

exit 0
