package utils

import (
	"net/url"
)

var allowedProtocols = map[string]bool{
	"http":  true,
	"https": true,
}

// IsValidURL проверяет переданный урл на предмет корректности.
func IsValidURL(maybeLink string) bool {
	parsedURL, err := url.Parse(maybeLink)
	_, protocolFound := allowedProtocols[parsedURL.Scheme]
	return err == nil && protocolFound
}
