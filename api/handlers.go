package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener"
	"io"
	"net/http"
)

type Handler struct {
	service shortener.URLService
	cfg     config.Config
}

func NewHandler(service shortener.URLService, cfg config.Config) Handler {
	return Handler{service: service, cfg: cfg}
}

func (h Handler) createShortcut(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
		res, err = h.jsonResponse(userID, body)
	default:
		res, err = h.plainResponse(userID, body)
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

func (h Handler) getShortcut(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	url := h.service.GetURL(urlID)
	if url == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var res []byte
	urls := h.service.GetUserURLs(userID)
	userURLs := fromCanonical(urls, h.cfg.BaseURL)

	if userURLs != nil {
		marshal, err := json.Marshal(userURLs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = marshal
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	if _, err := w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
