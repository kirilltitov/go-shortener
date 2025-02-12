package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	Config     config.Config
	Container  *container.Container
	Shortener  shortener.Shortener
	HTTPServer *http.Server

	wg *sync.WaitGroup
}

// New создает и возвращает сконфигурированный объект веб-приложения сервиса.
func New(ctx context.Context, cfg config.Config) (*Application, error) {
	cnt, err := container.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	a := &Application{
		Config:     cfg,
		Container:  cnt,
		Shortener:  shortener.New(cfg, cnt),
		HTTPServer: &http.Server{Addr: cfg.ServerAddress},
		wg:         &sync.WaitGroup{},
	}

	a.HTTPServer.Handler = utils.GzipHandle(a.createRouter())

	return a, nil
}

// Run запускает веб-сервер приложения. Может вернуть ошибку при ошибке конфигурации (занятый порт и т. д.).
func (a *Application) Run() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		var runFunc func() error

		if a.Config.EnableHTTPS == "" {
			runFunc = func() error {
				logger.Log.Infof("Starting a HTTP server")
				return a.HTTPServer.ListenAndServe()
			}
		} else {
			runFunc = func() error {
				logger.Log.Infof("Starting a HTTPS server")
				return a.HTTPServer.ListenAndServeTLS(
					"localhost.crt",
					"localhost.key",
				)
			}
		}

		if err := runFunc(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Log.Info("HTTP server shutdown")
			} else {
				logger.Log.Error(err)
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan
	logger.Log.Infof("Received signal: %v", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		logger.Log.Info("Shutting down HTTP server")
		if err := a.HTTPServer.Shutdown(shutdownCtx); err != nil {
			logger.Log.WithError(err).Error("Could not shutdown HTTP server properly")
		}
	}()

	a.wg.Wait()

	a.Container.Storage.Close()

	logger.Log.Info("Goodbye")
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
