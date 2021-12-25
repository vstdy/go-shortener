package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vstdy0/go-project/model"
	"io"
	"net/http"
	"strconv"
)

var id int

var urls = make(map[string]model.URLModel)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id++
			urlID := strconv.Itoa(id)
			urls[urlID] = model.URLModel{ID: urlID, URL: string(body)}
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("http://" + r.Host + "/" + urls[urlID].ID))
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			urlID := chi.URLParam(r, "id")
			fullURL := urls[urlID].URL
			if fullURL == "" {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			w.Header().Set("Location", fullURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		})
	})
	return r
}
