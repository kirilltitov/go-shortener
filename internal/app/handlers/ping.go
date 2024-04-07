package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

func HandlerPing(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	err := conn.Ping(ctx)
	code := 200

	if err != nil {
		logger.Log.Errorf("Could not ping PgSQL: %v\n", err)
		code = 500
	}

	w.WriteHeader(code)
}
