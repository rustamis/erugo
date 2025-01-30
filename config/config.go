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

func LoadConfig(configFile string) error {
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

	return nil
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
