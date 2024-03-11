package main

import (
	"fmt"
	"github.com/jxskiss/base62"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var cur int = 0
var storage = map[int]string{}

func isValidLink(maybeLink string) bool {
	return strings.HasPrefix(maybeLink, "https://")
}

func handler(w http.ResponseWriter, r *http.Request) {
	var result string
	var code int

	w.Header().Set("Content-Type", "text/plain")

	defer func() {
		fmt.Sprintf("Returning code %d with text %s", code, result)
		w.WriteHeader(code)
		io.WriteString(w, result)
	}()

	switch r.Method {
	case http.MethodGet:
		shortURL := strings.TrimPrefix(r.URL.Path, "/")
		decodedStringInt, err := base62.DecodeString(shortURL)
		if err != nil {
			result = fmt.Sprintf("Could not decode short url '%s'\n", shortURL)
			code = http.StatusBadRequest
			break
		}
		decodedInt, err := strconv.Atoi(string(decodedStringInt))
		if err != nil {
			result = fmt.Sprintf("Could not decode short url '%s'\n", shortURL)
			code = http.StatusBadRequest
			break
		}
		url, ok := storage[decodedInt]
		if !ok {
			result = fmt.Sprintf("URL '%s' not found\n", shortURL)
			code = http.StatusBadRequest
			break
		}

		code = http.StatusTemporaryRedirect

		http.Redirect(w, r, url, code)

		break
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			result = fmt.Sprintf("Could not read body: %s\n", err.Error())
			code = http.StatusBadRequest
			break
		}
		link := string(b)
		if !isValidLink(link) {
			result = fmt.Sprintf("Invalid link (must start with https://): %s\n", link)
			code = http.StatusBadRequest
			break
		}

		cur++
		storage[cur] = link

		curString := strconv.Itoa(cur)
		fmt.Printf("%s\n", curString)
		fmt.Printf("bytes: %v\n", []byte(curString))
		fmt.Printf("result: %v\n", base62.EncodeToString([]byte(curString)))
		result = fmt.Sprintf("http://localhost:8080/%s\n", base62.EncodeToString([]byte(curString)))
		code = http.StatusCreated

		break
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(fmt.Sprintf("Could not start server: %s\n", err.Error()))
	}
}
