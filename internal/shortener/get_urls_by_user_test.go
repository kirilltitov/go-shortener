package shortener

import (
	"fmt"

	"github.com/google/uuid"
)

func ExampleShortener_GetURLsByUser() {
	ctx, shortener := setupExample()
	userID, _ := uuid.NewV6()

	shortener.container.Storage.Set(ctx, userID, "https://ya.ru")

	result, err := shortener.GetURLsByUser(ctx, userID)

	if err != nil {
		fmt.Printf("Got error: %s", err)
		return
	}

	fmt.Printf("Your URLs are: %s\n", result)

	// Output:
	// Your URLs are: [{1 https://ya.ru http://localhost:8080/xA}]
}
