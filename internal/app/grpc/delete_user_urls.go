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
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		for err := range a.Shortener.DeleteUserURLs(context.Background(), doneCh, userID, req.UrlsToDel, a.wg) {
			if err != nil {
				logger.Log.Infof("Something went wrong during URL deletion: %s", err)
			}
		}

		close(doneCh)
	}()

	return &gen.DeleteUserURLsResponse{}, nil
}
