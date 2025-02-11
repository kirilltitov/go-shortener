package grpc

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/testhelpers"
	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/shortener"
)

func TestApplication_GetUserURLs(t *testing.T) {
	ctx := context.Background()
	authenticatedContext, _ := testhelpers.GetValidUserAndToken()

	type args struct {
		ctx context.Context
		req *gen.GetUserURLsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *gen.GetUserURLsResponse
		wantErr error
	}{
		{
			name: "No auth",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: ErrUnauthorized,
		},
		{
			name: "Positive",
			args: args{
				ctx: authenticatedContext,
				req: &gen.GetUserURLsRequest{},
			},
			want: &gen.GetUserURLsResponse{
				UserUrls: []*gen.URL{
					{
						ShortUrl:    "http://localhost:8080/xA",
						OriginalUrl: "https://ya.ru/1",
					},
					{
						ShortUrl:    "http://localhost:8080/yA",
						OriginalUrl: "https://ya.ru/2",
					},
				},
			},
			wantErr: nil,
		},
	}

	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, e := container.New(context.Background(), cfg)
	require.NoError(t, e)
	a := &Application{Shortener: shortener.New(cfg, cnt), wg: &sync.WaitGroup{}}

	_, e2 := a.CreateShortURL(authenticatedContext, &gen.CreateShortURLRequest{OriginalUrl: "https://ya.ru/1"})
	require.NoError(t, e2)
	_, e3 := a.CreateShortURL(authenticatedContext, &gen.CreateShortURLRequest{OriginalUrl: "https://ya.ru/2"})
	require.NoError(t, e3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.GetUserURLs(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
