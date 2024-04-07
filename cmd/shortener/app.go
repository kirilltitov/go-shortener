package main

import (
	"github.com/jackc/pgx/v5"

	"github.com/kirilltitov/go-shortener/internal/logger"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

type app struct {
	DB *pgx.Conn
}

func newApp(databaseDSN string) (*app, error) {
	db, err := internalStorage.NewConnection(databaseDSN)
	if err != nil {
		logger.Log.Warnf("Could not connect to DB: %v", err)
	}

	return &app{DB: db}, nil
}
