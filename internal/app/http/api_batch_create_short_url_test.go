package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/config"
)

func TestAPIHandlerBatchCreateShortURL(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, err := container.New(context.Background(), cfg)
	require.NoError(t, err)
	service := shortener.New(cfg, cnt)
	a := New(service, &sync.WaitGroup{})

	type want struct {
		code     int
		response batchResponse
	}
	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "Positive http/https",
			input: `[{"correlation_id":"lul","original_url":"https://ya.ru"},{"correlation_id":"kek","original_url":"http://ya.ru"}]`,
			want: want{
				code: 201,
				response: []batchResponseRow{
					{CorrelationID: "lul", ShortURL: fmt.Sprintf("%s/xA", a.Shortener.Config.BaseURL)},
					{CorrelationID: "kek", ShortURL: fmt.Sprintf("%s/yA", a.Shortener.Config.BaseURL)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader([]byte(tt.input)))
			w := httptest.NewRecorder()

			a.APIHandlerBatchCreateShortURL(w, r)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.response != nil {
				var res batchResponse
				require.NoError(t, json.Unmarshal(resultBody, &res))
				assert.Equal(t, tt.want.response, res)
			}
		})
	}
}
