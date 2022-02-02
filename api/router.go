package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/service/shortener"
)

func Router(svc shortener.URLService, cfg common.Config) chi.Router {
	h := NewHandler(svc, cfg)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(cfg.RequestTimeout))
	r.Use(gzipDecompressRequest)
	r.Use(gzipCompressResponse)
	r.Use(cookieAuth(cfg.SecretKey))

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.shortenURL)
		r.Post("/api/shorten", h.shortenURL)
		r.Post("/api/shorten/batch", h.shortenBatchURLs)
		r.Get("/{id}", h.getShortenURL)
		r.Get("/user/urls", h.getUserURLs)
		r.Get("/ping", h.getPing)
	})

	return r
}
