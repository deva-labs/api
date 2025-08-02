package users

import (
	"deva/src/functions"
	"deva/src/lib/dto"
	"deva/src/lib/interfaces"
	users "deva/src/modules/users/models"
	service "deva/src/modules/users/services"
	"deva/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

func GetUser(c *fiber.Ctx) error {
	userInterface := c.Locals("user")
	if userInterface == nil {
		return c.Status(http.StatusUnauthorized).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			},
			Error: nil,
		})
	}

	currentUser, ok := userInterface.(*users.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse user from context",
			},
			Error: nil,
		})
	}

	userInfo, serviceErr := service.GetUserInfo(currentUser.ID)
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

	return c.Status(http.StatusOK).JSON(interfaces.Response{
		Data: functions.ToUserResponse(userInfo),
		Status: interfaces.Status{
			Code:    http.StatusOK,
			Message: "Retrieved the profile of user successfully",
		},
		Error: nil,
	})
}

func Register(c *fiber.Ctx) error {
	var body dto.RegisterRequest
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

	response, serviceError := service.RegisterService(body)
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

	return c.Status(fiber.StatusCreated).JSON(interfaces.Response{
		Data: response,
		Status: interfaces.Status{
			Code:    http.StatusCreated,
			Message: "Registration successfully, redirecting to login...",
		},
		Error: nil,
	})
}

func Login(c *fiber.Ctx) error {
	var body dto.LoginRequest
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
	response, serviceError := service.LoginService(body)
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
			Code:    http.StatusOK,
			Message: "We just send to you an email verification code. Redirecting...",
		},
		Error: nil,
	})
}

func ForgotPassword(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}
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

	response, serviceError := service.ForgotPassword(request.Email)
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
			Code:    fiber.StatusOK,
			Message: "Verification email has been sent",
		},
		Error: nil,
	})
}

func RenewPassword(c *fiber.Ctx) error {
	var request dto.RenewPasswordRequest

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

	serviceError := service.RenewPassword(request)
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
		Data: nil,
		Status: interfaces.Status{
			Code:    fiber.StatusOK,
			Message: "Password has been renewed",
		},
		Error: nil,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var request struct {
		OldPassword string    `json:"old_password"`
		NewPassword string    `json:"new_password"`
		UserID      uuid.UUID `json:"user_id"`
	}

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

	serviceErr = service.ChangePassword(request.OldPassword, request.NewPassword, request.UserID)
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
		Data: nil,
		Status: interfaces.Status{
			Code:    fiber.StatusOK,
			Message: "Password has been renewed",
		},
		Error: nil,
	})
}

func GetUserAvatar(c *fiber.Ctx) error {
	var request struct {
		UserID uint `json:"user_id"`
	}

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

	avatarUser, serviceError := service.GetUserImageByID(request.UserID)
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
		Data: avatarUser,
		Status: interfaces.Status{
			Code:    fiber.StatusOK,
			Message: "Successfully retrieved user avatar",
		},
		Error: nil,
	})
}
