package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

type row struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func LoadStorageFromFile(path string, s Storage, cur *int) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		logger.Log.Infof("No storage file '%s'", path)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := json.NewDecoder(file)

	for {
		var r row
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			logger.Log.Fatal(err)
			continue
		}
		*cur++
		s.Set(*cur, r.OriginalURL)
		logger.Log.Infof("Loaded row %+v from file", r)
	}
}

func SaveRowToFile(path string, cur int, shortURL, URL string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(row{
		UUID:        strconv.Itoa(cur),
		ShortURL:    shortURL,
		OriginalURL: URL,
	}); err != nil {
		return err
	}

	return nil
}

func WipeFileStorage(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Log.Infof("Storage file '%s' doesn't exist, nothing to wipe", path)
		return
	} else if err != nil {
		panic(err)
	}
	if err := os.Remove(path); err != nil {
		panic(err)
	}
}
