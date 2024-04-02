package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jxskiss/base62"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/logger"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/kirilltitov/go-shortener/internal/utils"
)

func createShortURL(URL string, storage internalStorage.Storage, cur *int) (string, error) {
	if !utils.IsValidURL(URL) {
		return "", fmt.Errorf("invalid URL (must start with https:// or http://): %s", URL)
	}

	*cur++
	storage.Set(*cur, URL)

	shortURL := base62.EncodeToString([]byte(strconv.Itoa(*cur)))
	fullShortURL := fmt.Sprintf("%s/%s", config.GetBaseURL(), shortURL)

	if err := internalStorage.SaveRowToFile(config.GetFileStoragePath(), *cur, shortURL, URL); err != nil {
		return "", nil
	}

	return fullShortURL, nil
}

func HandlerCreateShortURL(w http.ResponseWriter, r *http.Request, storage internalStorage.Storage, cur *int) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not read body: %s\n", err.Error()))
		return
	}

	URL := string(b)
	shortURL, err := createShortURL(URL, storage, cur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, shortURL)
}

func APIHandlerCreateShortURL(w http.ResponseWriter, r *http.Request, storage internalStorage.Storage, cur *int) {
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

	shortURL, err := createShortURL(req.URL, storage, cur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Info(err.Error())
		return
	}

	resp := response{Result: shortURL}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}
