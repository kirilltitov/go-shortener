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

func TestApplication_GetURL(t *testing.T) {
	ctx := context.Background()
	authenticatedContext, _ := testhelpers.GetValidUserAndToken()

	type args struct {
		ctx context.Context
		req *gen.GetURLRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *gen.GetURLResponse
		wantErr error
	}{
		{
			name: "Positive",
			args: args{
				ctx: ctx,
				req: &gen.GetURLRequest{ShortUrl: "xA"},
			},
			want: &gen.GetURLResponse{
				OriginalUrl: "https://ya.ru/",
			},
		},
		{
			name: "Negative",
			args: args{
				ctx: ctx,
				req: &gen.GetURLRequest{ShortUrl: "invalid"},
			},
			want:    nil,
			wantErr: ErrBadRequest,
		},
	}

	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""
	cnt, e := container.New(context.Background(), cfg)
	require.NoError(t, e)
	a := &Application{Shortener: shortener.New(cfg, cnt), wg: &sync.WaitGroup{}}

	_, e2 := a.CreateShortURL(authenticatedContext, &gen.CreateShortURLRequest{OriginalUrl: "https://ya.ru/"})
	require.NoError(t, e2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.GetURL(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
