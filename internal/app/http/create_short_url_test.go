package http

import (
	"context"
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

func TestHandlerCreateShortURL(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(t, err)
	service := shortener.New(cfg, cnt)
	a := New(service, &sync.WaitGroup{})

	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "Positive http",
			input: "http://ya.ru",
			want: want{
				code:     201,
				response: fmt.Sprintf("%s/xA", a.Shortener.Config.BaseURL),
			},
		},
		{
			name:  "Positive https",
			input: "https://ya.ru",
			want: want{
				code:     201,
				response: fmt.Sprintf("%s/yA", a.Shortener.Config.BaseURL),
			},
		},
		{
			name:  "Negative invalid",
			input: "ya.ru",
			want: want{
				code:     400,
				response: "invalid URL (must start with https:// or http://): ya.ru\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.input))
			w := httptest.NewRecorder()

			a.HandlerCreateShortURL(w, r)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.response, string(resultBody))
		})
	}
}
