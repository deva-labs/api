#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
    export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

: "${RUN_WITH_DOCKER_COMPOSE:?RUN_WITH_DOCKER_COMPOSE enviroment variable is not set}"

cd fiber-with-docker

if [ "$RUN_WITH_DOCKER_COMPOSE" == "yes" ]; then
    echo "🚀 Running docker compose build and up..."
    docker compose build
    docker compose up -d
else
    echo "🚀 Running docker build and docker run..."
    make docker-build
    make docker-run
fi
