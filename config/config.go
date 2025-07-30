package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Days    int      `json:"days"`
	Limit   int      `json:"limit"`
	Types   []string `json:"types"`
	Verbose bool     `json:"verbose"`
	Keyword string   `json:"keyword"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Days:    30,
		Limit:   30,
		Types:   []string{"public"},
		Verbose: false,
		Keyword: "",
	}
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and set defaults for missing values
	if config.Days <= 0 {
		config.Days = 30
	}
	if config.Limit <= 0 {
		config.Limit = 30
	}
	if len(config.Types) == 0 {
		config.Types = []string{"public"}
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(path string, config *AppConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return "config/app.json"
}

// ValidateConfig validates the configuration
func ValidateConfig(config *AppConfig) error {
	if config.Days < 1 {
		return fmt.Errorf("days must be at least 1")
	}
	if config.Limit < 1 {
		return fmt.Errorf("limit must be at least 1")
	}
	if len(config.Types) == 0 {
		return fmt.Errorf("at least one channel type must be specified")
	}
	
	// Validate channel types
	for _, t := range config.Types {
		if t != "public" && t != "private" {
			return fmt.Errorf("invalid channel type: %s (must be 'public' or 'private')", t)
		}
	}
	
	return nil
} 