package main

import (
	"fediverse/jrd"
	"fediverse/webfinger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(webfinger.WebFinger(func(resource string) (jrd.JRD, error) {
		return jrd.JRD{}, nil
	}))
	http.ListenAndServe(":3000", r)
}
