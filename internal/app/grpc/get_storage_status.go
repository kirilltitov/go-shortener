package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (a *Application) GetStorageStatus(
	ctx context.Context,
	req *gen.GetStorageStatusRequest,
) (*gen.GetStorageStatusResponse, error) {
	switch v := a.Shortener.Container.Storage.(type) {
	case *storage.PgSQL:
		err := v.C.Ping(ctx)
		if err != nil {
			logger.Log.Errorf("Could not ping PgSQL: %v\n", err)
			return nil, status.Error(codes.Unavailable, "db is not available")
		}
	default:
		logger.Log.Info("Storage is not PgSQL")
		return nil, status.Error(codes.Unimplemented, "storage is not pgsql")
	}

	return nil, status.Error(codes.OK, "everything is fine")
}
