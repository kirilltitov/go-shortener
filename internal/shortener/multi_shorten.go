package shortener

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (s *Shortener) MultiShorten(ctx context.Context, items storage.Items) (storage.Items, error) {
	return s.container.Storage.MultiSet(ctx, items)
}
