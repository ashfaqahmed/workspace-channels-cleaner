package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvironment loads environment variables from .env file
func LoadEnvironment() error {
	return godotenv.Load()
}

// GetWorkspaceToken returns the workspace API token from environment
func GetWorkspaceToken() string {
	return os.Getenv("WORKSPACE_API_TOKEN")
}

// ValidateToken checks if the workspace token is set
func ValidateToken() error {
	token := GetWorkspaceToken()
	if token == "" {
		return &ConfigError{Message: "WORKSPACE_API_TOKEN not set in environment or .env file"}
	}
	return nil
}

// ConfigError represents configuration errors
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
