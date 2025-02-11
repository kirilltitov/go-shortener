package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func BenchmarkApplication_APIHandlerGetURL(b *testing.B) {
	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(b, err)
	service := shortener.New(cfg, cnt)
	a := New(service, &sync.WaitGroup{})

	createURLBytes, err := json.Marshal(request{URL: "https://ya.ru"})
	assert.NoError(b, err)
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(createURLBytes))
	w := httptest.NewRecorder()
	a.APIHandlerCreateShortURL(w, r)

	result := w.Result()
	defer result.Body.Close()
	resultBody, err := io.ReadAll(result.Body)
	require.NoError(b, err)
	var res response
	require.NoError(b, json.Unmarshal(resultBody, &res))

	inputBytes, err := json.Marshal(request{URL: res.Result[strings.LastIndex(res.Result, "/")+1:]})
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/get", bytes.NewReader(inputBytes))
		w := httptest.NewRecorder()

		a.APIHandlerGetURL(w, r)

		result := w.Result()
		defer result.Body.Close()
	}
}

func TestAPIHandlerCreateShortURL(t *testing.T) {
	cfg := config.New()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(t, err)
	service := shortener.New(cfg, cnt)
	a := New(service, &sync.WaitGroup{})

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
			name:  "Positive http",
			input: request{URL: "http://ya.ru"},
			want: want{
				code:     201,
				response: &response{Result: fmt.Sprintf("%s/xA", a.Shortener.Config.BaseURL)},
			},
		},
		{
			name:  "Positive https",
			input: request{URL: "https://ya.ru"},
			want: want{
				code:     201,
				response: &response{Result: fmt.Sprintf("%s/yA", a.Shortener.Config.BaseURL)},
			},
		},
		{
			name:  "Negative invalid",
			input: request{URL: "ya.ru"},
			want: want{
				code:     400,
				response: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputBytes, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(inputBytes))
			w := httptest.NewRecorder()

			a.APIHandlerCreateShortURL(w, r)

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
