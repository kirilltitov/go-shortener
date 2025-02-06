package storage

import (
	"context"
	"strconv"

	"github.com/google/uuid"
)

type inMemoryRow struct {
	URL    string
	cur    int
	userID uuid.UUID
}

// InMemory является хранилищем для сокращенных ссылок в памяти текущего процесса.
type InMemory struct {
	storage map[string]inMemoryRow
	cur     *int
}

// NewInMemoryStorage создает и возвращает экземпляр хранилища в памяти.
func NewInMemoryStorage(ctx context.Context) *InMemory {
	return &InMemory{
		storage: make(map[string]inMemoryRow),
		cur:     new(int),
	}
}

// Get загружает из хранилища информацию о сокращенной ссылке.
func (s InMemory) Get(ctx context.Context, shortURL string) (string, error) {
	var _err error = nil
	val, ok := s.storage[shortURL]
	if !ok {
		_err = ErrNotFound
	}

	return val.URL, _err
}

// Set записывает в хранилище информацию о сокращенной ссылке.
func (s InMemory) Set(ctx context.Context, userID uuid.UUID, URL string) (string, error) {
	*s.cur++
	shortURL := intToShortURL(*s.cur)
	s.storage[shortURL] = inMemoryRow{
		cur:    *s.cur,
		userID: userID,
		URL:    URL,
	}

	return shortURL, nil
}

// MultiSet записывает в хранилище информацию о нескольких сокращенных ссылках.
func (s InMemory) MultiSet(ctx context.Context, userID uuid.UUID, items Items) (Items, error) {
	var result Items

	for _, item := range items {
		shortURL, err := s.Set(ctx, userID, item.URL)
		if err != nil {
			return nil, err
		}
		result = append(result, Item{
			UUID: item.UUID,
			URL:  shortURL,
		})
	}

	return result, nil
}

// GetByUser загружает из хранилища все сокращенные ссылки пользователя.
func (s InMemory) GetByUser(ctx context.Context, userID uuid.UUID) (Items, error) {
	var result Items

	for shortURL, item := range s.storage {
		if item.userID == userID {
			result = append(result, Item{
				UUID:     strconv.Itoa(item.cur),
				URL:      item.URL,
				ShortURL: shortURL,
			})
		}
	}

	return result, nil
}

// DeleteByUser удаляет из хранилища все сокращенные ссылки пользователя.
func (s InMemory) DeleteByUser(ctx context.Context, userID uuid.UUID, shortURL string) error {
	if val, ok := s.storage[shortURL]; ok && val.userID == userID {
		delete(s.storage, shortURL)
	}

	return nil
}

// GetStats возвращает статистику хранилища.
func (s InMemory) GetStats(ctx context.Context) (*Stats, error) {
	stats := Stats{
		Users: 0,
		URLs:  0,
	}
	userStats := make(map[uuid.UUID]bool)

	for _, v := range s.storage {
		if _, ok := userStats[v.userID]; !ok {
			userStats[v.userID] = true
			stats.Users++
		}
		stats.URLs++
	}

	return &stats, nil
}
