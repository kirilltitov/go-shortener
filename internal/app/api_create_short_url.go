package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) APIHandlerCreateShortURL(w http.ResponseWriter, r *http.Request) {
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

	code := http.StatusCreated
	shortURL, err := a.Shortener.ShortenURL(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			code = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			log.Info(err.Error())
			return
		}
	}

	resp := response{Result: shortURL}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	w.Write(responseBytes)
}
