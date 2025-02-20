package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/config"
)

func TestAPIHandlerGetShortURL(t *testing.T) {
	type want struct {
		code     int
		response *response
	}
	tests := []struct {
		name  string
		input request
		want  want
	}{
		{
			name:  "Positive",
			input: request{URL: "xA"},
			want: want{
				code:     200,
				response: &response{Result: "https://ya.ru"},
			},
		},
		{
			name:  "Negative",
			input: request{URL: "asdf"},
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
	a := New(service, &sync.WaitGroup{})

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
	w := httptest.NewRecorder()
	a.HandlerCreateShortURL(w, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputBytes, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			r := httptest.NewRequest(http.MethodPost, "/api/get", bytes.NewReader(inputBytes))
			w := httptest.NewRecorder()

			a.APIHandlerGetURL(w, r)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.response != nil {
				var res response
				require.NoError(t, json.Unmarshal(resultBody, &res))
				assert.Equal(t, *tt.want.response, res)
			}
		})
	}
}
