package key_token

import (
	"github.com/gofiber/fiber/v2"
	key_token "skypipe/src/modules/key_token/services"
	"strings"
)

func RefreshAccessToken(c *fiber.Ctx) error {
	authHeader := c.Get("x-rtoken-id")
	clientID := c.Get("x-client-id")

	if authHeader == "" || clientID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "missing token or client id on header",
			},
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "invalid token header format",
			},
		})
	}

	accessToken := tokenParts[1]
	response, err := key_token.RefreshAccessToken(accessToken, clientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Failed to get user",
			},
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Verification successful",
		},
		"data": response,
	})
}
