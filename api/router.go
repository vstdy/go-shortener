package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener"
)

func Router(service shortener.URLService, cfg config.Config) chi.Router {
	h := NewHandler(service, cfg)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.createShortcut)
		r.Post("/api/shorten", h.createShortcut)
		r.Get("/{id}", h.getShortcut)
	})

	return r
}
