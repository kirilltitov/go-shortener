package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

func createShortURL(ctx context.Context, URL string, storage Storage) (string, error) {
	if !utils.IsValidURL(URL) {
		return "", fmt.Errorf("invalid URL (must start with https:// or http://): %s", URL)
	}

	shortURL, err := storage.Set(ctx, URL)

	return formatShortURL(shortURL), err
}

func formatShortURL(shortURL string) string {
	return fmt.Sprintf("%s/%s", config.GetBaseURL(), shortURL)
}

func HandlerCreateShortURL(w http.ResponseWriter, r *http.Request, storage Storage) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not read body: %s\n", err.Error()))
		return
	}

	code := http.StatusCreated
	URL := string(b)
	shortURL, err := createShortURL(r.Context(), URL, storage)
	if err != nil {
		if errors.Is(err, internalStorage.ErrDuplicate) {
			code = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
			return
		}
	}

	w.WriteHeader(code)
	io.WriteString(w, shortURL)
}

func APIHandlerCreateShortURL(w http.ResponseWriter, r *http.Request, storage Storage) {
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
	shortURL, err := createShortURL(r.Context(), req.URL, storage)
	if err != nil {
		if errors.Is(err, internalStorage.ErrDuplicate) {
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
