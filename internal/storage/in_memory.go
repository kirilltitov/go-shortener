package storage

type InMemory map[int]string

func (s InMemory) Get(key int) (string, bool) {
	val, ok := s[key]
	return val, ok
}

func (s InMemory) Set(key int, value string) {
	s[key] = value
}
