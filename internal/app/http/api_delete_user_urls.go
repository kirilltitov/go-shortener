package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// APIDeleteUserURLs является API-методом для удаления всех сокращенных ссылок для переданного пользователя.
func (a *Application) APIDeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logger.Log

	var req []string
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

	doneCh := make(chan struct{})
	go func() {
		for err := range a.Shortener.DeleteUserURLs(context.Background(), doneCh, *userID, req) {
			if err != nil {
				logger.Log.Infof("Something went wrong during URL deletion: %s", err)
			}
		}

		defer close(doneCh)
	}()

	w.WriteHeader(http.StatusAccepted)
}
