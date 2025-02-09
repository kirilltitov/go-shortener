package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

// MultiShorten сокращает множество ссылок для переданного пользователя.
// Возвращает список сокращенных ссылок.
func (s *Shortener) MultiShorten(ctx context.Context, userID uuid.UUID, items storage.Items) (storage.Items, error) {
	return s.Container.Storage.MultiSet(ctx, userID, items)
}
