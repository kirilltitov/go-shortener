package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) CreateShortURL(
	ctx context.Context,
	req *gen.CreateShortURLRequest,
) (*gen.CreateShortURLResponse, error) {
	userID, ok := getUserID(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	shortURL, err := a.Shortener.ShortenURL(ctx, userID, req.OriginalUrl)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			return nil, status.Error(codes.AlreadyExists, "short url already exists")
		} else {
			return nil, ErrInternal
		}
	}

	return &gen.CreateShortURLResponse{ShortUrl: shortURL}, nil
}
