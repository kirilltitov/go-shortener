package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kirilltitov/go-shortener/internal/config"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreateShortURL(t *testing.T) {
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
				response: fmt.Sprintf("%s/xA", config.GetBaseURL()),
			},
		},
		{
			name:  "Positive https",
			input: "https://ya.ru",
			want: want{
				code:     201,
				response: fmt.Sprintf("%s/yA", config.GetBaseURL()),
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

	storage := internalStorage.InMemory{}
	cur := 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.input))
			w := httptest.NewRecorder()

			HandlerCreateShortURL(w, r, storage, &cur)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.response, string(resultBody))
		})
	}
}

func TestAPIHandlerCreateShortURL(t *testing.T) {
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
				response: &response{Result: fmt.Sprintf("%s/xA", config.GetBaseURL())},
			},
		},
		{
			name:  "Positive https",
			input: request{URL: "https://ya.ru"},
			want: want{
				code:     201,
				response: &response{Result: fmt.Sprintf("%s/yA", config.GetBaseURL())},
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

	storage := internalStorage.InMemory{}
	cur := 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputBytes, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(inputBytes))
			w := httptest.NewRecorder()

			APIHandlerCreateShortURL(w, r, storage, &cur)

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
