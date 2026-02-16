package mysqlplugin

import (
	"encoding/json"
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

// LoadConfigFromFile loads MySQL configuration from a JSON file
func LoadConfigFromFile(filename string) (Config, error) {
	var config Config

	file, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.Password == "" {
		return config, fmt.Errorf("password is required in config file")
	}
	if config.Database == "" {
		return config, fmt.Errorf("database is required in config file")
	}

	// Set defaults
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.User == "" {
		config.User = "root"
	}

	return config, nil
}

// NewFromFile creates a new MySQL plugin using a JSON config file
func NewFromFile(filename, tableName string) (*Plugin, error) {
	config, err := LoadConfigFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	return New(config, tableName)
}

// SaveConfigToFile saves the configuration to a JSON file
func (c Config) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filename, data, 0600) // 0600 = read/write for owner only
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
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
