package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirilltitov/go-shortener/internal/config"
	"github.com/kirilltitov/go-shortener/internal/container"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

func setupExample() (context.Context, Shortener) {
	ctx := context.Background()

	cfg := config.NewWithoutParsing()
	cfg.DatabaseDSN = ""
	cfg.FileStoragePath = ""

	cnt, _ := container.New(ctx, cfg)

	shortener := New(cfg, cnt)

	return ctx, shortener
}

func ExampleShortener_ShortenURL() {
	ctx, shortener := setupExample()

	userID, _ := uuid.NewV6()

	shortURL, err := shortener.ShortenURL(ctx, userID, "https://ya.ru/")
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			fmt.Println("URL already exists")
			return
		} else {
			fmt.Println("Internal server error")
			return
		}
	}
	fmt.Printf("Your shortened URL is: %s\n", shortURL)

	// Output:
	// Your shortened URL is: http://localhost:8080/xA
}
