package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/shortener"
)

// APIHandlerInternalStats является API-методом для получения статистики по пользователям и ссылкам в сервисе.
// Доступен только из-под разрешенной подсети.
func (a *Application) APIHandlerInternalStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logger.Log

	result, err := a.Shortener.GetStats(r.Context(), r.Header.Get("X-Real-IP"))

	if err != nil {
		if errors.Is(err, shortener.ErrorUnauthorized) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		log.Info(err.Error())
		return
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.Write(responseBytes)
}
