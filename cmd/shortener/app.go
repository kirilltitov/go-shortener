package main

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	"github.com/kirilltitov/go-shortener/internal/logger"

	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

type app struct {
	Storage handlers.Storage
}

func newApp(ctx context.Context, databaseDSN string, fileStoragePath string) (*app, error) {
	s, err := newStorage(ctx, databaseDSN, fileStoragePath)
	if err != nil {
		return nil, err
	}

	return &app{Storage: s}, nil
}

func newStorage(ctx context.Context, databaseDSN string, fileStoragePath string) (handlers.Storage, error) {
	var s handlers.Storage

	if databaseDSN != "" {
		logger.Log.Info("Using storage PgSQL")
		_s, err := internalStorage.NewPgSQLStorage(ctx, databaseDSN)
		if err != nil {
			return nil, err
		}
		if err := _s.MigrateUp(ctx); err != nil {
			return nil, err
		}
		s = _s
	} else if fileStoragePath != "" {
		logger.Log.Info("Using storage file")
		_s, err := internalStorage.NewFileStorage(ctx, fileStoragePath)
		if err != nil {
			return nil, err
		}
		s = _s
	} else {
		logger.Log.Info("Using storage memory")
		s = internalStorage.NewInMemoryStorage(ctx)
	}

	return s, nil
}
