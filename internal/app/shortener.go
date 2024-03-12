package internal

import (
	"fmt"
	"github.com/jxskiss/base62"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var cur int = 0
var storage internalStorage.Storage = internalStorage.InMemory{}

func isValidLink(maybeLink string) bool {
	return strings.HasPrefix(maybeLink, "https://") || strings.HasPrefix(maybeLink, "http://")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var result string
	var code int

	w.Header().Set("Content-Type", "text/plain")

	defer func() {
		w.WriteHeader(code)
		io.WriteString(w, result)
	}()

	switch r.Method {
	case http.MethodGet:
		shortURL := strings.TrimPrefix(r.URL.Path, "/")
		decodedStringInt, err := base62.DecodeString(shortURL)
		if err != nil {
			result = fmt.Sprintf("Could not decode short url '%s'\n", shortURL)
			code = http.StatusBadRequest
			break
		}
		decodedInt, err := strconv.Atoi(string(decodedStringInt))
		if err != nil {
			result = fmt.Sprintf("Could not decode short url '%s'\n", shortURL)
			code = http.StatusBadRequest
			break
		}
		url, ok := storage.Get(decodedInt)
		if !ok {
			result = fmt.Sprintf("URL '%s' not found\n", shortURL)
			code = http.StatusBadRequest
			break
		}

		code = http.StatusTemporaryRedirect

		http.Redirect(w, r, url, code)
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			result = fmt.Sprintf("Could not read body: %s\n", err.Error())
			code = http.StatusBadRequest
			break
		}
		link := string(b)
		if !isValidLink(link) {
			result = fmt.Sprintf("Invalid link (must start with https://): %s\n", link)
			code = http.StatusBadRequest
			break
		}

		cur++
		storage.Set(cur, link)

		result = fmt.Sprintf("http://localhost:8080/%s", base62.EncodeToString([]byte(strconv.Itoa(cur))))
		code = http.StatusCreated
	}
}
