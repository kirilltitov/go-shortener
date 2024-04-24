package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

type Application struct {
	Config    config.Config
	Container *container.Container
	Shortener shortener.Shortener
}

func New(ctx context.Context, cfg config.Config) (*Application, error) {
	cnt, err := container.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Application{
		Config:    cfg,
		Container: cnt,
		Shortener: shortener.New(cfg, cnt),
	}, nil
}

func (a *Application) Run() error {
	return http.ListenAndServe(a.Config.ServerAddress, utils.GzipHandle(a.createRouter()))
}

func (a *Application) createRouter() chi.Router {
	router := chi.NewRouter()

	router.Post("/", logger.WithLogging(a.HandlerCreateShortURL))
	router.Get("/{short}", logger.WithLogging(a.HandlerGetURL))
	router.Get("/ping", logger.WithLogging(a.HandlerPing))

	router.Post("/api/get", logger.WithLogging(a.APIHandlerGetURL))
	router.Post("/api/shorten", logger.WithLogging(a.APIHandlerCreateShortURL))
	router.Post("/api/shorten/batch", logger.WithLogging(a.APIHandlerBatchCreateShortURL))

	return router
}
