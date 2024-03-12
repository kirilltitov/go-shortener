package config

import (
	"flag"
	"fmt"
)

const defaultPort = 8080

var Bind string
var BaseURL string

func ParseFlags() {
	flag.StringVar(&Bind, "a", fmt.Sprintf(":%d", defaultPort), "Host and port to bind")
	flag.StringVar(&BaseURL, "b", fmt.Sprintf("http://localhost:%d", defaultPort), "Base URL")

	flag.Parse()
}
