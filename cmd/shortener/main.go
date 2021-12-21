package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var id int

var urls = make(map[string]shortenedURL)

type shortenedURL struct {
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
		fullURL := urls[path[1]].URL
		if fullURL == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		body := r.PostForm["url"]
		if body == nil {
			http.Error(w, "Request must have 'url' field", http.StatusBadRequest)
			return
		}
		id++
		shortenedID := strconv.Itoa(id)
		urls[shortenedID] = shortenedURL{shortenedID, body[0]}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(r.Host + "/" + urls[shortenedID].ID))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", shortenURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
