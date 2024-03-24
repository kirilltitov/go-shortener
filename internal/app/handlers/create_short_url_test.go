package handlers

import (
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
				response: "Invalid link (must start with https:// or http://): ya.ru\n",
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
			result.Cookies()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.response, string(resultBody))
		})
	}
}
