#!/bin/bash

set -e


if [ -f fiber-with-docker/.env ]; then
    export $(grep -v '^#' fiber-with-docker/.env | xargs)
fi

# Verify APP_NAME is set
: "${APP_NAME:?APP_NAME environment variable not set}"

# Change into project directory
cd fiber-with-docker

# Initialize Go module only if it doesn't exist
if [ ! -f go.mod ]; then
    go mod init "$APP_NAME"
fi

# Tidy up dependencies (creates go.sum if needed)
go mod tidy

exit 0
