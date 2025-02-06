package shortener

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

// GetStats возвращает статистику использования сервиса, см. структуру [storage.Stats].
func (s *Shortener) GetStats(ctx context.Context) (*storage.Stats, error) {
	return s.container.Storage.GetStats(ctx)
}
