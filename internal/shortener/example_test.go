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

func ExampleShortener_GetURL() {
	ctx, shortener := setupExample()
	userID, _ := uuid.NewV6()

	shortURL, _ := shortener.Container.Storage.Set(ctx, userID, "https://ya.ru")

	result, err := shortener.GetURL(ctx, shortURL)

	if err != nil {
		fmt.Printf("Got error: %s", err)
		return
	}

	fmt.Printf("Your URL is: %s\n", result)

	// Output:
	// Your URL is: https://ya.ru
}

func ExampleShortener_MultiShorten() {
	ctx, shortener := setupExample()

	userID, _ := uuid.NewV6()

	result, err := shortener.MultiShorten(ctx, userID, storage.Items{
		storage.Item{
			URL: "https://ya.ru",
		},
		storage.Item{
			URL: "https://r0.ru",
		},
	})
	if err != nil {
		fmt.Printf("Could not batch insert: %v\n", err)
		return
	}
	for _, shortURL := range result {
		fmt.Printf("Your shortened URL is: %s\n", shortURL.URL)
	}

	// Output:
	// Your shortened URL is: xA
	// Your shortened URL is: yA
}

func ExampleShortener_GetURLsByUser() {
	ctx, shortener := setupExample()
	userID, _ := uuid.NewV6()

	shortener.Container.Storage.Set(ctx, userID, "https://ya.ru")

	result, err := shortener.GetURLsByUser(ctx, userID)

	if err != nil {
		fmt.Printf("Got error: %s", err)
		return
	}

	fmt.Printf("Your URLs are: %s\n", result)

	// Output:
	// Your URLs are: [{1 https://ya.ru http://localhost:8080/xA}]
}
