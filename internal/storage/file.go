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

type setter interface {
	Set(int, string)
}

func LoadStorageFromFile(path string, s setter, cur *int) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		logger.Log.Infof("No storage file '%s'", path)
		return nil
	}

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

	return file.Close()
}

func SaveRowToFile(path string, cur int, shortURL, URL string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(row{
		UUID:        strconv.Itoa(cur),
		ShortURL:    shortURL,
		OriginalURL: URL,
	}); err != nil {
		return err
	}

	return file.Close()
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
