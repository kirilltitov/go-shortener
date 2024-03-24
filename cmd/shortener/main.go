package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

var cur int = 0
var storage handlers.Storage = internalStorage.InMemory{}

func ShortenerRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/{short}", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetShortURL(w, r, storage)
	}))
	router.Post("/", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateShortURL(w, r, storage, &cur)
	}))

	return router
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config.Parse()

	logger.Log.Infof("Starting server at %s", config.GetServerAddress())

	return http.ListenAndServe(config.GetServerAddress(), ShortenerRouter())
}
