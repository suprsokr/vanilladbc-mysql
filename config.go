package mysqlplugin

import (
	"fmt"
	"os"
	"strconv"
)

// LoadConfigFromEnv loads MySQL configuration from environment variables
// Supported variables:
//   - MYSQL_HOST (default: localhost)
//   - MYSQL_PORT (default: 3306)
//   - MYSQL_USER (default: root)
//   - MYSQL_PASSWORD (required)
//   - MYSQL_DATABASE (required)
func LoadConfigFromEnv() (Config, error) {
	config := Config{
		Host: getEnvOrDefault("MYSQL_HOST", "localhost"),
		Port: getEnvIntOrDefault("MYSQL_PORT", 3306),
		User: getEnvOrDefault("MYSQL_USER", "root"),
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		return config, fmt.Errorf("MYSQL_PASSWORD environment variable is required")
	}
	config.Password = password

	database := os.Getenv("MYSQL_DATABASE")
	if database == "" {
		return config, fmt.Errorf("MYSQL_DATABASE environment variable is required")
	}
	config.Database = database

	return config, nil
}

// NewFromEnv creates a new MySQL plugin using environment variables for configuration
func NewFromEnv(tableName string) (*Plugin, error) {
	config, err := LoadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	return New(config, tableName)
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns the environment variable as an int or a default if not set/invalid
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// MaskPassword returns a config with the password masked for logging
func (c Config) MaskPassword() Config {
	masked := c
	if c.Password != "" {
		masked.Password = "********"
	}
	return masked
}

// String returns a safe string representation of the config (password masked)
func (c Config) String() string {
	return fmt.Sprintf("mysql://%s@%s:%d/%s", c.User, c.Host, c.Port, c.Database)
}
