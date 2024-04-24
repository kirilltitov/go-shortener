package config

import (
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

func New() Config {
	parseFlags()

	return Config{
		ServerAddress:   getServerAddress(),
		BaseURL:         getBaseURL(),
		FileStoragePath: getFileStoragePath(),
		DatabaseDSN:     getDatabaseDSN(),
	}
}

func getServerAddress() string {
	var result = flagBind

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	if envServerAddress != "" {
		result = envServerAddress
	}

	return result
}

func getBaseURL() string {
	var result = flagBaseURL

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		result = envBaseURL
	}

	return result
}

func getFileStoragePath() string {
	var result = flagFileStoragePath

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		result = envFileStoragePath
	}

	return result
}

func getDatabaseDSN() string {
	var result = flagDatabaseDSN

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		result = envDatabaseDSN
	}

	return result
}
