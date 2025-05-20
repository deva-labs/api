package projects

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"skypipe/src/lib/dto"
	"skypipe/src/lib/interfaces"
	projects "skypipe/src/modules/projects/services"
	"skypipe/src/utils"
	"skypipe/store"
	"time"
)

// CreateNewFiberProject is a controller function to handle create new fiber project
func CreateNewFiberProject(c *fiber.Ctx) error {
	var requestData dto.CreateFiberRequest

	// Bind the incoming JSON request data to requestData struct
	serviceErr := utils.BindJson(c, &requestData)
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

	// Get socket id from users
	conn, _ := store.GetUserSocket(requestData.UserID)
	if conn == nil {
		return c.Status(fiber.StatusNotFound).JSON(interfaces.Response{
			Data: nil,
			Status: interfaces.Status{
				Code:    fiber.StatusNotFound,
				Message: "WebSocket connection not found for users",
			},
			Error: nil,
		})
	}

	// Call the services function to create the fiber project
	zipPath, serviceError := projects.CreateFiberProject(conn, requestData.ProjectName, requestData.Env)
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

	return c.Status(fiber.StatusCreated).JSON(interfaces.Response{
		Data: responseData,
		Status: interfaces.Status{
			Code: fiber.StatusOK,
			Message: fmt.Sprintf("Project '%s' with framework '%s' created successfully",
				requestData.ProjectName,
				requestData.Env["FRAMEWORK"]),
		},
		Error: nil,
	})
}
