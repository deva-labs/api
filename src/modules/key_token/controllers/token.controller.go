package key_token

import (
	"deva/src/lib/interfaces"
	key_token "deva/src/modules/key_token/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RefreshAccessToken(c *fiber.Ctx) error {
	authHeader := c.Get("x-rtoken-id")
	clientID := c.Get("x-client-id")

	if authHeader == "" || clientID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    fiber.StatusUnauthorized,
				Message: "missing token or client id on header",
			},
			Error: nil,
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    fiber.StatusUnauthorized,
				Message: "invalid token header format",
			},
			Error: nil,
		})
	}

	accessToken := tokenParts[1]
	response, serviceErr := key_token.RefreshAccessToken(accessToken, clientID)
	if serviceErr != nil {
		s := serviceErr.Err.Error()
		errStr := &s
		return c.Status(serviceErr.StatusCode).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    serviceErr.StatusCode,
				Message: serviceErr.Message,
			},
			Error: errStr,
		})
	}

	return c.Status(fiber.StatusOK).JSON(interfaces.Response{
		Data: response,
		Status: interfaces.Status{
			Code:    fiber.StatusOK,
			Message: "Verification successful",
		},
		Error: nil,
	})
}
