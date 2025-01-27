package storage

// Item являет собою запись хранилища о сокращенной ссылке.
type Item struct {
	// UUID является первичным ключом сокращенной ссылки.
	UUID string

	// URL является настоящим адресом сокращенной ссылки.
	URL string

	// ShortURL является сокращенным идентификатором ссылки.
	ShortURL string
}

// Items является списком из объектом сокращенных ссылок.
type Items []Item
