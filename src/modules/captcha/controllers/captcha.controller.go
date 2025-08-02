package captcha

import (
	"deva/src/lib/dto"
	"deva/src/lib/interfaces"
	service "deva/src/modules/captcha/services"
	"deva/src/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func CreateCaptcha(c *fiber.Ctx) error {
	var body dto.CreateCaptchaRequest
	serviceErr := utils.BindJson(c, &body)
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
	response, serviceError := service.CreateCaptcha(body)
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
	return c.Status(fiber.StatusOK).JSON(interfaces.Response{
		Data: response,
		Status: interfaces.Status{
			Code:    http.StatusCreated,
			Message: "Captcha created",
		},
		Error: nil,
	})
}
