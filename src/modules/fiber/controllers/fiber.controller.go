package controllers

import (
	"dockerwizard-api/src/modules/fiber/services"
	"dockerwizard-api/src/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"time"
)

// CreateNewFiberProject is a controller function to handle create new fiber project
func CreateNewFiberProject(c *fiber.Ctx) error {
	var requestData struct {
		ProjectName  string `json:"project_name"`
		Framework    string `json:"framework"`
		RemoteConfig *struct {
			DockerHost    string `json:"docker_host"`
			UseDefaultTLS bool   `json:"use_default_tls"`
		} `json:"remote_config,omitempty"`
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

	var remoteConfig *services.RemoteBuildConfig
	if requestData.RemoteConfig != nil {
		// Set default TLS paths if requested
		tlsCA := ""
		tlsCert := ""
		tlsKey := ""

		if requestData.RemoteConfig.UseDefaultTLS {
			// Use the default paths where we copied the certs
			basePath := filepath.Join("store", "secrets")
			tlsCA = filepath.Join(basePath, "ca.pem")
			tlsCert = filepath.Join(basePath, "cert.pem")
			tlsKey = filepath.Join(basePath, "key.pem")

			// Verify the cert files exist
			if _, err := os.Stat(tlsCA); os.IsNotExist(err) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": fiber.Map{
						"code":    fiber.StatusBadRequest,
						"message": "Default TLS CA certificate not found",
					},
					"error": fmt.Sprintf("CA certificate not found at: %s", tlsCA),
				})
			}
			// Similar checks for other cert files...
		}

		remoteConfig = &services.RemoteBuildConfig{
			DockerHost:    requestData.RemoteConfig.DockerHost,
			TLSCACertPath: tlsCA,
			TLSCertPath:   tlsCert,
			TLSKeyPath:    tlsKey,
			ProjectName:   requestData.ProjectName,
		}
	}

	// Call the service function to create the fiber project
	zipPath, serviceErr := services.CreateFiberProject(requestData.ProjectName, requestData.Framework, remoteConfig)
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

	if remoteConfig != nil {
		remoteInfo := fiber.Map{
			"docker_host": remoteConfig.DockerHost,
			"status":      "configured",
		}
		if requestData.RemoteConfig.UseDefaultTLS {
			remoteInfo["tls_config"] = "using_default_certificates"
		}
		responseData["remote_build"] = remoteInfo
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
