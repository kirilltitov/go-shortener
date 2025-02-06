package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8080

var flagBind = fmt.Sprintf(":%d", defaultPort)
var flagBaseURL = fmt.Sprintf("http://localhost:%d", defaultPort)
var flagFileStoragePath = "/tmp/short-url-db.json"
var flagDatabaseDSN = ""
var flagEnableHTTPS = ""
var flagJSONConfig = ""

func parseFlags() {
	flag.StringVar(&flagBind, "a", flagBind, "Host and port to bind")
	flag.StringVar(&flagBaseURL, "b", flagBaseURL, "Base URL")
	flag.StringVar(&flagFileStoragePath, "f", flagFileStoragePath, "File storage path")
	flag.StringVar(&flagDatabaseDSN, "d", flagDatabaseDSN, "Database DSN")
	flag.StringVar(&flagEnableHTTPS, "s", flagEnableHTTPS, "Enable HTTPS")
	flag.StringVar(&flagJSONConfig, "c", flagJSONConfig, "Path to JSON config file")

	flag.Parse()
}
