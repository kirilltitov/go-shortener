package container

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

// Container хранит в себе зависимости сервиса.
type Container struct {
	Storage storage.Storage
}

// New создает, конфигурирует и возвращает объект контейнера зависимостей.
func New(ctx context.Context, cfg config.Config) (*Container, error) {
	s, err := newStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Container{Storage: s}, nil
}

func newStorage(ctx context.Context, cfg config.Config) (storage.Storage, error) {
	var s storage.Storage

	if cfg.DatabaseDSN != "" {
		logger.Log.Info("Using storage PgSQL")
		_s, err := storage.NewPgSQLStorage(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		if err := _s.MigrateUp(ctx); err != nil {
			return nil, err
		}
		s = _s
	} else if cfg.FileStoragePath != "" {
		logger.Log.Info("Using storage file")
		_s, err := storage.NewFileStorage(ctx, cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
		s = _s
	} else {
		logger.Log.Info("Using storage memory")
		s = storage.NewInMemoryStorage(ctx)
	}

	return s, nil
}
