package routes

import (
	"dockerwizard-api/src/modules/projects/controllers"
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

	projectsRoutes := api.Group("projects")
	{
		projectsRoutes.Post("create", projects.CreateNewFiberProject)
	}

	// Testing Routes
	//testingRoutes := api.Group("/testing")
	//{
	//	testingRoutes.Post("/animation", tests.TestProgressBar60Seconds)
	//}
}
