SHELL := /bin/bash
FRAMEWORK := golang-fiber

# Paths
SCRIPTS := ./scripts/$(FRAMEWORK)

# Targets
.PHONY: all init build \
        create-all install-all docker-all \
        clean help

all: init build

## Initialization Phase
init: create-all install-all

create-all: create-env create-dockerfile create-docker-compose create-entrypoint create-docker-bake create-main

install-all: install-go init-go-modules install-air air-init

create-env:
	@bash $(SCRIPTS)/create-env.sh

create-dockerfile:
	@bash $(SCRIPTS)/create-dockerfile.sh

create-docker-compose:
	@bash $(SCRIPTS)/create-docker-compose.sh

create-entrypoint:
	@bash $(SCRIPTS)/create-entrypoint.sh

create-docker-bake:
	@bash $(SCRIPTS)/create-docker-bake.sh

create-main:
	@bash $(SCRIPTS)/create-main.sh

install-go:
	@bash $(SCRIPTS)/install-golang.sh

init-go-modules:
	@bash $(SCRIPTS)/init-golang.sh

install-air:
	@bash $(SCRIPTS)/install-air.sh

air-init:
	@bash $(SCRIPTS)/init-air.sh

## Build Phase
build: docker-all clean

docker-all: docker-build docker-run docker-compose-up

docker-build:
	@bash $(SCRIPTS)/docker-build.sh

docker-run:
	@bash $(SCRIPTS)/docker-run.sh

docker-compose-up:
	@bash $(SCRIPTS)/docker-compose-up.sh

## Cleanup
clean:
	@bash $(SCRIPTS)/clean.sh

## Help
help:
	@echo "🛠️  Docker Project Management"
	@echo ""
	@echo "Initialization:"
	@echo "  init              Initialize project (env, Dockerfile, entrypoint, etc)"
	@echo ""
	@echo "Build Options:"
	@echo "  build             Full build (image + run + compose)"
	@echo ""
	@echo "Usage:"
	@echo "  make build"
	@echo "  make init"
	@echo "  make clean"
