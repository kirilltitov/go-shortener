package utils

import "strings"

func IsValidLink(maybeLink string) bool {
	return strings.HasPrefix(maybeLink, "https://") || strings.HasPrefix(maybeLink, "http://")
}
