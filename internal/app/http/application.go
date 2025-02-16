package http

import (
	"errors"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

// Application является объектом веб-приложения сервиса.
type Application struct {
	Shortener shortener.Shortener
	Server    *http.Server

	wg *sync.WaitGroup
}

// New создает и возвращает сконфигурированный объект веб-приложения сервиса.
func New(s shortener.Shortener, wg *sync.WaitGroup) *Application {
	a := &Application{
		Shortener: s,
		Server:    &http.Server{Addr: s.Config.ServerAddress},
		wg:        wg,
	}

	a.Server.Handler = utils.GzipHandle(a.createRouter())

	return a
}

// Run запускает веб-сервер приложения.
func (a *Application) Run() {
	defer a.wg.Done()

	var runFunc func() error

	if a.Shortener.Config.EnableHTTPS == "" {
		runFunc = func() error {
			logger.Log.Infof("Starting a HTTP server at %s", a.Shortener.Config.ServerAddress)
			return a.Server.ListenAndServe()
		}
	} else {
		runFunc = func() error {
			logger.Log.Infof("Starting a HTTPS server at %s", a.Shortener.Config.ServerAddress)
			return a.Server.ListenAndServeTLS(
				"localhost.crt",
				"localhost.key",
			)
		}
	}

	if err := runFunc(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server shutdown")
		} else {
			panic(err)
		}
	}
}

func (a *Application) createRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(logger.WithLogging, utils.GzipHandle)

	router.Mount("/debug", middleware.Profiler())

	router.Post("/", a.HandlerCreateShortURL)
	router.Get("/{short}", a.HandlerGetURL)
	router.Get("/ping", a.HandlerPing)

	router.Post("/api/get", a.APIHandlerGetURL)
	router.Get("/api/user/urls", a.APIUserURLs)
	router.Delete("/api/user/urls", a.APIDeleteUserURLs)
	router.Post("/api/shorten", a.APIHandlerCreateShortURL)
	router.Post("/api/shorten/batch", a.APIHandlerBatchCreateShortURL)

	router.Get("/api/internal/stats", a.APIHandlerInternalStats)

	return router
}
