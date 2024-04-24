package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *Application) HandlerGetURL(w http.ResponseWriter, r *http.Request) {
	result, err := a.Shortener.GetURL(r.Context(), chi.URLParam(r, "short"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	http.Redirect(w, r, result, http.StatusTemporaryRedirect)
}
