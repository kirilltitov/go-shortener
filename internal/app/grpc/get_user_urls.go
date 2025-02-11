package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func (a *Application) GetUserURLs(
	ctx context.Context,
	req *gen.GetUserURLsRequest,
) (*gen.GetUserURLsResponse, error) {
	userID, ok := getUserID(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	items, err := a.Shortener.GetURLsByUser(ctx, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Could not get URLs by user")
		return nil, ErrInternal
	}

	if len(items) == 0 {
		return nil, status.Error(codes.NotFound, "no urls")
	}

	result := &gen.GetUserURLsResponse{UserUrls: make([]*gen.URL, 0)}
	for _, item := range items {
		result.UserUrls = append(
			result.UserUrls,
			&gen.URL{
				ShortUrl:    item.ShortURL,
				OriginalUrl: item.URL,
			},
		)
	}

	return result, nil
}
