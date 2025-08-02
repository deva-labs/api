package routes

import (
	"deva/src/middlewares"
	captcha "deva/src/modules/captcha/controllers"
	key_token "deva/src/modules/key_token/controllers"
	"deva/src/modules/projects/controllers"
	users "deva/src/modules/users/controllers"
	verifications "deva/src/modules/verifications/controllers"
	VideoControllers "deva/src/modules/videos/controllers"
	"deva/src/utils"
	"github.com/gofiber/fiber/v2"
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
	authRoutes := api.Group("auth")
	{
		authRoutes.Get("user-info", authMiddleware(), authzMiddleware(needPermission["USER_READ"]), users.GetUser) // 1
		authRoutes.Post("register", users.Register)                                                                // 2
		authRoutes.Post("login", users.Login)                                                                      // 3
		authRoutes.Post("verifications", verifications.VerifyCodeAndGenerateToken)                                 // 6
		authRoutes.Post("refresh-access-token", key_token.RefreshAccessToken)                                      // 7
		authRoutes.Post("forgot-password", users.ForgotPassword)                                                   // 8
		authRoutes.Post("verify-forgot-password", verifications.VerifyCodeAndSetPasswordToken)                     // 9
		authRoutes.Post("change-password", users.ChangePassword)                                                   // 10
		authRoutes.Post("reset-password", users.RenewPassword)
	}

	captchaRoutes := api.Group("captcha")
	{
		captchaRoutes.Post("create", captcha.CreateCaptcha)
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

	// Videos Routes
	videos := api.Group("/videos")
	{
		videos.Post("", VideoControllers.ListAll)
		videos.Post("newest", VideoControllers.ListNewest)
		videos.Post("details", VideoControllers.Details)
	}
}
