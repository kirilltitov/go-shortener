package shortener

import (
	"context"
	"fmt"

	"github.com/kirilltitov/go-shortener/internal/utils"
)

func (s *Shortener) ShortenURL(ctx context.Context, URL string) (string, error) {
	if !utils.IsValidURL(URL) {
		return "", fmt.Errorf("invalid URL (must start with https:// or http://): %s", URL)
	}

	shortURL, err := s.container.Storage.Set(ctx, URL)

	return s.FormatShortURL(shortURL), err
}
