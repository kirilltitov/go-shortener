package handlers

import (
	"context"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func HandlerPing(w http.ResponseWriter, r *http.Request, ctx context.Context, s Storage) {
	code := 200

	switch v := s.(type) {
	case storage.PgSQL:
		err := v.C.Ping(ctx)

		if err != nil {
			logger.Log.Errorf("Could not ping PgSQL: %v\n", err)
			code = 500
		}
	default:
		logger.Log.Info("Storage is not PgSQL")
		code = 500
	}

	w.WriteHeader(code)
}
