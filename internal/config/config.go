package config

import (
	"os"
)

// Config является объектом для хранения конфигурации сервиса.
type Config struct {
	// ServerAddress является адресом (включая порт), на котором поднимется веб-сервер.
	ServerAddress string

	// BaseURL является полным адресом (протокол, домен, порт, путь), который будет использоваться при генерации коротких ссылок.
	BaseURL string

	// FileStoragePath является путем до файлового хранилища сервиса.
	// Пустой путь означает хранилище в памяти.
	// Если задано значение DatabaseDSN, это поле игнорируется.
	FileStoragePath string

	// DatabaseDSN является конфигурационной DSN-строкой для подключения к PostgreSQL.
	// Если поле не выставлено, используется файловое хранилище,
	// либо же (если у файлового хранилища не выставлен путь), хранилище в памяти.
	DatabaseDSN string
}

// New создает, автоматически заполняет и возвращает экземпляр конфигурации сервиса.
func New() Config {
	parseFlags()

	return NewWithoutParsing()
}

// NewWithoutParsing создает и возвращает экземпляр конфигурации сервиса без парсинга флагов.
func NewWithoutParsing() Config {
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
