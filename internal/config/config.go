package config

import (
	"os"
)

func Parse() {
	parseFlags()
}

func GetServerAddress() string {
	var result = flagBind

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	if envServerAddress != "" {
		result = envServerAddress
	}

	return result
}

func GetBaseURL() string {
	var result = flagBaseURL

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		result = envBaseURL
	}

	return result
}

func GetFileStoragePath() string {
	var result = flagFileStoragePath

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		result = envFileStoragePath
	}

	return result
}
