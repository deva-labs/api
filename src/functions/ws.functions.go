package functions

import (
	"dockerwizard-api/src/lib/interfaces"
	"dockerwizard-api/src/utils"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"strings"

	"time"
)

func RunProjectWorkflow(c *websocket.Conn, projectName string, env map[string]string) error {
	sc := &interfaces.SafeConn{Conn: c}

	baseEnv := map[string]string{
		"PROJECT_NAME":   projectName,
		"BASE_DIR":       "/app",
		"TLSCACERT_PATH": "/app/store/secrets/ca.pem",
		"TLSCERT_PATH":   "/app/store/secrets/cert.pem",
		"TLSKEY_PATH":    "/app/store/secrets/key.pem",
		"CONTEXT_NAME":   "myremote",
	}

	for k, v := range env {
		baseEnv[k] = v
	}

	steps := []interfaces.WorkflowStep{
		{Name: "environment file", Command: "make create-env", Action: "creating", EnvVars: baseEnv},
		{Name: "Dockerfile", Command: "make create-dockerfile", Action: "creating", EnvVars: baseEnv},
		{Name: "docker-compose.yml", Command: "make create-docker-compose", Action: "creating", EnvVars: baseEnv},
		{Name: "entrypoint.sh", Command: "make create-entrypoint", Action: "creating", EnvVars: baseEnv},
		{Name: "docker-bake.hcl", Command: "make create-docker-bake", Action: "creating", EnvVars: baseEnv},
		{Name: "main.go", Command: "make create-main", Action: "creating", EnvVars: baseEnv},
		{Name: "Golang", Command: "make install-go", Action: "installing", EnvVars: baseEnv},
		{Name: "Go modules", Command: "make init-go-modules", Action: "initializing", EnvVars: baseEnv},
		{Name: "Air live reload", Command: "make install-air", Action: "installing", EnvVars: baseEnv},
		{Name: "Air configuration", Command: "make air-init", Action: "initializing", EnvVars: baseEnv},
		{Name: "Docker Compose", Command: "make docker-compose-up", Action: "starting", EnvVars: baseEnv},
		{Name: "Project", Command: "make clean", Action: "exporting", EnvVars: baseEnv},
	}

	totalSteps := len(steps)
	startTime := time.Now()

	sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("🚀 Starting project workflow for %s...", projectName)))
	sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("📋 Total steps: %d\n", totalSteps)))

	for i, step := range steps {
		stepNumber := i + 1
		sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("[Step %d/%d] %s %s...", stepNumber, totalSteps, strings.Title(step.Action), step.Name)))

		if err := utils.ExecWithAnimation(sc, step.Name, step.Command, step.Action, step.EnvVars); err != nil {
			sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("❌ Error during %s: %v", step.Name, err)))
			return fmt.Errorf("step %d (%s) failed: %w", stepNumber, step.Name, err)
		}

		if stepNumber < totalSteps {
			remaining := time.Duration(float64(totalSteps-stepNumber)*time.Since(startTime).Seconds()/float64(stepNumber)) * time.Second
			sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("⏱️  Estimated time remaining: %ds", int(remaining.Seconds()))))
		}
	}

	elapsed := time.Since(startTime).Round(time.Second)
	sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("\n🎉 Workflow completed successfully in %s!", elapsed)))
	sc.SafeWrite(websocket.TextMessage, []byte("✅ Your project is ready to go!"))
	return nil
}
