package functions

import (
	"fmt"
	"net/http"
	"os"
	"skypipe/src/lib/interfaces"
	"skypipe/src/utils"
)

func validateCertPaths(config interfaces.RemoteBuildConfig) *utils.ServiceError {
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
