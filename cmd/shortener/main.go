package main

import (
	"fmt"
	internal "github.com/kirilltitov/go-shortener/internal/app"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", internal.Handler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(fmt.Sprintf("Could not start server: %s\n", err.Error()))
	}
}
