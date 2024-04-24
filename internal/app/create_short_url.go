package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) HandlerCreateShortURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not read body: %s\n", err.Error()))
		return
	}

	code := http.StatusCreated
	URL := string(b)
	shortURL, err := a.Shortener.ShortenURL(r.Context(), URL)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			code = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
			return
		}
	}

	w.WriteHeader(code)
	io.WriteString(w, shortURL)
}