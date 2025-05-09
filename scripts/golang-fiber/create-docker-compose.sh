#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
  export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

rm -f fiber-with-docker/docker-compose.yml

: "${WITH_DB:?WITH_DB environment variable not set}"
: "${APP_NAME:?APP_NAME environment variable not set}"
: "${DB_PASS:?DB_PASS environment variable not set}"
: "${DB_NAME:?DB_NAME environment variable not set}"
: "${DB_USER:?DB_USER environment variable not set}"
: "${DB_PORT:?DB_PORT environment variable not set}"
: "${APP_PORT:?APP_PORT environment variable not set}"
: "${FRAMEWORK:=golang}"

if [ "$WITH_DB" = "yes" ]; then
  cat <<EOF > fiber-with-docker/docker-compose.yml
version: "3.8"

services:
  db:
    image: mysql:latest
    container_name: ${APP_NAME}-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASS}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASS}
    ports:
      - "${DB_PORT}:3306"
    volumes:
      - db_data:/var/lib/mysql
    networks:
      - ${APP_NAME}-network

  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ${APP_NAME}-${FRAMEWORK}
    restart: always
    ports:
      - "${APP_PORT}:${APP_PORT}"
    environment:
      - DB_USER=${DB_USER}
      - DB_PASS=${DB_PASS}
      - DB_HOST=db
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
    volumes:
      - ./:/app
    depends_on:
      - db

volumes:
  db_data:

networks:
  ${APP_NAME}-network:
EOF

else
  cat <<EOF > fiber-with-docker/docker-compose.yml
version: "3.8"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ${APP_NAME}-${FRAMEWORK}
    restart: always
    ports:
      - "${APP_PORT}:${APP_PORT}"
    environment:
      - DB_USER=${DB_USER}
      - DB_PASS=${DB_PASS}
      - DB_HOST=localhost
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
    volumes:
      - ./:/app

volumes:
  db_data:

networks:
  ${APP_NAME}-network:
EOF
fi
