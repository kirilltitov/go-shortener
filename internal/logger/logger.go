package logger

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Log является логгером.
var Log = log.New()

func init() {
	Log.SetLevel(log.InfoLevel)
	Log.SetFormatter(&log.JSONFormatter{})
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write записывает переданный массив байтов в собственный буффер.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader записывает код статуса.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging осуществляет логирование переданной функции обработки HTTP-запросов.
func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.WithFields(log.Fields{
			"uri":         r.RequestURI,
			"method":      r.Method,
			"status":      responseData.status,
			"duration_μs": duration.Microseconds(),
			"size":        responseData.size,
		}).Info("Served HTTP request")
	}
}
