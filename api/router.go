package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-project/service/shortener"
)

func Router(svc shortener.URLService, config common.Config) chi.Router {
	h := NewHandler(svc, config)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(config.Timeout))
	r.Use(gzipDecompressRequest)
	r.Use(gzipCompressResponse)
	r.Use(cookieAuth(config.SecretKey))

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.shortenURL)
		r.Get("/{id}", h.getShortenedURL)
		r.Get("/ping", h.Ping)

		r.Route("/api", func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			r.Post("/shorten", h.shortenURL)
			r.Post("/shorten/batch", h.shortenBatchURLs)
			r.Get("/user/urls", h.getUserURLs)
			r.Delete("/user/urls", h.deleteUserURLs)
		})
	})

	return r
}
