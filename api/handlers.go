package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"io"
	"net/http"
)

func CreateShortcut(w http.ResponseWriter, r *http.Request) {
	var id string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id = shortener.AddURL(string(body))
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://" + r.Host + "/" + id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetShortcut(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	url := shortener.GetURL(urlID)
	if url == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
