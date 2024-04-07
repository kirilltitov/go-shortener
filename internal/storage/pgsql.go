package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewConnection(DSN string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), DSN)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
