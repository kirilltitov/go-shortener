package shortener

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/utils"
)

// ShortenURL сокращает ссылку для переданного пользователя и возвращает сокращенную ссылку.
func (s *Shortener) ShortenURL(ctx context.Context, userID uuid.UUID, URL string) (string, error) {
	if !utils.IsValidURL(URL) {
		return "", fmt.Errorf("invalid URL (must start with https:// or http://): %s", URL)
	}

	shortURL, err := s.container.Storage.Set(ctx, userID, URL)

	return s.FormatShortURL(shortURL), err
}
