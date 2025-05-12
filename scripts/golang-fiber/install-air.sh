#!/bin/bash

set -e
if command -v air &> /dev/null; then
    echo "✅ Air is already installed!"
    exit 0
fi

go install github.com/air-verse/air@latest

sudo mv "$(go env GOPATH)/bin/air" /usr/local/bin/

sleep 1
exit 0