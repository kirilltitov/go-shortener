package grpc

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/testhelpers"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
)

func TestApplication_GetInternalStats(t *testing.T) {
	_, trustedSubnet, e := net.ParseCIDR("195.248.161.0/24")
	require.NoError(t, e)

	ctx := context.Background()
	authenticatedContext, _ := testhelpers.GetValidUserAndToken()

	type want struct {
		err      error
		response *gen.GetInternalStatsResponse
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
				response: &gen.GetInternalStatsResponse{
					Stats: &gen.Stats{
						CountUrls:  uint32(1),
						CountUsers: uint32(1),
					},
				},
			},
		},
		{
			name:          "No IP",
			trustedSubnet: trustedSubnet,
			realIP:        "",
			want: want{
				err: ErrUnauthorized,
			},
		},
		{
			name:          "Disallowed IP",
			trustedSubnet: trustedSubnet,
			realIP:        "77.88.8.8",
			want: want{
				err: ErrUnauthorized,
			},
		},
		{
			name:          "No subnet",
			trustedSubnet: nil,
			realIP:        "195.248.161.225",
			want: want{
				err: ErrUnauthorized,
			},
		},
	}

	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, e := container.New(context.Background(), cfg)
	require.NoError(t, e)
	a := &Application{Shortener: shortener.New(cfg, cnt), wg: &sync.WaitGroup{}}

	_, err := a.CreateShortURL(authenticatedContext, &gen.CreateShortURLRequest{OriginalUrl: "https://ya.ru/1"})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Shortener.Config.TrustedSubnet = tt.trustedSubnet

			if tt.realIP != "" {
				ctx = testhelpers.NewContextWithValue("X-Real-IP", tt.realIP)
			} else {
				ctx = context.Background()
			}

			got, gotErr := a.GetInternalStats(ctx, nil)

			if tt.want.response != nil {
				require.NoError(t, gotErr)
				require.Equal(t, tt.want.response, got)
			}

			if tt.want.err != nil {
				require.Equal(t, tt.want.err, gotErr)
			}
		})
	}
}
