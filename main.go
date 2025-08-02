package main

import (
	"deva/src/config"
	"deva/src/routes"
	"deva/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠️ Could not load .env file: %v", err)
	} else {
		log.Println("✅ .env file loaded successfully")
	}
	app := fiber.New()

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "*",
		ExposeHeaders:    "*",
		AllowCredentials: false,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	app.Static("/public", "./public", fiber.Static{
		Browse:        true,
		CacheDuration: 10 * time.Minute,
	})

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Middleware to upgrade only if WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket handler
	app.Get("/ws", services.WebSocketUpgrader())

	// Connect to db
	config.ConnectDatabase()
	// Connect to redis
	config.ConnectRedis()
	// Register other routes
	routes.RegisterRoutes(app)
	err := app.Listen(":2350")
	if err != nil {
		return
	}
}
