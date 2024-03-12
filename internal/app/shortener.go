package internal

import (
	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"net/http"
)

var cur int = 0
var storage internalStorage.Storage = internalStorage.InMemory{}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	switch r.Method {
	case http.MethodGet:
		handlers.HandlerGetShortURL(w, r, storage)
	case http.MethodPost:
		handlers.HandlerCreateShortURL(w, r, storage, &cur)
	}
}
