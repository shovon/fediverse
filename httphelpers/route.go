package httphelpers

import (
	"context"
	"fediverse/pathhelpers"
	"net/http"
)

type contextValue struct {
	key string
}

// TODO: use the pathhelpers library

func Route(route string) Processor {
	return ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hasMatch, params := pathhelpers.Match(route, r.URL.Path)

				if !hasMatch {
					next.ServeHTTP(w, r)
					return
				}

				newR := r.WithContext(r.Context())
				for key, value := range params {
					newR = newR.WithContext(context.WithValue(newR.Context(), contextValue{key}, value))
				}

				middleware(next).ServeHTTP(w, newR)
			})
		}
	})
}
