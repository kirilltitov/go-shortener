package storage

import (
	"fmt"
	"strconv"

	"github.com/jxskiss/base62"
)

func intToShortURL(i int) string {
	return base62.EncodeToString([]byte(strconv.Itoa(i)))
}

func shortURLToInt(s string) (int, error) {
	decodedStringInt, err := base62.DecodeString(s)
	if err != nil {
		return 0, fmt.Errorf("could not decode short url '%s'", s)
	}

	decodedInt, err := strconv.Atoi(string(decodedStringInt))
	if err != nil {
		return 0, fmt.Errorf("could not decode short url '%s'", s)
	}

	return decodedInt, nil
}
