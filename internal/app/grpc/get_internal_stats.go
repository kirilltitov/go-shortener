package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/shortener"
)

func (a *Application) GetInternalStats(
	ctx context.Context,
	req *gen.GetInternalStatsRequest,
) (*gen.GetInternalStatsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("X-Real-IP")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}

	ip := values[0]
	result, err := a.Shortener.GetStats(ctx, ip)

	if err != nil {
		if errors.Is(err, shortener.ErrorUnauthorized) {
			return nil, ErrUnauthorized
		} else {
			return nil, ErrInternal
		}
	}

	return &gen.GetInternalStatsResponse{
		Stats: &gen.Stats{
			CountUrls:  uint32(result.URLs),
			CountUsers: uint32(result.Users),
		},
	}, nil
}
