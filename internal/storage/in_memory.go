package storage

import "context"

type inMemoryRow struct {
	cur int
	URL string
}

type InMemory struct {
	storage map[string]inMemoryRow
	cur     *int
}

func NewInMemoryStorage(ctx context.Context) *InMemory {
	return &InMemory{
		storage: make(map[string]inMemoryRow),
		cur:     new(int),
	}
}

func (s InMemory) Get(ctx context.Context, shortURL string) (string, error) {
	var _err error = nil
	val, ok := s.storage[shortURL]
	if !ok {
		_err = ErrNotFound
	}

	return val.URL, _err
}

func (s InMemory) Set(ctx context.Context, URL string) (string, error) {
	*s.cur++
	shortURL := intToShortURL(*s.cur)
	s.storage[shortURL] = inMemoryRow{
		cur: *s.cur,
		URL: URL,
	}

	return shortURL, nil
}

func (s InMemory) MultiSet(ctx context.Context, items Items) (Items, error) {
	var result Items

	for _, item := range items {
		shortURL, err := s.Set(ctx, item.URL)
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
