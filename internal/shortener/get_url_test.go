package shortener

import (
	"fmt"

	"github.com/google/uuid"
)

func ExampleShortener_GetURL() {
	ctx, shortener := setupExample()
	userID, _ := uuid.NewV6()

	shortURL, _ := shortener.container.Storage.Set(ctx, userID, "https://ya.ru")

	result, err := shortener.GetURL(ctx, shortURL)

	if err != nil {
		fmt.Printf("Got error: %s", err)
		return
	}

	fmt.Printf("Your URL is: %s\n", result)

	// Output:
	// Your URL is: https://ya.ru
}
