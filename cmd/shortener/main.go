package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var id int

var storage = make(map[string]storageURL)

type storageURL struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := strings.Split(r.URL.Path, "/")
		if len(path) == 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fullURL := storage[path[1]].URL
		if fullURL == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id++
		urlID := strconv.Itoa(id)
		storage[urlID] = storageURL{urlID, string(body)}
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://" + r.Host + "/" + storage[urlID].ID))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", shortenURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
