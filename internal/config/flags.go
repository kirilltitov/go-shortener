package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8080

var flagBind string
var flagBaseURL string

func parseFlags() {
	flag.StringVar(&flagBind, "a", fmt.Sprintf(":%d", defaultPort), "Host and port to bind")
	flag.StringVar(&flagBaseURL, "b", fmt.Sprintf("http://localhost:%d", defaultPort), "Base URL")

	flag.Parse()
}
