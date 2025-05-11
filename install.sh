#!/bin/bash

set -e

PROJECT_NAME="$1"

if [ ! -f Makefile ]; then
  echo "❌ Makefile not found in the current directory."
  exit 1
fi

if [ -t 1 ]; then
    # We're in a terminal, enable colors and animations
    export TERM=xterm-256color
else
    # Not in a terminal, disable animations
    export NO_COLOR=1
fi

echo "🚀 Running Makefile 'init' target..."
echo "---------------------------------------"

chmod +x ./scripts/*.sh

export PROJECT_NAME="$PROJECT_NAME"

make init

echo "✅ Project setup completed successfully."
