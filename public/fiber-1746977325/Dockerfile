FROM golang:1.24.2-alpine

RUN apk update && apk add --no-cache \
    mariadb-client \
    inotify-tools \
    bash

RUN go install github.com/air-verse/air@latest && \
    mv "$(go env GOPATH)/bin/air" /usr/local/bin/

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

COPY .air.toml ./
COPY . .

ENTRYPOINT ["/bin/sh", "entrypoint.sh"]
