package main

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/app"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func run() error {
	cfg := config.New()
	ctx := context.Background()

	app, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}

	logger.Log.Infof("Starting server at %s", cfg.ServerAddress)

	return app.Run()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
