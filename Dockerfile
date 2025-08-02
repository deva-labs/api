FROM golang:1.24.2-alpine

# Cài dependencies (thêm coreutils và tini để đảm bảo terminal tương tác)
RUN apk update && apk add --no-cache \
    postgresql-client \
    inotify-tools \
    curl \
    make \
    bash \
    git \
    docker-cli \
    zip \
    coreutils \
    tini \
    font-dejavu

# Cài air
RUN go install github.com/air-verse/air@latest && \
    mv "$(go env GOPATH)/bin/air" /usr/local/bin/

# Docker CLI, Buildx, Compose
ENV DOCKER_VERSION=28.1.1
ENV BUILDX_VERSION=v0.23.0
ENV COMPOSE_VERSION=v2.35.1

RUN mkdir -p /usr/local/bin && \
    curl -fsSL "https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz" -o docker.tgz && \
    tar xzvf docker.tgz --strip-components=1 -C /usr/local/bin && \
    rm docker.tgz

RUN mkdir -p /root/.docker/cli-plugins && \
    curl -SL "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION}/buildx-${BUILDX_VERSION}.linux-amd64" -o /root/.docker/cli-plugins/docker-buildx && \
    chmod +x /root/.docker/cli-plugins/docker-buildx && \
    curl -SL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-linux-x86_64" -o /root/.docker/cli-plugins/docker-compose && \
    chmod +x /root/.docker/cli-plugins/docker-compose

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY entrypoint.sh ./
RUN chmod +x entrypoint.sh

COPY .air.toml ./
COPY . .

# 💡 Sử dụng tini để đảm bảo terminal signals/TTY hoạt động đúng
ENTRYPOINT ["/sbin/tini", "--", "/bin/sh", "entrypoint.sh"]
