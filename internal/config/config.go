package config

import (
	"net"
	"os"

	"github.com/kirilltitov/go-shortener/internal/config/json"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

// Config является объектом для хранения конфигурации сервиса.
type Config struct {
	// ServerAddress является адресом (включая порт), на котором поднимется веб-сервер.
	ServerAddress string

	// GrpcAddress является адресом (включая порт), на котором поднимется gRPC-сервер
	GrpcAddress string

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

	// EnableHTTPS заставляет сервер запускаться в режиме HTTPS
	EnableHTTPS string

	// TrustedSubnet является доверенной подсетью, запросы из которой разрешены для внутренних методов
	TrustedSubnet *net.IPNet
}

// New создает, автоматически заполняет и возвращает экземпляр конфигурации сервиса.
func New() Config {
	parseFlags()

	return NewWithoutParsing()
}

// NewWithoutParsing создает и возвращает экземпляр конфигурации сервиса без парсинга флагов.
func NewWithoutParsing() Config {
	var jsonConfig json.Config
	if jsonFilePath := getJSONFilePath(); jsonFilePath != "" {
		data, err := json.Read(jsonFilePath)
		if err != nil {
			logger.Log.Infof("Could not load JSON config file: %s, ignoring", err)
		} else {
			c, err2 := json.Load(data)
			if err2 != nil {
				logger.Log.Infof("Could not load JSON config file: %s, ignoring", err2)
			} else {
				jsonConfig = *c
				logger.Log.Infof("Loaded config file %s with config %+v", jsonFilePath, jsonConfig)
			}
		}
	}

	return Config{
		ServerAddress:   getServerAddress(jsonConfig),
		GrpcAddress:     getGrpcServerAddress(jsonConfig),
		BaseURL:         getBaseURL(jsonConfig),
		FileStoragePath: getFileStoragePath(jsonConfig),
		DatabaseDSN:     getDatabaseDSN(jsonConfig),
		EnableHTTPS:     getEnableHTTPS(jsonConfig),
		TrustedSubnet:   getTrustedSubnet(jsonConfig),
	}
}

func getServerAddress(jsonConfig json.Config) string {
	var result = flagBind

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	if envServerAddress != "" {
		result = envServerAddress
	}

	if result == "" && jsonConfig.ServerAddress != "" {
		result = jsonConfig.ServerAddress
	}

	return result
}

func getGrpcServerAddress(jsonConfig json.Config) string {
	var result = flagGrpcBind

	envGrpcServerAddress := os.Getenv("GRPC_SERVER")
	if envGrpcServerAddress != "" {
		result = envGrpcServerAddress
	}

	if result == "" && jsonConfig.GrpcAddress != "" {
		result = jsonConfig.GrpcAddress
	}

	return result
}

func getBaseURL(jsonConfig json.Config) string {
	var result = flagBaseURL

	envBaseURL := os.Getenv("BASE_URL")
	if envBaseURL != "" {
		result = envBaseURL
	}

	if result == "" && jsonConfig.BaseURL != "" {
		result = jsonConfig.BaseURL
	}

	return result
}

func getFileStoragePath(jsonConfig json.Config) string {
	var result = flagFileStoragePath

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		result = envFileStoragePath
	}

	if result == "" && jsonConfig.FileStoragePath != "" {
		result = jsonConfig.FileStoragePath
	}

	return result
}

func getDatabaseDSN(jsonConfig json.Config) string {
	var result = flagDatabaseDSN

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		result = envDatabaseDSN
	}

	if result == "" && jsonConfig.DatabaseDSN != "" {
		result = jsonConfig.DatabaseDSN
	}

	return result
}

func getEnableHTTPS(jsonConfig json.Config) string {
	var result = flagEnableHTTPS

	envEnableHTTPS := os.Getenv("ENABLE_HTTPS")
	if envEnableHTTPS != "" {
		result = envEnableHTTPS
	}

	if result == "" && jsonConfig.EnableHTTPS {
		result = "true"
	}

	return result
}

func getTrustedSubnet(jsonConfig json.Config) *net.IPNet {
	var result = flagTrustedSubnet

	envTrustedSubnet := os.Getenv("TRUSTED_SUBNET")
	if envTrustedSubnet != "" {
		result = envTrustedSubnet
	}

	if result == "" && jsonConfig.TrustedSubnet != "" {
		result = jsonConfig.TrustedSubnet
	}

	if result == "" {
		return nil
	}

	_, ipNet, err := net.ParseCIDR(result)
	if err != nil {
		logger.Log.WithError(err).Error("Could not parse trusted subnet")
		return nil
	}

	return ipNet
}

func getJSONFilePath() string {
	var result = flagJSONConfig

	JSONConfigPath := os.Getenv("CONFIG")
	if JSONConfigPath != "" {
		result = JSONConfigPath
	}

	return result
}
