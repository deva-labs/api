# Include spinner utils
UTILS := ./scripts/utils.sh
SHELL := /bin/bash
FRAMEWORK := golang-fiber
export


# Targets
.PHONY: all init build \
        create-env create-dockerfile create-docker-compose create-entrypoint \
        create-docker-bake create-main install-go init-go-modules install-air air-init \
        docker-build docker-run docker-compose-up clean help

all: init build

## Initialization Phase - Full setup as original
init: create-env create-dockerfile create-docker-compose create-entrypoint \
      create-docker-bake create-main install-go init-go-modules install-air air-init

create-env:
	@source $(UTILS) && with_progress_bar ".env" "./scripts/$(FRAMEWORK)/create-env.sh" "creating"

create-dockerfile:
	@source $(UTILS) && with_progress_bar "Dockerfile" "./scripts/$(FRAMEWORK)/create-dockerfile.sh" "creating"

create-docker-compose:
	@source $(UTILS) && with_progress_bar "docker-compose.yml" "./scripts/$(FRAMEWORK)/create-docker-compose.sh" "creating"

create-entrypoint:
	@source $(UTILS) && with_progress_bar "entrypoint.sh" "./scripts/$(FRAMEWORK)/create-entrypoint.sh" "creating"

create-docker-bake:
	@source $(UTILS) && with_progress_bar "docker-bake.hcl" "./scripts/$(FRAMEWORK)/create-docker-bake.sh" "creating"

create-main:
	@source $(UTILS) && with_progress_bar "main.go" "./scripts/$(FRAMEWORK)/create-main.sh" "creating"

install-go:
	@source $(UTILS) && with_progress_bar "Golang" "./scripts/$(FRAMEWORK)/install-golang.sh" "installing"

init-go-modules:
	@source $(UTILS) && with_progress_bar "Go module" "./scripts/$(FRAMEWORK)/init-golang.sh" "initializing"

install-air:
	@source $(UTILS) && with_progress_bar "Air" "./scripts/$(FRAMEWORK)/install-air.sh" "installing"

air-init:
	@source $(UTILS) && with_progress_bar "Air config" "./scripts/$(FRAMEWORK)/init-air.sh" "initializing"

## Build Phase - Complete workflow
build: docker-compose-up clean

docker-build:
	@source $(UTILS) && with_progress_bar "Docker image (Remote)" "./scripts/$(FRAMEWORK)/docker-build.sh" "building"

docker-run:
	@source $(UTILS) && with_progress_bar "Container (Remote)" "./scripts/$(FRAMEWORK)/docker-run.sh" "running"

docker-compose-up:
	@source $(UTILS) && with_progress_bar "Compose (Remote)" "./scripts/$(FRAMEWORK)/docker-compose-up.sh" "starting"

## Cleanup
clean:
	@source $(UTILS) && with_progress_bar "Cleanup" "./scripts/$(FRAMEWORK)/clean.sh" "cleaning"

## Help
help:
	@echo "Docker Project Management"
	@echo ""
	@echo "Initialization:"
	@echo "  init                    Full initialization (env, configs, dependencies)"
	@echo ""
	@echo "Build Options:"
	@echo "  build                   Complete build workflow (build+run+compose)"
	@echo "  docker-build            Build Docker image (auto-detects local/remote)"
	@echo "  docker-run              Run container (auto-detects local/remote)"
	@echo "  docker-compose-up       Start with compose (auto-detects local/remote)"
	@echo ""
	@echo "Configuration:"
	@echo "  DOCKER_HOST            Set remote Docker host (default: tcp://localhost:2376)"
	@echo "  DOCKER_TLS_VERIFY      Enable TLS verification (default: 1)"
	@echo "  DOCKER_CERT_PATH       Path to TLS certs (default: /etc/docker/certs)"
	@echo ""
	@echo "Example remote build:"
	@echo "  make build DOCKER_HOST=tcp://my-remote-host:2376"
	@echo ""
	@echo "Example local build:"
	@echo "  make build DOCKER_TLS_VERIFY="