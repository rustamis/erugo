package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/DeanWard/erugo/models"
)

var AppConfig models.Config

// Environment variable names
const (
	ENV_BASE_STORAGE_PATH  = "ERUGO_BASE_STORAGE_PATH"
	ENV_APP_URL            = "ERUGO_APP_URL"
	ENV_BIND_PORT          = "ERUGO_BIND_PORT"
	ENV_JWT_SECRET         = "ERUGO_JWT_SECRET"
	ENV_MAX_SHARE_SIZE     = "ERUGO_MAX_SHARE_SIZE"
	ENV_DATABASE_FILE_PATH = "ERUGO_DATABASE_FILE_PATH"
	ENV_PRIVATE_DATA_PATH  = "ERUGO_PRIVATE_DATA_PATH"
)

func LoadConfig(configFile string) error {
	// First load from config file
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}

	// Then override with environment variables if they exist
	if envVal := os.Getenv(ENV_BASE_STORAGE_PATH); envVal != "" {
		AppConfig.BaseStoragePath = envVal
	}

	if envVal := os.Getenv(ENV_APP_URL); envVal != "" {
		AppConfig.AppUrl = envVal
	}

	if envVal := os.Getenv(ENV_BIND_PORT); envVal != "" {
		if port, err := strconv.Atoi(envVal); err == nil {
			AppConfig.BindPort = port
		}
	}

	if envVal := os.Getenv(ENV_JWT_SECRET); envVal != "" {
		AppConfig.JwtSecret = envVal
	}

	if envVal := os.Getenv(ENV_MAX_SHARE_SIZE); envVal != "" {
		AppConfig.MaxShareSize = envVal
	}

	if envVal := os.Getenv(ENV_DATABASE_FILE_PATH); envVal != "" {
		AppConfig.DatabaseFilePath = envVal
	}

	if envVal := os.Getenv(ENV_PRIVATE_DATA_PATH); envVal != "" {
		AppConfig.PrivateDataPath = envVal
	}

	return nil
}

// GetEnvWithDefault gets an environment variable value or returns the default if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetMaxShareSize() int64 {
	rawSize := AppConfig.MaxShareSize
	convertedSize, err := ConvertSizeToBytes(rawSize)
	if err != nil {
		return 0
	}
	return int64(convertedSize)
}

// ConvertSizeToBytes converts a human-readable file size string to bytes
func ConvertSizeToBytes(sizeStr string) (uint64, error) {
	// Regular expression to match size patterns like "2G", "250M", "512K"
	re := regexp.MustCompile(`(?i)^(\d+)([KMGT]?B?)$`)
	matches := re.FindStringSubmatch(strings.ToUpper(sizeStr))

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	value, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %s", matches[1])
	}

	var multiplier uint64 = 1
	switch matches[2] {
	case "K", "KB":
		multiplier = 1024
	case "M", "MB":
		multiplier = 1024 * 1024
	case "G", "GB":
		multiplier = 1024 * 1024 * 1024
	case "T", "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	return value * multiplier, nil
}
