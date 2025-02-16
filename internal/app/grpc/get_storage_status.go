package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

func (a *Application) GetStorageStatus(
	ctx context.Context,
	req *gen.GetStorageStatusRequest,
) (*gen.GetStorageStatusResponse, error) {
	if err := a.Shortener.Container.Storage.Status(ctx); err != nil {
		logger.Log.Errorf("Could not ping storage: %v\n", err)
		return nil, status.Error(codes.Unavailable, "db is not available")
	}

	return nil, status.Error(codes.OK, "everything is fine")
}
