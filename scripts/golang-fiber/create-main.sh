#!/bin/bash

set -e

# Define fallback logic for ENV and MAIN file
if [ -d "/app/public/${PROJECT_NAME}" ]; then
  BASE_DIR="/app"
else
  BASE_DIR="."
fi

ENV_PATH="${BASE_DIR}/public/${PROJECT_NAME}/.env"
MAIN_FILE="${BASE_DIR}/public/${PROJECT_NAME}/main.go"

# Load .env file
if [ -f "$ENV_PATH" ]; then
  export $(grep -v "^#" "$ENV_PATH" | xargs)
else
  exit 1
fi

# Ensure APP_PORT is set
: "${APP_PORT:?❌ APP_PORT environment variable is not set}"

# Ensure directory exists
mkdir -p "$(dirname "$MAIN_FILE")"

# Write the Go application
cat <<EOF > "$MAIN_FILE"
package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":${APP_PORT}")
}
EOF

sleep 1
exit 0