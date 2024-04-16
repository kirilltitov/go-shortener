package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

type File struct {
	InMemory
	path string
}

type fileRow struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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

func (f *File) Set(ctx context.Context, URL string) (string, error) {
	shortURL, err := f.InMemory.Set(ctx, URL)
	if err != nil {
		return "", err
	}

	if err := f.saveRowToFile(shortURL, URL); err != nil {
		return "", err
	}

	return shortURL, nil
}

func (f *File) LoadStorageFromFile(ctx context.Context) error {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0444)
	if err != nil {
		logger.Log.Infof("No storage file '%s'", f.path)
		return nil
	}

	decoder := json.NewDecoder(file)

	for {
		var r fileRow
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			logger.Log.Fatal(err)
			continue
		}
		if _, err := f.InMemory.Set(ctx, r.OriginalURL); err != nil {
			return err
		}
		logger.Log.Infof("Loaded row %+v from file", r)
	}

	return file.Close()
}

func (f *File) saveRowToFile(shortURL, URL string) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(fileRow{
		UUID:        strconv.Itoa(*f.cur),
		ShortURL:    shortURL,
		OriginalURL: URL,
	}); err != nil {
		return err
	}

	return file.Close()
}

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
