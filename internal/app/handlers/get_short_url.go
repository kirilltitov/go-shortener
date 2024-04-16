package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kirilltitov/go-shortener/internal/logger"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
)

func getShortURL(ctx context.Context, shortURL string, storage Storage) (string, error) {
	url, err := storage.Get(ctx, shortURL)
	if err != nil {
		if errors.Is(err, internalStorage.ErrNotFound) {
			return "", fmt.Errorf("URL '%s' not found", shortURL)
		} else {
			return "", err
		}
	}

	return url, nil
}

func HandlerGetShortURL(w http.ResponseWriter, r *http.Request, storage Storage) {
	result, err := getShortURL(r.Context(), chi.URLParam(r, "short"), storage)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	http.Redirect(w, r, result, http.StatusTemporaryRedirect)
}

func APIHandlerGetShortURL(w http.ResponseWriter, r *http.Request, storage Storage) {
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

	result, err := getShortURL(r.Context(), req.URL, storage)

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
