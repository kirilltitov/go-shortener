package app

import (
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) HandlerPing(w http.ResponseWriter, r *http.Request) {
	code := 200

	switch v := a.Container.Storage.(type) {
	case *storage.PgSQL:
		err := v.C.Ping(r.Context())

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
