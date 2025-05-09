#!/bin/bash

set -e

# Load environment variables from fiber-with-docker/.env
if [ -f fiber-with-docker/.env ]; then
  export $(grep -v "^#" fiber-with-docker/.env | xargs)
fi

: "${APP_NAME:?APP_NAME environment variable is not set}"

# Create .air.toml in the correct directory if it doesn't exist
[ -f fiber-with-docker/.air.toml ] || (cd fiber-with-docker && mkdir bin && air init)

# Update the Air config with custom values
sed -i '/^\[build\]/,/^\[.*\]/ s|^\( *cmd = \).*|\1"go build -buildvcs=false -o ./bin/'"${APP_NAME}"'-api ."|' fiber-with-docker/.air.toml
sed -i '/^\[build\]/,/^\[.*\]/ s|^\( *bin = \).*|\1"./bin/'"${APP_NAME}"'-api"|' fiber-with-docker/.air.toml
sed -i 's|^tmp_dir = .*|tmp_dir = "bin"|' fiber-with-docker/.air.toml

exit 0
