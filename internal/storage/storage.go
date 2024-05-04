package storage

import (
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Set(ctx context.Context, userID uuid.UUID, URL string) (string, error)
	MultiSet(ctx context.Context, userID uuid.UUID, items Items) (Items, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (Items, error)
}
