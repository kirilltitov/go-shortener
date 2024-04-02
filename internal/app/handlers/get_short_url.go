package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jxskiss/base62"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func getShortURL(shortURL string, storage storage.Storage) (string, error) {
	decodedStringInt, err := base62.DecodeString(shortURL)
	if err != nil {
		return "", fmt.Errorf("could not decode short url '%s'", shortURL)
	}

	decodedInt, err := strconv.Atoi(string(decodedStringInt))
	if err != nil {
		return "", fmt.Errorf("could not decode short url '%s'", shortURL)
	}

	url, ok := storage.Get(decodedInt)
	if !ok {
		return "", fmt.Errorf("URL '%s' not found", shortURL)
	}

	return url, nil
}

func HandlerGetShortURL(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	result, err := getShortURL(chi.URLParam(r, "short"), storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	http.Redirect(w, r, result, http.StatusTemporaryRedirect)
}

func APIHandlerGetShortURL(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
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

	result, err := getShortURL(req.URL, storage)

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
