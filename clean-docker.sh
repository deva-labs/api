#!/bin/bash

# Script to clean all Docker artifacts
echo "=== Starting Docker Cleanup ==="

# Stop all running containers
echo "Stopping all containers..."
docker stop $(docker ps -aq) 2>/dev/null || echo "No containers to stop"

# Remove all containers
echo "Removing all containers..."
docker rm $(docker ps -aq) 2>/dev/null || echo "No containers to remove"

# Remove all Docker images
echo "Removing all images..."
docker rmi $(docker images -q) 2>/dev/null || echo "No images to remove"

# Remove all Docker networks (except predefined ones)
echo "Removing custom networks..."
docker network prune -f

# Remove all Docker volumes
echo "Force removing all volumes (including attached)..."
docker volume rm $(docker volume ls -q) 2>/dev/null || echo "No volumes to remove"

# Remove all Docker build cache
echo "Removing build cache..."
docker builder prune -af

# Clean up dangling images
echo "Removing dangling images..."
docker image prune -f

# For more aggressive cleanup (uncomment if needed)
# echo "Removing all unused objects..."
# docker system prune -af --volumes

echo "=== Docker Cleanup Complete ==="
echo "Current Docker status:"
docker system df
