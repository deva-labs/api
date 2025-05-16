package verifications

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	verifications "skypipe/src/modules/verifications/models"
	services "skypipe/src/modules/verifications/services"
)

// VerifyCode handles the verifications code process and token generation (simple version)
func VerifyCode(c *fiber.Ctx) error {
	var code verifications.VerificationCode

	// Parse JSON input
	if err := c.BodyParser(&code); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "Invalid input",
			},
			"error": err.Error(),
		})
	}

	// Call the service to verify the code and generate tokens
	token, err := services.VerifyCode(code.Code, code.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    http.StatusUnauthorized,
				"message": err.Error(),
			},
		})
	}

	// Success response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    http.StatusOK,
			"message": "Verification successful",
		},
		"data": token,
	})
}

func VerifyCodeAndGenerateToken(c *fiber.Ctx) error {
	var code verifications.VerificationCode

	// Parse JSON input
	if err := c.BodyParser(&code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
			"error": err.Error(),
		})
	}

	// Call the service to verify the code and generate tokens
	response, statusCode, err := services.VerifyCodeAndGenerateTokens(code)
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    statusCode,
				"message": err.Error(),
			},
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Verification successful",
		},
		"data": response,
	})
}
