package main

import (
	"github.com/jackc/pgx/v5"

	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

type app struct {
	DB *pgx.Conn
}

func newApp(databaseDSN string) (*app, error) {
	var db *pgx.Conn
	if databaseDSN != "" {
		_db, err := internalStorage.NewConnection(databaseDSN)
		if err != nil {
			panic(err)
		}
		db = _db
	}

	return &app{DB: db}, nil
}
