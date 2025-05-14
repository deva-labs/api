package projects

import (
	projects "dockerwizard-api/src/modules/projects/services"
	"dockerwizard-api/src/utils"
	"dockerwizard-api/store"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"time"
)

// CreateNewFiberProject is a controller function to handle create new fiber project
func CreateNewFiberProject(c *fiber.Ctx) error {
	var requestData struct {
		ProjectName string            `json:"project_name"`
		UserID      uint              `json:"user_id"`
		Env         map[string]string `json:"env"`
	}

	// Bind the incoming JSON request data to requestData struct
	if err := utils.BindJson(c, &requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Invalid request data",
			},
			"error": err.Error(),
		})
	}

	// Get socket id from users
	conn, _ := store.GetUserSocket(requestData.UserID)
	if conn == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "WebSocket connection not found for users",
		})
	}

	// Call the service function to create the fiber project
	zipPath, serviceErr := projects.CreateFiberProject(conn, requestData.ProjectName, requestData.Env)
	if serviceErr != nil {
		return c.Status(serviceErr.StatusCode).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    serviceErr.StatusCode,
				"message": serviceErr.Message,
			},
			"error": serviceErr.Err.Error(),
		})
	}

	// Response
	responseData := fiber.Map{
		"project_name": requestData.ProjectName,
		"framework":    requestData.Env["FRAMEWORK"],
		"created_at":   time.Now().Format(time.RFC3339),
		"download": fiber.Map{
			"file_name": filepath.Base(zipPath),
			"path":      zipPath,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.Map{
			"code": fiber.StatusCreated,
			"message": fmt.Sprintf("Project '%s' with framework '%s' created successfully",
				requestData.ProjectName,
				requestData.Env["FRAMEWORK"]),
		},
		"data": responseData,
	})
}
