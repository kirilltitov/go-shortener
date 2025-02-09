package grpc

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func (a *Application) DeleteUserURLs(
	ctx context.Context,
	req *gen.DeleteUserURLsRequest,
) (*gen.DeleteUserURLsResponse, error) {
	userID, ok := getUserID(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	doneCh := make(chan struct{})
	go func() {
		for err := range a.Shortener.DeleteUserURLs(context.Background(), doneCh, userID, req.UrlsToDel) {
			if err != nil {
				logger.Log.Infof("Something went wrong during URL deletion: %s", err)
			}
		}

		defer close(doneCh)
	}()

	return &gen.DeleteUserURLsResponse{}, nil
}
