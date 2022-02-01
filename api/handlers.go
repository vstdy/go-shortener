package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/service/shortener"
)

type Handler struct {
	service shortener.URLService
	cfg     common.Config
}

func NewHandler(service shortener.URLService, cfg common.Config) Handler {
	return Handler{service: service, cfg: cfg}
}

func (h Handler) createShortcut(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "context: failed to retrieve user_id", http.StatusInternalServerError)
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
		res, err = h.jsonResponse(r.Context(), userID, body)
	default:
		res, err = h.plainResponse(r.Context(), userID, body)
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
	urlID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.service.GetURL(r.Context(), urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	urls, err := h.service.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userURLs := userURLsFromCanonical(urls, h.cfg.BaseURL)

	var res []byte
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

func (h Handler) getPing(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", h.cfg.PSQLStorage.DSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
