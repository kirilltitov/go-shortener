package http

import (
	"context"
	"encoding/json"
	"io"
	"net"
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

func TestApplication_APIHandlerInternalStats(t *testing.T) {
	_, trustedSubnet, err := net.ParseCIDR("195.248.161.0/24")
	require.NoError(t, err)

	type statsResponse struct {
		Users int `json:"users"`
		URLs  int `json:"urls"`
	}

	type want struct {
		code     int
		response *statsResponse
	}
	tests := []struct {
		name          string
		trustedSubnet *net.IPNet
		realIP        string
		want          want
	}{
		{
			name:          "Positive",
			trustedSubnet: trustedSubnet,
			realIP:        "195.248.161.225",
			want: want{
				code: 200,
				response: &statsResponse{
					Users: 1,
					URLs:  1,
				},
			},
		},
		{
			name:          "No IP",
			trustedSubnet: trustedSubnet,
			realIP:        "",
			want: want{
				code: 401,
			},
		},
		{
			name:          "Disallowed IP",
			trustedSubnet: trustedSubnet,
			realIP:        "77.88.8.8",
			want: want{
				code: 401,
			},
		},
		{
			name:          "No subnet",
			trustedSubnet: nil,
			realIP:        "195.248.161.225",
			want: want{
				code: 401,
			},
		},
	}

	//a, err := New(context.Background(), config.Config{})
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
			a.Shortener.Config.TrustedSubnet = tt.trustedSubnet

			r := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			if tt.realIP != "" {
				r.Header.Set("X-Real-IP", tt.realIP)
			}
			w := httptest.NewRecorder()

			a.APIHandlerInternalStats(w, r)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.response != nil {
				var res statsResponse
				require.NoError(t, json.Unmarshal(resultBody, &res))
				assert.Equal(t, tt.want.response, &res)
			}
		})
	}
}
