#!/bin/bash

set -e

if [ -f fiber-with-docker/.env ]; then
  export $(grep -v "^#" fiber-with-docker/.env | xargs)
fi

: "${APP_PORT:?APP_PORT environment variable is not set}"

cat <<EOF > fiber-with-docker/main.go
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
