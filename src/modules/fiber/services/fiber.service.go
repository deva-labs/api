package services

import (
	"dockerwizard-api/src/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/creack/pty"
)

type RemoteBuildConfig struct {
	DockerHost    string `json:"dockerHost"`
	TLSCACertPath string `json:"tlsCACertPath"`
	TLSCertPath   string `json:"tlsCertPath"`
	TLSKeyPath    string `json:"tlsKeyPath"`
	ProjectName   string `json:"projectName"`
}

// CreateFiberProject creates a new Fiber project with optional remote build configuration
func CreateFiberProject(projectName, framework string, remoteConfig *RemoteBuildConfig) (string, *utils.ServiceError) {
	TestProgressBar60Seconds()
	return "", &utils.ServiceError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    fmt.Sprintf("Service is unavailable. Please try again later."),
	}
	// Validate and sanitize inputs
	projectName = strings.TrimSpace(projectName)
	framework = strings.TrimSpace(framework)

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
	if err := runInstallation(finalProjectName); err != nil {
		return "", err
	}

	// 3. Optionally run remote build
	if remoteConfig != nil {
		remoteConfig.ProjectName = finalProjectName
		if err := runRemoteBuild(*remoteConfig); err != nil {
			return "", err
		}
	}

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

func runInstallation(projectName string) *utils.ServiceError {
	// Prepare command based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "install.bat", projectName)
	} else {
		cmd = exec.Command("bash", "install.sh", projectName)
	}

	// Check terminal capabilities
	if term.IsTerminal(uintptr(int(os.Stdout.Fd()))) {
		return runWithPTY(cmd)
	}
	return runWithoutPTY(cmd)
}

func runWithPTY(cmd *exec.Cmd) *utils.ServiceError {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start PTY",
			Err:        err,
		}
	}
	defer ptmx.Close()

	// Handle terminal state
	oldState, err := term.MakeRaw(uintptr(int(os.Stdin.Fd())))
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to set terminal raw mode",
			Err:        err,
		}
	}
	defer term.Restore(uintptr(int(os.Stdin.Fd())), oldState)

	// Handle terminal resizing
	ch := make(chan os.Signal, 1)
	defer close(ch)

	// Copy streams
	go func() { io.Copy(os.Stdout, ptmx) }()
	go func() { io.Copy(ptmx, os.Stdin) }()

	if err := cmd.Wait(); err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "installation failed",
			Err:        err,
		}
	}

	return nil
}

func runWithoutPTY(cmd *exec.Cmd) *utils.ServiceError {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "installation failed",
			Err:        err,
		}
	}

	return nil
}

func runRemoteBuild(config RemoteBuildConfig) *utils.ServiceError {
	if err := validateCertPaths(config); err != nil {
		return err
	}

	// Prepare command based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "build.bat", config.ProjectName, config.DockerHost)
	} else {
		cmd = exec.Command("bash", "build.sh", config.ProjectName, config.DockerHost)
	}

	// Run with appropriate terminal handling
	if term.IsTerminal(uintptr(int(os.Stdout.Fd()))) {
		return runWithPTY(cmd)
	}
	return runWithoutPTY(cmd)
}

func validateCertPaths(config RemoteBuildConfig) *utils.ServiceError {
	requiredFiles := []struct {
		desc string
		path string
	}{
		{"CA certificate", config.TLSCACertPath},
		{"Client certificate", config.TLSCertPath},
		{"Client key", config.TLSKeyPath},
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file.path); os.IsNotExist(err) {
			return &utils.ServiceError{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("%s not found at %s", file.desc, file.path),
				Err:        err,
			}
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

func TestProgressBar60Seconds() {
	totalDuration := 60 * time.Second
	start := time.Now()
	barWidth := 50

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Tắt dòng buffering nếu bị ảnh hưởng bởi log framework
	stdout := os.Stdout

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start)
			if elapsed >= totalDuration {
				goto Done
			}

			percent := float64(elapsed) / float64(totalDuration)
			filled := int(percent * float64(barWidth))
			empty := barWidth - filled

			bar := fmt.Sprintf("\r[%s%s] %3.0f%%",
				strings.Repeat("=", filled),
				strings.Repeat(" ", empty),
				percent*100)

			// Ghi ra stderr để tránh bị Fiber/Air che stdout
			_, err := fmt.Fprint(stdout, bar)
			if err != nil {
				return
			}
		}
	}
Done:
	fmt.Fprintf(stdout, "\r[%s] 100%%\n", strings.Repeat("=", barWidth))
	_, err := fmt.Fprintln(stdout, "✔ Test complete in 60 seconds.")
	
	if err != nil {
		return
	}
}
