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

func TestApplication_BatchCreateShortURL(t *testing.T) {
	ctx := context.Background()
	authenticatedContext, _ := testhelpers.GetValidUserAndToken()

	type args struct {
		ctx context.Context
		req *gen.BatchCreateShortURLRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *gen.BatchCreateShortURLResponse
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
				req: &gen.BatchCreateShortURLRequest{
					BatchUrlRequests: []*gen.URLRequest{
						{
							CorrelationId: "1",
							OriginalUrl:   "https://ya.ru/1",
						},
						{
							CorrelationId: "2",
							OriginalUrl:   "https://ya.ru/2",
						},
					},
				},
			},
			want: &gen.BatchCreateShortURLResponse{
				BatchUrlResponses: []*gen.URLResponse{
					{
						CorrelationId: "1",
						ShortUrl:      "http://localhost:8080/xA",
					},
					{
						CorrelationId: "2",
						ShortUrl:      "http://localhost:8080/yA",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.BatchCreateShortURL(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
