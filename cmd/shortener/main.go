package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

func ShortenerRouter(a *app) chi.Router {
	router := chi.NewRouter()
	ctx := context.Background()

	router.Get("/{short}", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetShortURL(w, r, a.Storage)
	}))
	router.Post("/api/get", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.APIHandlerGetShortURL(w, r, a.Storage)
	}))

	router.Post("/", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateShortURL(w, r, a.Storage)
	}))
	router.Post("/api/shorten", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.APIHandlerCreateShortURL(w, r, a.Storage)
	}))

	router.Get("/ping", logger.WithLogging(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerPing(w, r, ctx, a.Storage)
	}))

	return router
}

func run() error {
	config.Parse()

	ctx := context.Background()

	a, err := newApp(ctx, config.GetDatabaseDSN(), config.GetFileStoragePath())
	if err != nil {
		return err
	}

	logger.Log.Infof("Starting server at %s", config.GetServerAddress())

	return http.ListenAndServe(config.GetServerAddress(), utils.GzipHandle(ShortenerRouter(a)))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
