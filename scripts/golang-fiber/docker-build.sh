#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
    export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

# Verify APP_NAME is set
: "${APP_NAME:?APP_NAME environment variable not set}"
docker build -t "${APP_NAME}:latest" .