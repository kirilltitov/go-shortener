package main

import (
	"github.com/jackc/pgx/v5"

	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

type app struct {
	DB *pgx.Conn
}

func newApp(databaseDSN string) (*app, error) {
	db, err := internalStorage.NewConnection(databaseDSN)
	if err != nil {
		return nil, err
	}

	return &app{DB: db}, nil
}
