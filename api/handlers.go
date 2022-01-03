package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener"
	"github.com/vstdy0/go-project/storage/inmemory"
	"io"
	"net/http"
)

func CreateShortcut(service shortener.URLService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var url, res string
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var urlModel inmemory.URLModel
			if err := json.Unmarshal(body, &urlModel); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			url = urlModel.URL
			res = `{"result": "%s/%s"}`
		default:
			url = string(body)
			res = "%s/%s"
		}
		id := service.AddURL(url)
		res = fmt.Sprintf(res, cfg.BaseURL, id)

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write([]byte(res)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetShortcut(service shortener.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "id")
		url := service.GetURL(urlID)
		if url == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
