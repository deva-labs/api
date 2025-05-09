#!/bin/bash

set -e

# Kiểm tra xem 'air' đã được cài đặt chưa
if command -v air &> /dev/null; then
    echo "✅ Air is already installed!"
    exit 0
fi

# Cài đặt Air (latest)
echo "⬇️ Installing Air..."
go install github.com/air-verse/air@latest

# Di chuyển binary tới /usr/local/bin (yêu cầu sudo)
echo "📦 Moving Air binary to /usr/local/bin..."
sudo mv "$(go env GOPATH)/bin/air" /usr/local/bin/

echo "✅ Air installed successfully!"

exit 0
