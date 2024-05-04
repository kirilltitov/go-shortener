package shortener

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (s *Shortener) MultiShorten(ctx context.Context, userID uuid.UUID, items storage.Items) (storage.Items, error) {
	return s.container.Storage.MultiSet(ctx, userID, items)
}
