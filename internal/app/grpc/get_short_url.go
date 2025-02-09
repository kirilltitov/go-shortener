package grpc

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func (a *Application) GetURL(
	ctx context.Context,
	req *gen.GetURLRequest,
) (*gen.GetURLResponse, error) {
	result, err := a.Shortener.GetURL(ctx, req.ShortUrl)
	if err != nil {
		logger.Log.Error(err)
		return nil, ErrInternal
	}

	return &gen.GetURLResponse{OriginalUrl: result}, nil
}
