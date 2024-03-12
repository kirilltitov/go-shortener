package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kirilltitov/go-shortener/internal/app/handlers"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"log"
	"net/http"
)

var cur int = 0
var storage internalStorage.Storage = internalStorage.InMemory{}

func ShortenerRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/{short}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetShortURL(w, r, storage)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateShortURL(w, r, storage, &cur)
	})

	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", ShortenerRouter()))
}
