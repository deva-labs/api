package routes

import (
	"dockerwizard-api/src/modules/fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes is a Router Controller
func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Health Check Routes
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ready"})
	})

	newProject := api.Group("/create-project")
	{
		newProject.Post("", controllers.CreateNewFiberProject)
	}
}
