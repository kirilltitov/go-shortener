package storage

import "errors"

var ErrNotFound = errors.New("key was not found")
var ErrDeleted = errors.New("URL has been deleted")
