package storage

import "context"

type InMemory struct {
	s   map[int]string
	cur *int
}

func NewInMemoryStorage(ctx context.Context) *InMemory {
	return &InMemory{
		s:   make(map[int]string),
		cur: new(int),
	}
}

func (s InMemory) Get(ctx context.Context, shortURL string) (string, error) {
	i, err := shortURLToInt(shortURL)
	if err != nil {
		return "", err
	}

	var _err error = nil
	val, ok := s.s[i]
	if !ok {
		_err = ErrNotFound
	}

	return val, _err
}

func (s InMemory) Set(ctx context.Context, URL string) (string, error) {
	*s.cur++
	s.s[*s.cur] = URL

	return intToShortURL(*s.cur), nil
}
