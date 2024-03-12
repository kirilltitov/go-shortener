package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	storage2 "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerGetShortURL(t *testing.T) {
	type want struct {
		code        int
		headerName  string
		headerValue string
	}
	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "Positive",
			input: "xA",
			want: want{
				code:        307,
				headerName:  "Location",
				headerValue: "https://ya.ru",
			},
		},
		{
			name:  "Negative",
			input: "asdf",
			want: want{
				code: 400,
			},
		},
	}

	storage := storage2.InMemory{}
	cur := 0

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
	w := httptest.NewRecorder()
	HandlerCreateShortURL(w, r, storage, &cur)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.input), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("short", tt.input)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			HandlerGetShortURL(w, r, storage)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.headerName != "" {
				assert.Equal(t, tt.want.headerValue, result.Header.Get(tt.want.headerName))
			}
		})
	}
}
