package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirilltitov/go-shortener/internal/app/grpc/gen"
	"github.com/kirilltitov/go-shortener/internal/logger"
	"github.com/kirilltitov/go-shortener/internal/storage"
	log "github.com/sirupsen/logrus"
)

func (a *Application) BatchCreateShortURL(
	ctx context.Context,
	req *gen.BatchCreateShortURLRequest,
) (*gen.BatchCreateShortURLResponse, error) {
	userID, ok := ctx.Value(ctxUserIDKey).(uuid.UUID)
	if !ok {
		return nil, ErrUnauthorized
	}

	var items storage.Items
	for _, r := range req.BatchUrlRequests {
		items = append(items, storage.Item{
			UUID: r.CorrelationId,
			URL:  r.OriginalUrl,
		})
		logger.Log.Infof("Loaded row %+v from body", r)
	}

	result, err := a.Shortener.MultiShorten(ctx, userID, items)
	if err != nil {
		log.Infof("Could not batch insert: %v", err)
		return nil, ErrInternal
	}
	log.Infof("Inserted items %+v", result)

	response := gen.BatchCreateShortURLResponse{
		BatchUrlResponses: make([]*gen.URLResponse, len(result)),
	}
	for _, item := range result {
		response.BatchUrlResponses = append(response.BatchUrlResponses, &gen.URLResponse{
			CorrelationId: item.UUID,
			ShortUrl:      a.Shortener.FormatShortURL(item.URL),
		})
	}
	return &response, nil
}
