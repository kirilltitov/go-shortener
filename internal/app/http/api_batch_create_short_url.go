package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

type batchRequest []batchRequestRow
type batchRequestRow struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type batchResponse []batchResponseRow
type batchResponseRow struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// APIHandlerBatchCreateShortURL является API-методом для сокращения множества ссылок.
func (a *Application) APIHandlerBatchCreateShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logger.Log

	var b batchRequest
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		log.Infof("Could not get body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &b); err != nil {
		log.Infof("Could not parse request JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := a.authenticate(r, w, false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not authenticate user: %s\n", err.Error()))
		return
	}

	if userID == nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "User not authenticated")
		return
	}

	var items storage.Items
	for _, r := range b {
		items = append(items, storage.Item{
			UUID: r.CorrelationID,
			URL:  r.OriginalURL,
		})
		logger.Log.Infof("Loaded row %+v from body", r)
	}

	result, err := a.Shortener.MultiShorten(r.Context(), *userID, items)
	if err != nil {
		log.Infof("Could not batch insert: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Infof("Inserted items %+v", result)

	var response batchResponse
	for _, item := range result {
		response = append(response, batchResponseRow{
			CorrelationID: item.UUID,
			ShortURL:      a.Shortener.FormatShortURL(item.URL),
		})
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}
