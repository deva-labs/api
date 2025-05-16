package users

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"skypipe/src/lib/dto"
	users "skypipe/src/modules/users/models"
	service "skypipe/src/modules/users/services"
	"skypipe/src/utils"
	"strconv"
)

func GetUser(c *fiber.Ctx) error {
	userInterface := c.Locals("user")
	if userInterface == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			},
		})
	}

	currentUser, ok := userInterface.(*users.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "Failed to parse user from context",
			},
		})
	}

	userInfo, err := service.GetUserInfo(currentUser.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "Failed to get user info",
			},
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    http.StatusOK,
			"message": "Retrieved the profile of user successfully",
		},
		"data": userInfo,
	})
}

func Register(c *fiber.Ctx) error {
	var body dto.RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
		})
	}
	if !utils.VerifyCaptcha(body.Captcha) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusForbidden,
				"message": "Captcha verifications failed",
			},
		})
	}

	response, serviceError := service.RegisterService(body)
	if serviceError != nil {
		return c.Status(serviceError.StatusCode).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    serviceError.StatusCode,
				"message": serviceError.Message,
			},
			"error": serviceError.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusCreated,
			"message": "User created successfully",
		},
		"data": response,
	})
}

func Login(c *fiber.Ctx) error {
	var body dto.LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
		})
	}
	response, serviceError := service.LoginService(body)
	if serviceError != nil {
		return c.Status(serviceError.StatusCode).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    serviceError.StatusCode,
				"message": serviceError.Message,
			},
			"error": serviceError.Error,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "We just send an email verification code to you",
		},
		"data": response,
	})
}

func ForgotPassword(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
			"error": err.Error(),
		})
	}

	status, err := service.ForgotPassword(request.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    status,
				"message": "Forbidden",
			},
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Verification email has been sent",
		},
	})
}

func RenewPassword(c *fiber.Ctx) error {
	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		UserID      uint   `json:"user_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
		})
	}

	status, err := service.RenewPassword(request.NewPassword, strconv.Itoa(int(request.UserID)))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    status,
				"message": "Failed to change password",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Successfully changed password",
		},
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		UserID      uint   `json:"user_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
		})
	}

	status, err := service.ChangePassword(request.OldPassword, request.NewPassword, strconv.Itoa(int(request.UserID)))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    status,
				"message": "Failed to change password",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Successfully changed password",
		},
	})
}

func GetUserAvatar(c *fiber.Ctx) error {
	var request struct {
		UserID uint `json:"user_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid input",
			},
			"error": err.Error(),
		})
	}

	avatarUser, err := service.GetUserImageByID(request.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": err.Message,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    fiber.StatusOK,
			"message": "Successfully retrieved user avatar",
		},
		"data": avatarUser,
	})
}
