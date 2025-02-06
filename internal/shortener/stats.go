package shortener

import (
	"context"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (s *Shortener) GetStats(ctx context.Context) (*storage.Stats, error) {
	return s.container.Storage.GetStats(ctx)
}
