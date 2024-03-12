package handlers

import (
	"fmt"
	"github.com/jxskiss/base62"
	"github.com/kirilltitov/go-shortener/internal/config"
	internalStorage "github.com/kirilltitov/go-shortener/internal/storage"
	"github.com/kirilltitov/go-shortener/internal/utils"
	"io"
	"net/http"
	"strconv"
)

func HandlerCreateShortURL(w http.ResponseWriter, r *http.Request, storage internalStorage.Storage, cur *int) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Could not read body: %s\n", err.Error()))
		return
	}

	link := string(b)
	if !utils.IsValidLink(link) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("Invalid link (must start with https:// or http://): %s\n", link))
		return
	}

	*cur++
	storage.Set(*cur, link)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("%s/%s", config.BaseURL, base62.EncodeToString([]byte(strconv.Itoa(*cur)))))
}
