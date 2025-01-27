package shortener

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kirilltitov/go-shortener/internal/storage"
)

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
