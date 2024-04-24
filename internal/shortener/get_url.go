package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/kirilltitov/go-shortener/internal/storage"
)

func (s *Shortener) GetURL(ctx context.Context, shortURL string) (string, error) {
	url, err := s.container.Storage.Get(ctx, shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", fmt.Errorf("URL '%s' not found", shortURL)
		} else {
			return "", err
		}
	}

	return url, nil
}
