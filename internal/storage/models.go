package storage

// Item являет собою запись хранилища о сокращенной ссылке.
type Item struct {
	UUID     string // UUID является первичным ключом сокращенной ссылки.
	URL      string // URL является настоящим адресом сокращенной ссылки.
	ShortURL string // ShortURL является сокращенным идентификатором ссылки.
}

// Items является списком из объектом сокращенных ссылок.
type Items []Item

// Stats хранит в себе статистику хранилища
type Stats struct {
	Users int `json:"users" db:"users"` // Users является количеством пользователей сокращателя урлов
	URLs  int `json:"urls"  db:"urls"`  // URLs является количеством сокращенных урлов
}
