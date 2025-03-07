package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Config holds all configuration for the application
type Config struct {
	// MongoDB connection
	MongoDBURI      string
	MongoDBDatabase string

	// LINE API
	LineChannelSecret string
	LineChannelToken  string

	// UR API
	URAPIKey string

	// Monitoring settings
	MonitorInterval int // in minutes
	DanchiList      []string

	// Server settings
	Port int
}

var (
	config Config
	once   sync.Once
	err    error
)

// GetConfig returns the application configuration, loading it if necessary
func GetConfig() (Config, error) {
	once.Do(func() {
		config = Config{
			MongoDBURI:       getEnv("MONGODB_URI", ""),
			MongoDBDatabase:  getEnv("MONGODB_DATABASE", "ur_monitor"),
			LineChannelToken: getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),
			URAPIKey:         getEnv("UR_API_KEY", ""),
			MonitorInterval:  getEnvAsInt("MONITOR_INTERVAL", 60),
			DanchiList:       getEnvAsStringSlice("DANCHI_LIST", []string{}),
			Port:             getEnvAsInt("PORT", 8080),
		}

		// Validate required configuration
		if config.MongoDBURI == "" {
			err = fmt.Errorf("MONGODB_URI environment variable is required")
			return
		}

		if config.LineChannelSecret == "" || config.LineChannelToken == "" {
			err = fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN environment variables are required")
			return
		}

		if config.URAPIKey == "" {
			err = fmt.Errorf("UR_API_KEY environment variable is required")
			return
		}

		if len(config.DanchiList) == 0 {
			err = fmt.Errorf("DANCHI_LIST environment variable is required")
			return
		}
	})

	return config, err
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as int with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as string slice with a default value
func getEnvAsStringSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	return strings.Split(valueStr, ",")
}
