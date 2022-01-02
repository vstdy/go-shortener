package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vstdy0/go-project/service/shortener"
)

func Router(service shortener.URLService) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Route("/", func(r chi.Router) {
		r.Post("/", CreateShortcut(service))
		r.Post("/api/shorten", CreateShortcut(service))
		r.Get("/{id}", GetShortcut(service))
	})
	return r
}
