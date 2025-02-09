package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/config"
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

	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(t, err)
	service := shortener.New(cfg, cnt)
	a := New(service)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
	w := httptest.NewRecorder()
	a.HandlerCreateShortURL(w, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.input), nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("short", tt.input)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			a.HandlerGetURL(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.headerName != "" {
				assert.Equal(t, tt.want.headerValue, result.Header.Get(tt.want.headerName))
			}
		})
	}
}
