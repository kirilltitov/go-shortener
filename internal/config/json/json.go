package json

import (
	"encoding/json"
	"os"
)

// Config хранит в себе конфиг, загружаемый из JSON-файла
type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
	GrpcAddress     string `json:"grpc_address"`
}

// Load загружает JSON-файл с конфигурацией и, в случае успеха, возвращает структуру с конфигурацией
func Load(data []byte) (*Config, error) {
	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Read открывает файл с JSON-конфигом и возвращает его содержимое в виде байтов
func Read(filename string) ([]byte, error) {
	result, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return result, nil
}
