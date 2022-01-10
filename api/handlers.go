package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener"
	"io"
	"net/http"
)

func jsonResponse(service shortener.URLService, cfg config.Config, body []byte) ([]byte, error) {
	var urlRequest ShortenURLRequest
	err := json.Unmarshal(body, &urlRequest)
	if err != nil {
		return nil, err
	}
	id, err := service.AddURL(urlRequest.URL)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(ShortenURLResponse{Result: cfg.BaseURL + "/" + id})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func plainResponse(service shortener.URLService, cfg config.Config, body []byte) ([]byte, error) {
	id, err := service.AddURL(string(body))
	if err != nil {
		return nil, err
	}
	return []byte(cfg.BaseURL + "/" + id), nil
}

func CreateShortcut(service shortener.URLService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var res []byte
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			res, err = jsonResponse(service, cfg, body)
		default:
			res, err = plainResponse(service, cfg, body)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write(res); err != nil {
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
