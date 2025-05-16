package projects

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"skypipe/src/functions"
	"skypipe/src/utils"
	"strings"
	"time"
)

// CreateFiberProject creates a new Fiber project with optional remote build configuration
func CreateFiberProject(conn *websocket.Conn, projectName string, env map[string]string) (string, *utils.ServiceError) {

	//return "", &utils.ServiceError{
	//	StatusCode: http.StatusServiceUnavailable,
	//	Message:    fmt.Sprintf("Service is unavailable. Please try again later."),
	//}
	// Validate and sanitize inputs
	framework := env["LANGUAGE"] + "-" + env["FRAMEWORK"]

	if projectName == "" || framework == "" {
		return "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "project name and framework cannot be empty",
			Err:        nil,
		}
	}
	if !isValidProjectName(projectName) {
		return "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid project name (only alphanumeric and hyphens allowed)",
			Err:        nil,
		}
	}
	finalProjectName := generateProjectName(projectName)

	// 1. Modify Makefile
	if err := updateFrameworkInMakefile("./Makefile", framework); err != nil {
		return "", err
	}

	// 2. Run installation with proper terminal handling
	if err := functions.RunProjectWorkflow(conn, finalProjectName, env); err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "project creation failed",
			Err:        err,
		}
	}
	// 3. Return a zip file path
	zipPath, err := findProjectZip(finalProjectName)
	if err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to locate project zip",
			Err:        err,
		}
	}

	return zipPath, nil
}

// Helper Functions
func isValidProjectName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(name)
}

func generateProjectName(baseName string) string {
	return fmt.Sprintf("%s-%d", strings.ToLower(baseName), time.Now().Unix())
}

func updateFrameworkInMakefile(path, framework string) *utils.ServiceError {
	data, err := os.ReadFile(path)
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to read Makefile at %s", path),
			Err:        err,
		}
	}

	modifiedData, err := replaceFrameworkVariable(string(data), framework)
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to modify Makefile for framework %s", framework),
			Err:        err,
		}
	}

	if err := os.WriteFile(path, []byte(modifiedData), 0644); err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to write Makefile at %s", path),
			Err:        err,
		}
	}

	return nil
}

func replaceFrameworkVariable(data, framework string) (string, error) {
	pattern := regexp.MustCompile(`(?m)^FRAMEWORK\s*:=\s*.*$`)
	if !pattern.MatchString(data) {
		return "", fmt.Errorf("FRAMEWORK variable not found in Makefile")
	}
	return pattern.ReplaceAllString(data, fmt.Sprintf("FRAMEWORK := %s", framework)), nil
}

func findProjectZip(projectName string) (string, error) {
	targetName := projectName + ".zip"
	publicDir := "./public"

	entries, err := os.ReadDir(publicDir)
	if err != nil {
		return "", fmt.Errorf("error reading directory %s: %v", publicDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == targetName {
			return filepath.Join(publicDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("file %s not found in %s", targetName, publicDir)
}
