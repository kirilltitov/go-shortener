package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func TestWithLogging(t *testing.T) {
	ja := jsonassert.New(t)

	buf := bytes.Buffer{}
	Log.SetOutput(&buf)

	s := storage.InMemory{}
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()

	WithLogging(func(writer http.ResponseWriter, reader *http.Request) {
		handlers.HandlerGetShortURL(writer, reader, s)
	})(w, r)

	ja.Assertf(
		buf.String(),
		`{
			"duration_μs": "<<PRESENCE>>",
			"level": "info",
			"method": "GET",
			"msg": "Served HTTP request",
			"size": 30,
			"status": 400,
			"time": "<<PRESENCE>>",
			"uri": "/abc"
		}`,
	)
}
