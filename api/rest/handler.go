package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"

	"github.com/vstdy0/go-shortener/api/rest/model"
	"github.com/vstdy0/go-shortener/pkg"
	"github.com/vstdy0/go-shortener/pkg/logging"
	"github.com/vstdy0/go-shortener/pkg/tracing"
	"github.com/vstdy0/go-shortener/service/shortener"
)

const (
	serviceName = "Shortener server"
)

// Handler keeps handler dependencies.
type Handler struct {
	config  Config
	service shortener.Service
}

// NewHandler returns a new Handler instance.
func NewHandler(service shortener.Service, config Config) Handler {
	return Handler{service: service, config: config}
}

// Logger returns logger with service field set.
func (h Handler) Logger(ctx context.Context) (context.Context, zerolog.Logger) {
	ctx, logger := logging.GetCtxLogger(ctx, logging.WithLogLevel(h.config.LogLevel))
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return ctx, logger
}

// shortenURL creates shortcut for given url.
func (h Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())
	ctx, span := tracing.StartSpanFromCtx(ctx, "Shortening URL")
	defer tracing.FinishSpan(span, nil)

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
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
		res, err = h.jsonURLResponse(ctx, userID, body)
	default:
		res, err = h.plainURLResponse(ctx, userID, body)
	}
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !errors.Is(err, pkg.ErrAlreadyExists) {
			logger.Warn().Err(err).Msg("Shortening URL:")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

// shortenURLsBatch creates shortcuts for given urls batch.
func (h Handler) shortenURLsBatch(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())
	ctx, span := tracing.StartSpanFromCtx(ctx, "Shortening URLs batch")
	defer tracing.FinishSpan(span, nil)

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
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

	res, err := h.urlsBatchResponse(ctx, userID, body)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !errors.Is(err, pkg.ErrAlreadyExists) {
			logger.Warn().Err(err).Msg("Shortening URLs batch:")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusConflict)
		if _, err = w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

// getOriginalURL returns origin url from shortcut.
func (h Handler) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())
	ctx, span := tracing.StartSpanFromCtx(ctx, "Getting original URL")
	defer tracing.FinishSpan(span, nil)

	urlID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.service.GetURL(ctx, urlID)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Warn().Err(err).Msg("Getting original URL:")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if url == "" {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// getUsersURLs returns urls created by current user.
func (h Handler) getUsersURLs(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())
	ctx, span := tracing.StartSpanFromCtx(ctx, "Getting user URLs")
	defer tracing.FinishSpan(span, nil)

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "context: failed to retrieve user_id", http.StatusInternalServerError)
		return
	}

	urls, err := h.service.GetUserURLs(ctx, userID)
	if err != nil {
		logger.Warn().Err(err).Msg("Getting user URLs:")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userURLs := model.NewUserURLsFromCanon(urls, h.config.BaseURL)

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

// deleteUserURLs removes urls created by current user.
func (h Handler) deleteUserURLs(w http.ResponseWriter, r *http.Request) {
	ctx, logger := h.Logger(r.Context())
	ctx, span := tracing.StartSpanFromCtx(ctx, "Deleting user URLs")
	defer tracing.FinishSpan(span, nil)

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
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

	objs, err := model.URLsToDelToCanon(ids, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.RemoveUserURLs(ctx, objs)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Warn().Err(err).Msg("Deleting user URLs:")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// ping checks connection to database.
func (h Handler) ping(w http.ResponseWriter, r *http.Request) {
	_, logger := h.Logger(r.Context())

	if err := h.service.Ping(); err != nil {
		logger.Warn().Err(err).Msg("Ping:")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
