package middlewares

import (
	"deva/src/lib/dto"
	"deva/src/lib/interfaces"
	key_token "deva/src/modules/key_token/services"
	roles "deva/src/modules/roles/services"
	users "deva/src/modules/users/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token missing",
			})
		}

		token, err := SplitToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		userInfo, serviceErr := key_token.VerifyToken(token)
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

		c.Locals("user", userInfo)
		return c.Next()
	}
}

func Authorization(requiredPermission dto.PermissionInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInterface := c.Locals("user")
		if userInterface == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		user, ok := userInterface.(*users.User)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "User type assertion failed",
			})
		}

		// Get user role from user_id
		role, err := roles.GetRoleByUserID(user.ID)
		if err != nil {
			return err
		}
		if !roles.CheckPermission(role.ID, requiredPermission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: insufficient permissions",
			})
		}

		return c.Next()
	}
}

// SplitToken is the function to extract bearer token from Authorization
func SplitToken(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("token format must be 'Bearer <token>'")
	}
	return parts[1], nil
}
