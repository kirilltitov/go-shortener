package app

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// APIHandlerGetURL является API-методом для получения полной ссылки из переданного короткого идентификатора.
func (a *Application) APIHandlerGetURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logger.Log

	var req request
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		log.Infof("Could not get body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		log.Infof("Could not parse request JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := a.Shortener.GetURL(r.Context(), req.URL)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Info(err.Error())
		return
	}

	responseBytes, err := json.Marshal(response{Result: result})
	if err != nil {
		panic(err)
	}

	w.Write(responseBytes)
}
