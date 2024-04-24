package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/config"
)

func TestAPIHandlerCreateShortURL(t *testing.T) {
	a, err := New(context.Background(), config.Config{})
	require.NoError(t, err)

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
				response: &response{Result: fmt.Sprintf("%s/xA", a.Config.BaseURL)},
			},
		},
		{
			name:  "Positive https",
			input: request{URL: "https://ya.ru"},
			want: want{
				code:     201,
				response: &response{Result: fmt.Sprintf("%s/yA", a.Config.BaseURL)},
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
