package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

// GetURLsByUser возвращает все сокращенные ссылки для переданного пользователя.
func (s *Shortener) GetURLsByUser(ctx context.Context, userID uuid.UUID) (storage.Items, error) {
	result, err := s.Container.Storage.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, item := range result {
		item.ShortURL = s.FormatShortURL(item.ShortURL)
		result[i] = item
	}

	return result, nil
}
