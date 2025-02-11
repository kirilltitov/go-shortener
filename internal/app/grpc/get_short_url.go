package grpc

import (
	"context"
	"errors"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) GetURL(
	ctx context.Context,
	req *gen.GetURLRequest,
) (*gen.GetURLResponse, error) {
	result, err := a.Shortener.GetURL(ctx, req.ShortUrl)
	if err != nil {
		logger.Log.Error(err)

		if errors.Is(err, storage.ErrDeleted) {
			return nil, ErrGone
		} else {
			return nil, ErrBadRequest
		}
	}

	return &gen.GetURLResponse{OriginalUrl: result}, nil
}
