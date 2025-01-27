package storage

import "errors"

// ErrNotFound возвращается пользователю когда искомой записи нет в хранилище.
var ErrNotFound = errors.New("key was not found")

// ErrDeleted возвращается пользователю когда короткая ссылка была удалена из хранилища.
var ErrDeleted = errors.New("URL has been deleted")
