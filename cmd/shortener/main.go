// Модуль main запускает веб-сервер сокращателя ссылок.
// См. [github.com/kirilltitov/go-shortener/internal/config.Config]
// на предмет конфигурационных параметров сервиса.
package main

import (
	"context"
	_ "net/http/pprof"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func run() error {
	cfg := config.New()
	ctx := context.Background()

	a, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}

	logger.Log.Infof("Starting server at %s", cfg.ServerAddress)

	return a.Run()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
