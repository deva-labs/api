package controllers

import (
	"dockerwizard-api/src/modules/fiber/services"
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
		ProjectName string `json:"project_name"`
		Framework   string `json:"framework"`
		UserID      uint   `json:"user_id"`
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

	// Get socket id from user
	conn, _ := store.GetUserSocket(requestData.UserID)
	if conn == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "WebSocket connection not found for user",
		})
	}

	// Call the service function to create the fiber project
	zipPath, serviceErr := services.CreateFiberProject(conn, requestData.ProjectName, requestData.Framework)
	if serviceErr != nil {
		return c.Status(serviceErr.StatusCode).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    serviceErr.StatusCode,
				"message": serviceErr.Message,
			},
			"error": serviceErr.Err.Error(),
		})
	}

	// Prepare response data
	responseData := fiber.Map{
		"project_name": requestData.ProjectName,
		"framework":    requestData.Framework,
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
				requestData.Framework),
		},
		"data": responseData,
	})
}
