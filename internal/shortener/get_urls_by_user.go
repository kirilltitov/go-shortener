package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (s *Shortener) GetURLsByUser(ctx context.Context, userID uuid.UUID) (storage.Items, error) {
	return s.container.Storage.GetByUser(ctx, userID)
}
