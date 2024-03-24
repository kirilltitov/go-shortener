package utils

import (
	"net/url"
)

var allowedProtocols = map[string]bool{
	"http":  true,
	"https": true,
}

func IsValidLink(maybeLink string) bool {
	parsedURL, err := url.Parse(maybeLink)
	_, protocolFound := allowedProtocols[parsedURL.Scheme]
	return err == nil && protocolFound
}
