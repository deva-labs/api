#!/bin/bash

set -e

if [ ! -f Makefile ]; then
  echo "❌ Makefile not found in the current directory."
  exit 1
fi

echo "🚀 Running Makefile 'install' target..."
echo "---------------------------------------"

chmod +x ./scripts/*.sh

make install

echo "✅ Project setup completed successfully."
