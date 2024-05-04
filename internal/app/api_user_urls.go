package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *Application) APIUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := a.authenticate(r, w)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not authenticate user: %s\n", err.Error()))
		return
	}

	if userID == nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, fmt.Sprintf("User not authenticated: %s\n", err.Error()))
		return
	}

	items, err := a.Shortener.GetURLsByUser(r.Context(), *userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not get URLs by user: %s\n", err.Error()))
		return
	}

	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type row struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var result []row
	for _, item := range items {
		result = append(
			result,
			row{ShortURL: item.ShortURL, OriginalURL: item.URL},
		)
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	w.Write(responseBytes)
}
