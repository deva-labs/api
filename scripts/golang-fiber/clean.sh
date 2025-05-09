#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
  export $(grep -v "^#" fiber-with-docker/.env | xargs)
fi
: "${APP_NAME:?APP_NAME enviroment variable is not set}"

# 1. Dừng và xoá tất cả các container Docker
echo "🔴 Stopping and removing Docker containers..."
cd fiber-with-docker
docker compose down

# 2. Xoá image Docker với tên $APP_NAME:latest
echo "🧹 Removing Docker image $APP_NAME:latest..."
docker rmi "$APP_NAME:latest" || echo "⚠️ Docker image $APP_NAME:latest does not exist or is in use."

# 3. Nén thư mục public thành file zip và lưu file zip trong thư mục hiện tại
echo "📦 Zipping public folder to ./public.zip..."
cd ..
zip -r "./public/fiber-with-docker.zip" "./fiber-with-docker"

# 4. Xoá thư mục làm việc fiber-with-docker
echo "🧹 Deleting working directory fiber-with-docker..."
rm -rf fiber-with-docker

echo "✅ All tasks completed successfully!"
