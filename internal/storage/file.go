package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// File является файловым хранилищем сокращенных ссылок.
type File struct {
	InMemory
	path string
}

type fileRow struct {
	UUID        string    `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	UserID      uuid.UUID `json:"user_id"`
}

// NewFileStorage создает, конфигурирует и возвращает экземпляр объекта файлового хранилища для заданного пути хранения.
func NewFileStorage(ctx context.Context, path string) (*File, error) {
	result := &File{
		InMemory: *NewInMemoryStorage(ctx),
		path:     path,
	}

	if err := result.LoadStorageFromFile(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

// Set записывает в хранилище информацию о сокращенной ссылке.
func (f *File) Set(ctx context.Context, userID uuid.UUID, URL string) (string, error) {
	shortURL, err := f.InMemory.Set(ctx, userID, URL)
	if err != nil {
		return "", err
	}

	if err := f.saveRowToFile(*f.cur, userID, shortURL, URL); err != nil {
		return "", err
	}

	return shortURL, nil
}

// MultiSet записывает в хранилище информацию о нескольких сокращенных ссылках.
func (f *File) MultiSet(ctx context.Context, userID uuid.UUID, items Items) (Items, error) {
	var result Items

	for _, item := range items {
		shortURL, err := f.Set(ctx, userID, item.URL)
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

// LoadStorageFromFile загружает хранилище из файла.
func (f *File) LoadStorageFromFile(ctx context.Context) error {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0444)
	if err != nil {
		logger.Log.Infof("No storage file '%s'", f.path)
		return nil
	}

	decoder := json.NewDecoder(file)

	for {
		var r fileRow
		if err := decoder.Decode(&r); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			logger.Log.Fatal(err)
			continue
		}
		if _, err := f.InMemory.Set(ctx, r.UserID, r.OriginalURL); err != nil {
			return err
		}
		logger.Log.Infof("Loaded row %+v from file", r)
	}

	return file.Close()
}

func (f *File) saveRowToFile(idx int, userID uuid.UUID, shortURL, URL string) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(fileRow{
		UUID:        strconv.Itoa(idx),
		UserID:      userID,
		ShortURL:    shortURL,
		OriginalURL: URL,
	}); err != nil {
		return err
	}

	return file.Close()
}

// WipeFileStorage безусловно очищает файловое хранилище.
func (f *File) WipeFileStorage() {
	if _, err := os.Stat(f.path); errors.Is(err, os.ErrNotExist) {
		logger.Log.Infof("Storage file '%s' doesn't exist, nothing to wipe", f.path)
		return
	} else if err != nil {
		panic(err)
	}
	if err := os.Remove(f.path); err != nil {
		panic(err)
	}
}
