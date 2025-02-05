package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

// Application является объектом веб-приложения сервиса.
type Application struct {
	Config    config.Config
	Container *container.Container
	Shortener shortener.Shortener
}

// New создает и возвращает сконфигурированный объект веб-приложения сервиса.
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

// Run запускает веб-сервер приложения. Может вернуть ошибку при ошибке конфигурации (занятый порт и т. д.).
func (a *Application) Run() error {
	handler := utils.GzipHandle(a.createRouter())

	if a.Config.EnableHTTPS == "" {
		logger.Log.Infof("Starting a HTTP server")
		return http.ListenAndServe(a.Config.ServerAddress, handler)
	} else {
		logger.Log.Infof("Starting a HTTPS server")
		return http.ListenAndServeTLS(
			a.Config.ServerAddress,
			"localhost.crt",
			"localhost.key",
			handler,
		)
	}
}

func (a *Application) createRouter() chi.Router {
	router := chi.NewRouter()

	router.Mount("/debug", middleware.Profiler())

	router.Post("/", logger.WithLogging(a.HandlerCreateShortURL))
	router.Get("/{short}", logger.WithLogging(a.HandlerGetURL))
	router.Get("/ping", logger.WithLogging(a.HandlerPing))

	router.Post("/api/get", logger.WithLogging(a.APIHandlerGetURL))
	router.Get("/api/user/urls", logger.WithLogging(a.APIUserURLs))
	router.Delete("/api/user/urls", logger.WithLogging(a.APIDeleteUserURLs))
	router.Post("/api/shorten", logger.WithLogging(a.APIHandlerCreateShortURL))
	router.Post("/api/shorten/batch", logger.WithLogging(a.APIHandlerBatchCreateShortURL))

	return router
}
