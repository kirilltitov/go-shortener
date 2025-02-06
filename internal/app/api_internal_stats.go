package app

import (
	"encoding/json"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// APIHandlerInternalStats является API-методом для получения статистики по пользователям и ссылкам в сервисе.
// Доступен только из-под разрешенной подсети.
func (a *Application) APIHandlerInternalStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logger.Log

	if !a.isTrustedClientIP(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := a.Shortener.GetStats(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
