package main

import (
	"errors"
	"fediverse/jrd"
	"fediverse/nullable"
	"fediverse/webfinger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func toHandlerFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use("/.well_known", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	r.Route("/.well-known", func(r chi.Router) {
		r.Get("/webfinger", toHandlerFunc(webfinger.CreateHandler(func(resource string) (jrd.JRD, error) {
			return jrd.JRD{
				Subject: nullable.Just(resource),

				Links: nullable.Just([]jrd.Link{}),
			}, errors.New("not yet implemented")
		})))
		r.Get("/nodeinfo", func(w http.ResponseWriter, r *http.Request) {

		})
	})
	http.ListenAndServe(":3000", r)
}
