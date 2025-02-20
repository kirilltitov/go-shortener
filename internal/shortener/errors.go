package shortener

import "errors"

// ErrorUnauthorized является ошибкой доступа к ресурсу
var ErrorUnauthorized = errors.New("unauthorized")
