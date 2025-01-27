// Пакет storage инкапсулирует логику обращения к хранилищу сокращенных ссылок.
package storage

import (
	"context"

	"github.com/google/uuid"
)

// Storage описывает интерфейс хранилища для сокращенных ссылок.
type Storage interface {
	// Get возвращает сокращенную ссылку по её короткому идентификатору, либо ошибку [ErrNotFound].
	Get(ctx context.Context, shortURL string) (string, error)

	// Set создает новую сокращенную ссылку и возвращает её сокращенный идентификатор.
	Set(ctx context.Context, userID uuid.UUID, URL string) (string, error)

	// MultiSet создает множество сокращенных ссылок и возвращает их сокращенные идентификаторы.
	MultiSet(ctx context.Context, userID uuid.UUID, items Items) (Items, error)

	// GetByUser возвращает все сокращенные ссылки для данного пользователя.
	GetByUser(ctx context.Context, userID uuid.UUID) (Items, error)

	// DeleteByUser удаляет сокращенную ссылку для данного пользователя.
	DeleteByUser(ctx context.Context, userID uuid.UUID, shortURL string) error
}
