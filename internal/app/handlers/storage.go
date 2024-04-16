package handlers

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

type Storage interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Set(ctx context.Context, URL string) (string, error)
	MultiSet(ctx context.Context, items storage.Items) (storage.Items, error)
}
