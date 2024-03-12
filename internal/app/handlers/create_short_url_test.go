package handlers

import (
	storage2 "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
				response: "http://localhost:8080/xA",
			},
		},
		{
			name:  "Positive https",
			input: "https://ya.ru",
			want: want{
				code:     201,
				response: "http://localhost:8080/yA",
			},
		},
		{
			name:  "Negative invalid",
			input: "ya.ru",
			want: want{
				code:     400,
				response: "Invalid link (must start with https:// or http://): ya.ru\n",
			},
		},
	}

	storage := storage2.InMemory{}
	cur := 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.input))
			w := httptest.NewRecorder()

			HandlerCreateShortURL(w, r, storage, &cur)

			result := w.Result()
			result.Cookies()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.response, string(resultBody))
		})
	}
}
