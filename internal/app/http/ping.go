package http

import (
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// HandlerPing является методом, возвращающим текущее здоровье сервиса.
func (a *Application) HandlerPing(w http.ResponseWriter, r *http.Request) {
	code := 200

	if err := a.Shortener.Container.Storage.Status(r.Context()); err != nil {
		logger.Log.Errorf("Could not ping storage: %v\n", err)
		code = 500
	}

	w.WriteHeader(code)
}
