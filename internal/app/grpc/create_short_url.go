package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) CreateShortURL(
	ctx context.Context,
	req *gen.CreateShortURLRequest,
) (*gen.CreateShortURLResponse, error) {
	logger.Log.Info("ctx %+v", ctx)
	userID, ok := ctx.Value(ctxUserIDKey).(uuid.UUID)
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
