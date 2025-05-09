#!/bin/bash

set -e

rm -f .env

mkdir -p fiber-with-docker

cat <<EOF > fiber-with-docker/.env
GO_VERSION=1.24.2
APP_NAME=dockerwizard-api
GO_TAR=go\${GO_VERSION}.linux-amd64.tar.gz
WITH_DB=yes
RUN_WITH_DOCKER_COMPOSE=yes
ENV=dev
IMAGE_NAME=\${APP_NAME}:latest
DB_PASS=dockerwizard-admin
DB_NAME=dockerwizard
DB_USER=admin
DB_PORT=3307
APP_PORT=2350
EOF
