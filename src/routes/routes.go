package routes

import (
	"github.com/gofiber/fiber/v2"
	"skypipe/src/middlewares"
	key_token "skypipe/src/modules/key_token/controllers"
	"skypipe/src/modules/projects/controllers"
	users "skypipe/src/modules/users/controllers"
	verifications "skypipe/src/modules/verifications/controllers"
	"skypipe/src/utils"
)

// RegisterRoutes is a Router Controller
func RegisterRoutes(app *fiber.App) {
	// Middlewares
	authMiddleware := middlewares.AuthMiddleware
	authzMiddleware := middlewares.Authorization
	needPermission := utils.Permissions
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

	// Authentication Routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.Get("user-info", authMiddleware(), authzMiddleware(needPermission["USER_READ"]), users.GetUser) // 1
		authRoutes.Post("register", users.Register)                                                                // 2
		authRoutes.Post("login", users.Login)                                                                      // 3
		authRoutes.Post("verifications", verifications.VerifyCodeAndGenerateToken)                                 // 6
		authRoutes.Post("refresh-access-token", key_token.RefreshAccessToken)                                      // 7
		authRoutes.Post("forgot-password", users.ForgotPassword)                                                   // 8
		authRoutes.Post("confirm-forgot-password", verifications.VerifyCode)                                       // 9
		authRoutes.Post("change-password", users.ChangePassword)                                                   // 10
		authRoutes.Post("reset-password", users.RenewPassword)
	}

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
