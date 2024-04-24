package storage

import (
	"context"
)

type Storage interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Set(ctx context.Context, URL string) (string, error)
	MultiSet(ctx context.Context, items Items) (Items, error)
}
