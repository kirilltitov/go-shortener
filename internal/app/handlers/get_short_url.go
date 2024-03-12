package handlers

import (
	"fmt"
	"github.com/jxskiss/base62"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func HandlerGetShortURL(w http.ResponseWriter, r *http.Request, storage internalStorage.Storage) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	decodedStringInt, err := base62.DecodeString(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not decode short url '%s'\n", shortURL))
		return
	}

	decodedInt, err := strconv.Atoi(string(decodedStringInt))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not decode short url '%s'\n", shortURL))
		return
	}

	url, ok := storage.Get(decodedInt)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("URL '%s' not found\n", shortURL))
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}