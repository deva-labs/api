package verifications

import (
	"deva/src/lib/dto"
	"deva/src/lib/interfaces"
	services "deva/src/modules/verifications/services"
	"deva/src/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// VerifyCodeAndSetPasswordToken handles the verifications code process and token generation (simple version)
func VerifyCodeAndSetPasswordToken(c *fiber.Ctx) error {
	var request dto.VerificationCodeRequest

	// Parse JSON input
	serviceErr := utils.BindJson(c, &request)
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

	// Call the services to verify the code and generate tokens
	token, serviceError := services.VerifyCodeAndSetPasswordToken(request)
	if serviceError != nil {
		s := serviceError.Err.Error()
		errStr := &s
		return c.Status(serviceError.StatusCode).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    serviceError.StatusCode,
				Message: serviceError.Message,
			},
			Error: errStr,
		})
	}

	// Success response
	return c.Status(http.StatusOK).JSON(interfaces.Response{
		Data: token,
		Status: interfaces.Status{
			Code:    http.StatusOK,
			Message: "Verification successfully.",
		},
		Error: nil,
	})
}

func VerifyCodeAndGenerateToken(c *fiber.Ctx) error {
	var code dto.VerificationCodeRequest
	// Parse JSON input
	serviceErr := utils.BindJson(c, &code)
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

	// Call the services to verify the code and generate tokens
	response, serviceError := services.VerifyCodeAndGenerateTokens(code)
	if serviceError != nil {
		s := serviceError.Err.Error()
		errStr := &s
		return c.Status(serviceError.StatusCode).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    serviceError.StatusCode,
				Message: serviceError.Message,
			},
			Error: errStr,
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(interfaces.Response{
		Data: response,
		Status: interfaces.Status{
			Code:    http.StatusOK,
			Message: "Login successful! Redirecting to home.",
		},
		Error: nil,
	})
}
