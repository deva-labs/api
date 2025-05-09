#!/bin/bash

set -e
if [ -f fiber-with-docker/.env ]; then
  export $(grep -v "^#" fiber-with-docker/.env | xargs)
fi

: "${APP_NAME:?APP_NAME enviroment variable is not set}"
: "${APP_PORT:?APP_PORT enviroment variable is not set}"
docker run -d -p ${APP_PORT}:${APP_PORT} --name ${APP_NAME} ${APP_NAME}:latest