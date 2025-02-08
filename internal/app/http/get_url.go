package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

// HandlerGetURL является методом для перехода на оригинальную ссылку с сокращенной.
func (a *Application) HandlerGetURL(w http.ResponseWriter, r *http.Request) {
	result, err := a.Shortener.GetURL(r.Context(), chi.URLParam(r, "short"))

	if err != nil {
		if errors.Is(err, storage.ErrDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
			return
		}
	}

	http.Redirect(w, r, result, http.StatusTemporaryRedirect)
}
