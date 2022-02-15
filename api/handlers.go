package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/vstdy0/go-project/api/model"
	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/pkg"
	"github.com/vstdy0/go-project/service/shortener"
)

type Handler struct {
	service shortener.URLService
	config  common.Config
}

func NewHandler(service shortener.URLService, config common.Config) Handler {
	return Handler{service: service, config: config}
}

func (h Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
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
		res, err = h.jsonURLResponse(r.Context(), userID, body)
	default:
		res, err = h.plainURLResponse(r.Context(), userID, body)
	}
	if err != nil {
		if !errors.Is(err, pkg.ErrIntegrityViolation) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) shortenBatchURLs(w http.ResponseWriter, r *http.Request) {
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

	contentType := r.Header.Get("Content-Type")

	res, err := h.urlsBatchResponse(r.Context(), userID, body)
	if err != nil {
		if !errors.Is(err, pkg.ErrIntegrityViolation) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusConflict)
		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) getShortenedURL(w http.ResponseWriter, r *http.Request) {
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

	if url == "" {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h Handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "context: failed to retrieve user_id", http.StatusInternalServerError)
		return
	}

	urls, err := h.service.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userURLs := model.UserURLsFromCanonical(urls, h.config.BaseURL)

	if userURLs == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	res, err := json.Marshal(userURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) deleteUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "context: failed to retrieve user_id", http.StatusInternalServerError)
		return
	}

	var ids []string
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	objs, err := model.URLsToDeleteToCanonical(ids, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUserURLs(objs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
